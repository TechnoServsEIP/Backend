<<<<<<< HEAD
// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package proto

/*
 * Types and routines for supporting protocol buffer extensions.
 */

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
)

// ErrMissingExtension is the error returned by GetExtension if the named extension is not in the message.
var ErrMissingExtension = errors.New("proto: missing extension")

// ExtensionRange represents a range of message extensions for a protocol buffer.
// Used in code generated by the protocol compiler.
type ExtensionRange struct {
	Start, End int32 // both inclusive
}

// extendableProto is an interface implemented by any protocol buffer generated by the current
// proto compiler that may be extended.
type extendableProto interface {
	Message
	ExtensionRangeArray() []ExtensionRange
	extensionsWrite() map[int32]Extension
	extensionsRead() (map[int32]Extension, sync.Locker)
}

// extendableProtoV1 is an interface implemented by a protocol buffer generated by the previous
// version of the proto compiler that may be extended.
type extendableProtoV1 interface {
	Message
	ExtensionRangeArray() []ExtensionRange
	ExtensionMap() map[int32]Extension
}

// extensionAdapter is a wrapper around extendableProtoV1 that implements extendableProto.
type extensionAdapter struct {
	extendableProtoV1
}

func (e extensionAdapter) extensionsWrite() map[int32]Extension {
	return e.ExtensionMap()
}

func (e extensionAdapter) extensionsRead() (map[int32]Extension, sync.Locker) {
	return e.ExtensionMap(), notLocker{}
}

// notLocker is a sync.Locker whose Lock and Unlock methods are nops.
type notLocker struct{}

func (n notLocker) Lock()   {}
func (n notLocker) Unlock() {}

// extendable returns the extendableProto interface for the given generated proto message.
// If the proto message has the old extension format, it returns a wrapper that implements
// the extendableProto interface.
func extendable(p interface{}) (extendableProto, error) {
	switch p := p.(type) {
	case extendableProto:
		if isNilPtr(p) {
			return nil, fmt.Errorf("proto: nil %T is not extendable", p)
		}
		return p, nil
	case extendableProtoV1:
		if isNilPtr(p) {
			return nil, fmt.Errorf("proto: nil %T is not extendable", p)
		}
		return extensionAdapter{p}, nil
	}
	// Don't allocate a specific error containing %T:
	// this is the hot path for Clone and MarshalText.
	return nil, errNotExtendable
}

var errNotExtendable = errors.New("proto: not an extendable proto.Message")

