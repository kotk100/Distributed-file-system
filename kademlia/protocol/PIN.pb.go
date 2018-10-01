// Code generated by protoc-gen-go. DO NOT EDIT.
// source: PIN.proto

package protocol

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Pin struct {
	FileHash             []byte   `protobuf:"bytes,1,opt,name=fileHash,proto3" json:"fileHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pin) Reset()         { *m = Pin{} }
func (m *Pin) String() string { return proto.CompactTextString(m) }
func (*Pin) ProtoMessage()    {}
func (*Pin) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddcfcb921b582cdc, []int{0}
}

func (m *Pin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pin.Unmarshal(m, b)
}
func (m *Pin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pin.Marshal(b, m, deterministic)
}
func (m *Pin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pin.Merge(m, src)
}
func (m *Pin) XXX_Size() int {
	return xxx_messageInfo_Pin.Size(m)
}
func (m *Pin) XXX_DiscardUnknown() {
	xxx_messageInfo_Pin.DiscardUnknown(m)
}

var xxx_messageInfo_Pin proto.InternalMessageInfo

func (m *Pin) GetFileHash() []byte {
	if m != nil {
		return m.FileHash
	}
	return nil
}

func init() {
	proto.RegisterType((*Pin)(nil), "protocol.Pin")
}

func init() { proto.RegisterFile("PIN.proto", fileDescriptor_ddcfcb921b582cdc) }

var fileDescriptor_ddcfcb921b582cdc = []byte{
	// 74 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x0c, 0xf0, 0xf4, 0xd3,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0xc9, 0xf9, 0x39, 0x4a, 0x8a, 0x5c, 0xcc,
	0x01, 0x99, 0x79, 0x42, 0x52, 0x5c, 0x1c, 0x69, 0x99, 0x39, 0xa9, 0x1e, 0x89, 0xc5, 0x19, 0x12,
	0x8c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x70, 0x7e, 0x12, 0x1b, 0x58, 0xb1, 0x31, 0x20, 0x00, 0x00,
	0xff, 0xff, 0xc8, 0xca, 0x3f, 0x83, 0x40, 0x00, 0x00, 0x00,
}