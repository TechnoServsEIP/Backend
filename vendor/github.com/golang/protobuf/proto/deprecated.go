<<<<<<< HEAD
// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2018 The Go Authors.  All rights reserved.
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

import "errors"

// Deprecated: do not use.
type Stats struct{ Emalloc, Dmalloc, Encode, Decode, Chit, Cmiss, Size uint64 }

// Deprecated: do not use.
func GetStats() Stats { return Stats{} }

// Deprecated: do not use.
=======
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	// Deprecated: No longer returned.
	ErrNil = errors.New("proto: Marshal called with nil")

	// Deprecated: No longer returned.
	ErrTooLarge = errors.New("proto: message encodes to over 2 GB")

	// Deprecated: No longer returned.
	ErrInternalBadWireType = errors.New("proto: internal error: bad wiretype for oneof")
)

// Deprecated: Do not use.
type Stats struct{ Emalloc, Dmalloc, Encode, Decode, Chit, Cmiss, Size uint64 }

// Deprecated: Do not use.
func GetStats() Stats { return Stats{} }

// Deprecated: Do not use.
>>>>>>> clientGRPCBilling
func MarshalMessageSet(interface{}) ([]byte, error) {
	return nil, errors.New("proto: not implemented")
}

<<<<<<< HEAD
// Deprecated: do not use.
=======
// Deprecated: Do not use.
>>>>>>> clientGRPCBilling
func UnmarshalMessageSet([]byte, interface{}) error {
	return errors.New("proto: not implemented")
}

<<<<<<< HEAD
// Deprecated: do not use.
=======
// Deprecated: Do not use.
>>>>>>> clientGRPCBilling
func MarshalMessageSetJSON(interface{}) ([]byte, error) {
	return nil, errors.New("proto: not implemented")
}

<<<<<<< HEAD
// Deprecated: do not use.
=======
// Deprecated: Do not use.
>>>>>>> clientGRPCBilling
func UnmarshalMessageSetJSON([]byte, interface{}) error {
	return errors.New("proto: not implemented")
}

<<<<<<< HEAD
// Deprecated: do not use.
func RegisterMessageSetType(Message, int32, string) {}
=======
// Deprecated: Do not use.
func RegisterMessageSetType(Message, int32, string) {}

// Deprecated: Do not use.
func EnumName(m map[int32]string, v int32) string {
	s, ok := m[v]
	if ok {
		return s
	}
	return strconv.Itoa(int(v))
}

// Deprecated: Do not use.
func UnmarshalJSONEnum(m map[string]int32, data []byte, enumName string) (int32, error) {
	if data[0] == '"' {
		// New style: enums are strings.
		var repr string
		if err := json.Unmarshal(data, &repr); err != nil {
			return -1, err
		}
		val, ok := m[repr]
		if !ok {
			return 0, fmt.Errorf("unrecognized enum %s value %q", enumName, repr)
		}
		return val, nil
	}
	// Old style: enums are ints.
	var val int32
	if err := json.Unmarshal(data, &val); err != nil {
		return 0, fmt.Errorf("cannot unmarshal %#q into enum %s", data, enumName)
	}
	return val, nil
}

// Deprecated: Do not use.
type InternalMessageInfo struct{}

func (*InternalMessageInfo) DiscardUnknown(Message)                        { panic("not implemented") }
func (*InternalMessageInfo) Marshal([]byte, Message, bool) ([]byte, error) { panic("not implemented") }
func (*InternalMessageInfo) Merge(Message, Message)                        { panic("not implemented") }
func (*InternalMessageInfo) Size(Message) int                              { panic("not implemented") }
func (*InternalMessageInfo) Unmarshal(Message, []byte) error               { panic("not implemented") }
>>>>>>> clientGRPCBilling