func isNilPtr(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

// XXX_InternalExtensions is an internal representation of proto extensions.
//
// Each generated message struct type embeds an anonymous XXX_InternalExtensions field,
// thus gaining the unexported 'extensions' method, which can be called only from the proto package.
//
// The methods of XXX_InternalExtensions are not concurrency safe in general,
// but calls to logically read-only methods such as has and get may be executed concurrently.
type XXX_InternalExtensions struct {
	// The struct must be indirect so that if a user inadvertently copies a
	// generated message and its embedded XXX_InternalExtensions, they
	// avoid the mayhem of a copied mutex.
	//
	// The mutex serializes all logically read-only operations to p.extensionMap.
	// It is up to the client to ensure that write operations to p.extensionMap are
	// mutually exclusive with other accesses.
	p *struct {
		mu           sync.Mutex
		extensionMap map[int32]Extension
	}
}

// extensionsWrite returns the extension map, creating it on first use.
func (e *XXX_InternalExtensions) extensionsWrite() map[int32]Extension {
	if e.p == nil {
		e.p = new(struct {
			mu           sync.Mutex
			extensionMap map[int32]Extension
		})
		e.p.extensionMap = make(map[int32]Extension)
	}
	return e.p.extensionMap
}

// extensionsRead returns the extensions map for read-only use.  It may be nil.
// The caller must hold the returned mutex's lock when accessing Elements within the map.
func (e *XXX_InternalExtensions) extensionsRead() (map[int32]Extension, sync.Locker) {
	if e.p == nil {
		return nil, nil
	}
	return e.p.extensionMap, &e.p.mu
}

// ExtensionDesc represents an extension specification.
// Used in generated code from the protocol compiler.
type ExtensionDesc struct {
	ExtendedType  Message     // nil pointer to the type that is being extended
	ExtensionType interface{} // nil pointer to the extension type
	Field         int32       // field number
	Name          string      // fully-qualified name of extension, for text formatting
	Tag           string      // protobuf tag style
	Filename      string      // name of the file in which the extension is defined
}

func (ed *ExtensionDesc) repeated() bool {
	t := reflect.TypeOf(ed.ExtensionType)
	return t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
}

// Extension represents an extension in a message.
type Extension struct {
	// When an extension is stored in a message using SetExtension
	// only desc and value are set. When the message is marshaled
	// enc will be set to the encoded form of the message.
	//
	// When a message is unmarshaled and contains extensions, each
	// extension will have only enc set. When such an extension is
	// accessed using GetExtension (or GetExtensions) desc and value
	// will be set.
	desc *ExtensionDesc

	// value is a concrete value for the extension field. Let the type of
	// desc.ExtensionType be the "API type" and the type of Extension.value
	// be the "storage type". The API type and storage type are the same except:
	//	* For scalars (except []byte), the API type uses *T,
	//	while the storage type uses T.
	//	* For repeated fields, the API type uses []T, while the storage type
	//	uses *[]T.
	//
	// The reason for the divergence is so that the storage type more naturally
	// matches what is expected of when retrieving the values through the
	// protobuf reflection APIs.
	//
	// The value may only be populated if desc is also populated.
	value interface{}

	// enc is the raw bytes for the extension field.
	enc []byte
}

// SetRawExtension is for testing only.
func SetRawExtension(base Message, id int32, b []byte) {
	epb, err := extendable(base)
	if err != nil {
		return
	}
	extmap := epb.extensionsWrite()
	extmap[id] = Extension{enc: b}
}

// isExtensionField returns true iff the given field number is in an extension range.
func isExtensionField(pb extendableProto, field int32) bool {
	for _, er := range pb.ExtensionRangeArray() {
		if er.Start <= field && field <= er.End {
			return true
		}
	}
	return false
}

// checkExtensionTypes checks that the given extension is valid for pb.
func checkExtensionTypes(pb extendableProto, extension *ExtensionDesc) error {
	var pbi interface{} = pb
	// Check the extended type.
	if ea, ok := pbi.(extensionAdapter); ok {
		pbi = ea.extendableProtoV1
	}
	if a, b := reflect.TypeOf(pbi), reflect.TypeOf(extension.ExtendedType); a != b {
		return fmt.Errorf("proto: bad extended type; %v does not extend %v", b, a)
	}
	// Check the range.
	if !isExtensionField(pb, extension.Field) {
		return errors.New("proto: bad extension number; not in declared ranges")
	}
	return nil
}

// extPropKey is sufficient to uniquely identify an extension.
type extPropKey struct {
	base  reflect.Type
	field int32
}

var extProp = struct {
	sync.RWMutex
	m map[extPropKey]*Properties
}{
	m: make(map[extPropKey]*Properties),
}

func extensionProperties(ed *ExtensionDesc) *Properties {
	key := extPropKey{base: reflect.TypeOf(ed.ExtendedType), field: ed.Field}

	extProp.RLock()
	if prop, ok := extProp.m[key]; ok {
		extProp.RUnlock()
		return prop
	}
	extProp.RUnlock()

	extProp.Lock()
	defer extProp.Unlock()
	// Check again.
	if prop, ok := extProp.m[key]; ok {
		return prop
	}

	prop := new(Properties)
	prop.Init(reflect.TypeOf(ed.ExtensionType), "unknown_name", ed.Tag, nil)
	extProp.m[key] = prop
	return prop
}

// HasExtension returns whether the given extension is present in pb.
func HasExtension(pb Message, extension *ExtensionDesc) bool {
	// TODO: Check types, field numbers, etc.?
	epb, err := extendable(pb)
	if err != nil {
		return false
	}
	extmap, mu := epb.extensionsRead()
	if extmap == nil {
		return false
	}
	mu.Lock()
	_, ok := extmap[extension.Field]
	mu.Unlock()
	return ok
}

// ClearExtension removes the given extension from pb.
func ClearExtension(pb Message, extension *ExtensionDesc) {
	epb, err := extendable(pb)
	if err != nil {
		return
	}
	// TODO: Check types, field numbers, etc.?
	extmap := epb.extensionsWrite()
	delete(extmap, extension.Field)
=======
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"errors"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type (
	// ExtensionDesc represents an extension descriptor and
	// is used to interact with an extension field in a message.
	//
	// Variables of this type are generated in code by protoc-gen-go.
	ExtensionDesc = protoimpl.ExtensionInfo

	// ExtensionRange represents a range of message extensions.
	// Used in code generated by protoc-gen-go.
	ExtensionRange = protoiface.ExtensionRangeV1

	// Deprecated: Do not use; this is an internal type.
	Extension = protoimpl.ExtensionFieldV1

	// Deprecated: Do not use; this is an internal type.
	XXX_InternalExtensions = protoimpl.ExtensionFields
)

// ErrMissingExtension reports whether the extension was not present.
var ErrMissingExtension = errors.New("proto: missing extension")

var errNotExtendable = errors.New("proto: not an extendable proto.Message")

// HasExtension reports whether the extension field is present in m
// either as an explicitly populated field or as an unknown field.
func HasExtension(m Message, xt *ExtensionDesc) (has bool) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() {
		return false
	}

	// Check whether any populated known field matches the field number.
	xtd := xt.TypeDescriptor()
	if isValidExtension(mr.Descriptor(), xtd) {
		has = mr.Has(xtd)
	} else {
		mr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
			has = int32(fd.Number()) == xt.Field
			return !has
		})
	}

	// Check whether any unknown field matches the field number.
	for b := mr.GetUnknown(); !has && len(b) > 0; {
		num, _, n := protowire.ConsumeField(b)
		has = int32(num) == xt.Field
		b = b[n:]
	}
	return has
}

