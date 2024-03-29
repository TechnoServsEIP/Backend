// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server_action.proto

package game_servers

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Start Server Action
type StartServerActionRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartServerActionRequest) Reset()         { *m = StartServerActionRequest{} }
func (m *StartServerActionRequest) String() string { return proto.CompactTextString(m) }
func (*StartServerActionRequest) ProtoMessage()    {}
func (*StartServerActionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d14b02b978e672f, []int{0}
}

func (m *StartServerActionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartServerActionRequest.Unmarshal(m, b)
}
func (m *StartServerActionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartServerActionRequest.Marshal(b, m, deterministic)
}
func (m *StartServerActionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartServerActionRequest.Merge(m, src)
}
func (m *StartServerActionRequest) XXX_Size() int {
	return xxx_messageInfo_StartServerActionRequest.Size(m)
}
func (m *StartServerActionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartServerActionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartServerActionRequest proto.InternalMessageInfo

func (m *StartServerActionRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type StartServerActionResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartServerActionResponse) Reset()         { *m = StartServerActionResponse{} }
func (m *StartServerActionResponse) String() string { return proto.CompactTextString(m) }
func (*StartServerActionResponse) ProtoMessage()    {}
func (*StartServerActionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d14b02b978e672f, []int{1}
}

func (m *StartServerActionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartServerActionResponse.Unmarshal(m, b)
}
func (m *StartServerActionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartServerActionResponse.Marshal(b, m, deterministic)
}
func (m *StartServerActionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartServerActionResponse.Merge(m, src)
}
func (m *StartServerActionResponse) XXX_Size() int {
	return xxx_messageInfo_StartServerActionResponse.Size(m)
}
func (m *StartServerActionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartServerActionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartServerActionResponse proto.InternalMessageInfo

// Stop Server Action
type StopServerActionRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopServerActionRequest) Reset()         { *m = StopServerActionRequest{} }
func (m *StopServerActionRequest) String() string { return proto.CompactTextString(m) }
func (*StopServerActionRequest) ProtoMessage()    {}
func (*StopServerActionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d14b02b978e672f, []int{2}
}

func (m *StopServerActionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopServerActionRequest.Unmarshal(m, b)
}
func (m *StopServerActionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopServerActionRequest.Marshal(b, m, deterministic)
}
func (m *StopServerActionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopServerActionRequest.Merge(m, src)
}
func (m *StopServerActionRequest) XXX_Size() int {
	return xxx_messageInfo_StopServerActionRequest.Size(m)
}
func (m *StopServerActionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopServerActionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopServerActionRequest proto.InternalMessageInfo

func (m *StopServerActionRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type StopServerActionResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopServerActionResponse) Reset()         { *m = StopServerActionResponse{} }
func (m *StopServerActionResponse) String() string { return proto.CompactTextString(m) }
func (*StopServerActionResponse) ProtoMessage()    {}
func (*StopServerActionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d14b02b978e672f, []int{3}
}

func (m *StopServerActionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopServerActionResponse.Unmarshal(m, b)
}
func (m *StopServerActionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopServerActionResponse.Marshal(b, m, deterministic)
}
func (m *StopServerActionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopServerActionResponse.Merge(m, src)
}
func (m *StopServerActionResponse) XXX_Size() int {
	return xxx_messageInfo_StopServerActionResponse.Size(m)
}
func (m *StopServerActionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StopServerActionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StopServerActionResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*StartServerActionRequest)(nil), "technoservs.game_servers.StartServerActionRequest")
	proto.RegisterType((*StartServerActionResponse)(nil), "technoservs.game_servers.StartServerActionResponse")
	proto.RegisterType((*StopServerActionRequest)(nil), "technoservs.game_servers.StopServerActionRequest")
	proto.RegisterType((*StopServerActionResponse)(nil), "technoservs.game_servers.StopServerActionResponse")
}

func init() { proto.RegisterFile("server_action.proto", fileDescriptor_8d14b02b978e672f) }

var fileDescriptor_8d14b02b978e672f = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0x8a, 0x4f, 0x4c, 0x2e, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x92, 0x28, 0x49, 0x4d, 0xce, 0xc8, 0xcb, 0x07, 0x49, 0x15, 0xeb, 0xa5, 0x27, 0xe6, 0xa6, 0xc6,
	0x43, 0x54, 0x15, 0x2b, 0xe9, 0x71, 0x49, 0x04, 0x97, 0x24, 0x16, 0x95, 0x04, 0x83, 0xf9, 0x8e,
	0x60, 0x4d, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42, 0x5c, 0x2c, 0xa5, 0xa5, 0x99,
	0x29, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92, 0x34, 0x97, 0x24, 0x16, 0xf5,
	0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x4a, 0xba, 0x5c, 0xe2, 0xc1, 0x25, 0xf9, 0x05, 0xc4, 0x9a,
	0x25, 0x05, 0xb2, 0x1b, 0x5d, 0x39, 0xc4, 0xa8, 0x24, 0x36, 0xb0, 0xc3, 0x8d, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xba, 0x7f, 0x9e, 0x81, 0xcf, 0x00, 0x00, 0x00,
}