// ClearExtension removes the the exntesion field from m
// either as an explicitly populated field or as an unknown field.
func ClearExtension(m Message, xt *ExtensionDesc) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() {
		return
	}

	xtd := xt.TypeDescriptor()
	if isValidExtension(mr.Descriptor(), xtd) {
		mr.Clear(xtd)
	} else {
		mr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
			if int32(fd.Number()) == xt.Field {
				mr.Clear(fd)
				return false
			}
			return true
		})
	}
	clearUnknown(mr, fieldNum(xt.Field))
}

// ClearAllExtensions clears all extensions from m.
// This includes populated fields and unknown fields in the extension range.
func ClearAllExtensions(m Message) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() {
		return
	}

	mr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		if fd.IsExtension() {
			mr.Clear(fd)
		}
		return true
	})
	clearUnknown(mr, mr.Descriptor().ExtensionRanges())
>>>>>>> clientGRPCBilling
}

// GetExtension retrieves a proto2 extended field from pb.
//
// If the descriptor is type complete (i.e., ExtensionDesc.ExtensionType is non-nil),
// then GetExtension parses the encoded field and returns a Go value of the specified type.
// If the field is not present, then the default value is returned (if one is specified),
// otherwise ErrMissingExtension is reported.
//
<<<<<<< HEAD
// If the descriptor is not type complete (i.e., ExtensionDesc.ExtensionType is nil),
// then GetExtension returns the raw encoded bytes of the field extension.
func GetExtension(pb Message, extension *ExtensionDesc) (interface{}, error) {
	epb, err := extendable(pb)
	if err != nil {
		return nil, err
	}

	if extension.ExtendedType != nil {
		// can only check type if this is a complete descriptor
		if err := checkExtensionTypes(epb, extension); err != nil {
			return nil, err
		}
	}

	emap, mu := epb.extensionsRead()
	if emap == nil {
		return defaultExtensionValue(extension)
	}
	mu.Lock()
	defer mu.Unlock()
	e, ok := emap[extension.Field]
	if !ok {
		// defaultExtensionValue returns the default value or
		// ErrMissingExtension if there is no default.
		return defaultExtensionValue(extension)
	}

	if e.value != nil {
		// Already decoded. Check the descriptor, though.
		if e.desc != extension {
			// This shouldn't happen. If it does, it means that
			// GetExtension was called twice with two different
			// descriptors with the same field number.
			return nil, errors.New("proto: descriptor conflict")
		}
		return extensionAsLegacyType(e.value), nil
	}

	if extension.ExtensionType == nil {
		// incomplete descriptor
		return e.enc, nil
	}

	v, err := decodeExtension(e.enc, extension)
	if err != nil {
		return nil, err
	}

	// Remember the decoded version and drop the encoded version.
	// That way it is safe to mutate what we return.
	e.value = extensionAsStorageType(v)
	e.desc = extension
	e.enc = nil
	emap[extension.Field] = e
	return extensionAsLegacyType(e.value), nil
}

// defaultExtensionValue returns the default value for extension.
// If no default for an extension is defined ErrMissingExtension is returned.
func defaultExtensionValue(extension *ExtensionDesc) (interface{}, error) {
	if extension.ExtensionType == nil {
		// incomplete descriptor, so no default
		return nil, ErrMissingExtension
	}

	t := reflect.TypeOf(extension.ExtensionType)
	props := extensionProperties(extension)

	sf, _, err := fieldDefault(t, props)
	if err != nil {
		return nil, err
	}

	if sf == nil || sf.value == nil {
		// There is no default value.
		return nil, ErrMissingExtension
	}

	if t.Kind() != reflect.Ptr {
		// We do not need to return a Ptr, we can directly return sf.value.
		return sf.value, nil
	}

	// We need to return an interface{} that is a pointer to sf.value.
	value := reflect.New(t).Elem()
	value.Set(reflect.New(value.Type().Elem()))
	if sf.kind == reflect.Int32 {
		// We may have an int32 or an enum, but the underlying data is int32.
		// Since we can't set an int32 into a non int32 reflect.value directly
		// set it as a int32.
		value.Elem().SetInt(int64(sf.value.(int32)))
	} else {
		value.Elem().Set(reflect.ValueOf(sf.value))
	}
	return value.Interface(), nil
}

// decodeExtension decodes an extension encoded in b.
func decodeExtension(b []byte, extension *ExtensionDesc) (interface{}, error) {
	t := reflect.TypeOf(extension.ExtensionType)
	unmarshal := typeUnmarshaler(t, extension.Tag)

	// t is a pointer to a struct, pointer to basic type or a slice.
	// Allocate space to store the pointer/slice.
	value := reflect.New(t).Elem()

	var err error
	for {
		x, n := decodeVarint(b)
		if n == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		b = b[n:]
		wire := int(x) & 7

		b, err = unmarshal(b, valToPointer(value.Addr()), wire)
		if err != nil {
			return nil, err
		}

		if len(b) == 0 {
			break
		}
	}
	return value.Interface(), nil
}

// GetExtensions returns a slice of the extensions present in pb that are also listed in es.
// The returned slice has the same length as es; missing extensions will appear as nil elements.
func GetExtensions(pb Message, es []*ExtensionDesc) (extensions []interface{}, err error) {
	epb, err := extendable(pb)
	if err != nil {
		return nil, err
	}
	extensions = make([]interface{}, len(es))
	for i, e := range es {
		extensions[i], err = GetExtension(epb, e)
		if err == ErrMissingExtension {
			err = nil
		}
		if err != nil {
			return
		}
	}
	return
}

// ExtensionDescs returns a new slice containing pb's extension descriptors, in undefined order.
// For non-registered extensions, ExtensionDescs returns an incomplete descriptor containing
// just the Field field, which defines the extension's field number.
func ExtensionDescs(pb Message) ([]*ExtensionDesc, error) {
	epb, err := extendable(pb)
	if err != nil {
		return nil, err
	}
	registeredExtensions := RegisteredExtensions(pb)

	emap, mu := epb.extensionsRead()
	if emap == nil {
		return nil, nil
	}
	mu.Lock()
	defer mu.Unlock()
	extensions := make([]*ExtensionDesc, 0, len(emap))
	for extid, e := range emap {
		desc := e.desc
		if desc == nil {
			desc = registeredExtensions[extid]
			if desc == nil {
				desc = &ExtensionDesc{Field: extid}
			}
		}

		extensions = append(extensions, desc)
	}
	return extensions, nil
}

// SetExtension sets the specified extension of pb to the specified value.
func SetExtension(pb Message, extension *ExtensionDesc, value interface{}) error {
	epb, err := extendable(pb)
	if err != nil {
		return err
	}
	if err := checkExtensionTypes(epb, extension); err != nil {
		return err
	}
	typ := reflect.TypeOf(extension.ExtensionType)
	if typ != reflect.TypeOf(value) {
		return fmt.Errorf("proto: bad extension value type. got: %T, want: %T", value, extension.ExtensionType)
	}
	// nil extension values need to be caught early, because the
	// encoder can't distinguish an ErrNil due to a nil extension
	// from an ErrNil due to a missing field. Extensions are
	// always optional, so the encoder would just swallow the error
	// and drop all the extensions from the encoded message.
	if reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("proto: SetExtension called with nil value of type %T", value)
	}

	extmap := epb.extensionsWrite()
	extmap[extension.Field] = Extension{desc: extension, value: extensionAsStorageType(value)}
	return nil
}

// ClearAllExtensions clears all extensions from pb.
func ClearAllExtensions(pb Message) {
	epb, err := extendable(pb)
	if err != nil {
		return
	}
	m := epb.extensionsWrite()
	for k := range m {
		delete(m, k)
	}
}

// A global registry of extensions.
// The generated code will register the generated descriptors by calling RegisterExtension.

var extensionMaps = make(map[reflect.Type]map[int32]*ExtensionDesc)

// RegisterExtension is called from the generated code.
func RegisterExtension(desc *ExtensionDesc) {
	st := reflect.TypeOf(desc.ExtendedType).Elem()
	m := extensionMaps[st]
	if m == nil {
		m = make(map[int32]*ExtensionDesc)
		extensionMaps[st] = m
	}
	if _, ok := m[desc.Field]; ok {
		panic("proto: duplicate extension registered: " + st.String() + " " + strconv.Itoa(int(desc.Field)))
	}
	m[desc.Field] = desc
}

// RegisteredExtensions returns a map of the registered extensions of a
// protocol buffer struct, indexed by the extension number.
// The argument pb should be a nil pointer to the struct type.
func RegisteredExtensions(pb Message) map[int32]*ExtensionDesc {
	return extensionMaps[reflect.TypeOf(pb).Elem()]
}

// extensionAsLegacyType converts an value in the storage type as the API type.
// See Extension.value.
func extensionAsLegacyType(v interface{}) interface{} {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		// Represent primitive types as a pointer to the value.
		rv2 := reflect.New(rv.Type())
		rv2.Elem().Set(rv)
		v = rv2.Interface()
	case reflect.Ptr:
		// Represent slice types as the value itself.
		switch rv.Type().Elem().Kind() {
		case reflect.Slice:
			if rv.IsNil() {
				v = reflect.Zero(rv.Type().Elem()).Interface()
			} else {
				v = rv.Elem().Interface()
			}
		}
	}
	return v
}

// extensionAsStorageType converts an value in the API type as the storage type.
// See Extension.value.
func extensionAsStorageType(v interface{}) interface{} {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Ptr:
		// Represent slice types as the value itself.
		switch rv.Type().Elem().Kind() {
		case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
			if rv.IsNil() {
				v = reflect.Zero(rv.Type().Elem()).Interface()
			} else {
				v = rv.Elem().Interface()
			}
		}
	case reflect.Slice:
		// Represent slice types as a pointer to the value.
		if rv.Type().Elem().Kind() != reflect.Uint8 {
			rv2 := reflect.New(rv.Type())
			rv2.Elem().Set(rv)
			v = rv2.Interface()
		}
	}
	return v
=======
// If the descriptor is type incomplete (i.e., ExtensionDesc.ExtensionType is nil),
// then GetExtension returns the raw encoded bytes for the extension field.
func GetExtension(m Message, xt *ExtensionDesc) (interface{}, error) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() || mr.Descriptor().ExtensionRanges().Len() == 0 {
		return nil, errNotExtendable
	}

	// Retrieve the unknown fields for this extension field.
	var bo protoreflect.RawFields
	for bi := mr.GetUnknown(); len(bi) > 0; {
		num, _, n := protowire.ConsumeField(bi)
		if int32(num) == xt.Field {
			bo = append(bo, bi[:n]...)
		}
		bi = bi[n:]
	}

	// For type incomplete descriptors, only retrieve the unknown fields.
	if xt.ExtensionType == nil {
		return []byte(bo), nil
	}

	// If the extension field only exists as unknown fields, unmarshal it.
	// This is rarely done since proto.Unmarshal eagerly unmarshals extensions.
	xtd := xt.TypeDescriptor()
	if !isValidExtension(mr.Descriptor(), xtd) {
		return nil, fmt.Errorf("proto: bad extended type; %T does not extend %T", xt.ExtendedType, m)
	}
	if !mr.Has(xtd) && len(bo) > 0 {
		m2 := mr.New()
		if err := (proto.UnmarshalOptions{
			Resolver: extensionResolver{xt},
		}.Unmarshal(bo, m2.Interface())); err != nil {
			return nil, err
		}
		if m2.Has(xtd) {
			mr.Set(xtd, m2.Get(xtd))
			clearUnknown(mr, fieldNum(xt.Field))
		}
	}

	// Check whether the message has the extension field set or a default.
	var pv protoreflect.Value
	switch {
	case mr.Has(xtd):
		pv = mr.Get(xtd)
	case xtd.HasDefault():
		pv = xtd.Default()
	default:
		return nil, ErrMissingExtension
	}

	v := xt.InterfaceOf(pv)
	rv := reflect.ValueOf(v)
	if isScalarKind(rv.Kind()) {
		rv2 := reflect.New(rv.Type())
		rv2.Elem().Set(rv)
		v = rv2.Interface()
	}
	return v, nil
}

// extensionResolver is a custom extension resolver that stores a single
// extension type that takes precedence over the global registry.
type extensionResolver struct{ xt protoreflect.ExtensionType }

func (r extensionResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	if xtd := r.xt.TypeDescriptor(); xtd.FullName() == field {
		return r.xt, nil
	}
	return protoregistry.GlobalTypes.FindExtensionByName(field)
}

func (r extensionResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	if xtd := r.xt.TypeDescriptor(); xtd.ContainingMessage().FullName() == message && xtd.Number() == field {
		return r.xt, nil
	}
	return protoregistry.GlobalTypes.FindExtensionByNumber(message, field)
}

// GetExtensions returns a list of the extensions values present in m,
// corresponding with the provided list of extension descriptors, xts.
// If an extension is missing in m, the corresponding value is nil.
func GetExtensions(m Message, xts []*ExtensionDesc) ([]interface{}, error) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() {
		return nil, errNotExtendable
	}

	vs := make([]interface{}, len(xts))
	for i, xt := range xts {
		v, err := GetExtension(m, xt)
		if err != nil {
			if err == ErrMissingExtension {
				continue
			}
			return vs, err
		}
		vs[i] = v
	}
	return vs, nil
}

// SetExtension sets an extension field in m to the provided value.
func SetExtension(m Message, xt *ExtensionDesc, v interface{}) error {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() || mr.Descriptor().ExtensionRanges().Len() == 0 {
		return errNotExtendable
	}

	rv := reflect.ValueOf(v)
	if reflect.TypeOf(v) != reflect.TypeOf(xt.ExtensionType) {
		return fmt.Errorf("proto: bad extension value type. got: %T, want: %T", v, xt.ExtensionType)
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("proto: SetExtension called with nil value of type %T", v)
		}
		if isScalarKind(rv.Elem().Kind()) {
			v = rv.Elem().Interface()
		}
	}

	xtd := xt.TypeDescriptor()
	if !isValidExtension(mr.Descriptor(), xtd) {
		return fmt.Errorf("proto: bad extended type; %T does not extend %T", xt.ExtendedType, m)
	}
	mr.Set(xtd, xt.ValueOf(v))
	clearUnknown(mr, fieldNum(xt.Field))
	return nil
}

// SetRawExtension inserts b into the unknown fields of m.
//
// Deprecated: Use Message.ProtoReflect.SetUnknown instead.
func SetRawExtension(m Message, fnum int32, b []byte) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() {
		return
	}

	// Verify that the raw field is valid.
	for b0 := b; len(b0) > 0; {
		num, _, n := protowire.ConsumeField(b0)
		if int32(num) != fnum {
			panic(fmt.Sprintf("mismatching field number: got %d, want %d", num, fnum))
		}
		b0 = b0[n:]
	}

	ClearExtension(m, &ExtensionDesc{Field: fnum})
	mr.SetUnknown(append(mr.GetUnknown(), b...))
}

// ExtensionDescs returns a list of extension descriptors found in m,
// containing descriptors for both populated extension fields in m and
// also unknown fields of m that are in the extension range.
// For the later case, an type incomplete descriptor is provided where only
// the ExtensionDesc.Field field is populated.
// The order of the extension descriptors is undefined.
func ExtensionDescs(m Message) ([]*ExtensionDesc, error) {
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() || mr.Descriptor().ExtensionRanges().Len() == 0 {
		return nil, errNotExtendable
	}

	// Collect a set of known extension descriptors.
	extDescs := make(map[protoreflect.FieldNumber]*ExtensionDesc)
	mr.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.IsExtension() {
			xt := fd.(protoreflect.ExtensionTypeDescriptor)
			if xd, ok := xt.Type().(*ExtensionDesc); ok {
				extDescs[fd.Number()] = xd
			}
		}
		return true
	})

	// Collect a set of unknown extension descriptors.
	extRanges := mr.Descriptor().ExtensionRanges()
	for b := mr.GetUnknown(); len(b) > 0; {
		num, _, n := protowire.ConsumeField(b)
		if extRanges.Has(num) && extDescs[num] == nil {
			extDescs[num] = nil
		}
		b = b[n:]
	}

	// Transpose the set of descriptors into a list.
	var xts []*ExtensionDesc
	for num, xt := range extDescs {
		if xt == nil {
			xt = &ExtensionDesc{Field: int32(num)}
		}
		xts = append(xts, xt)
	}
	return xts, nil
}

// isValidExtension reports whether xtd is a valid extension descriptor for md.
func isValidExtension(md protoreflect.MessageDescriptor, xtd protoreflect.ExtensionTypeDescriptor) bool {
	return xtd.ContainingMessage() == md && md.ExtensionRanges().Has(xtd.Number())
}

// isScalarKind reports whether k is a protobuf scalar kind (except bytes).
// This function exists for historical reasons since the representation of
// scalars differs between v1 and v2, where v1 uses *T and v2 uses T.
func isScalarKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		return true
	default:
		return false
	}
}

// clearUnknown removes unknown fields from m where remover.Has reports true.
func clearUnknown(m protoreflect.Message, remover interface {
	Has(protoreflect.FieldNumber) bool
}) {
	var bo protoreflect.RawFields
	for bi := m.GetUnknown(); len(bi) > 0; {
		num, _, n := protowire.ConsumeField(bi)
		if !remover.Has(num) {
			bo = append(bo, bi[:n]...)
		}
		bi = bi[n:]
	}
	if bi := m.GetUnknown(); len(bi) != len(bo) {
		m.SetUnknown(bo)
	}
}

type fieldNum protoreflect.FieldNumber

func (n1 fieldNum) Has(n2 protoreflect.FieldNumber) bool {
	return protoreflect.FieldNumber(n1) == n2
>>>>>>> clientGRPCBilling
}
