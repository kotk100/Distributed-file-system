// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Protocol files/SEND_FILE.proto

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

type SendFile struct {
	IPaddress            string   `protobuf:"bytes,1,opt,name=IPaddress,proto3" json:"IPaddress,omitempty"`
	OriginalSender       []byte   `protobuf:"bytes,2,opt,name=originalSender,proto3" json:"originalSender,omitempty"`
	KademliaID           []byte   `protobuf:"bytes,3,opt,name=KademliaID,proto3" json:"KademliaID,omitempty"`
	FileHash             []byte   `protobuf:"bytes,4,opt,name=fileHash,proto3" json:"fileHash,omitempty"`
	FileSize             int64    `protobuf:"varint,5,opt,name=fileSize,proto3" json:"fileSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendFile) Reset()         { *m = SendFile{} }
func (m *SendFile) String() string { return proto.CompactTextString(m) }
func (*SendFile) ProtoMessage()    {}
func (*SendFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_17ee577262736f51, []int{0}
}

func (m *SendFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendFile.Unmarshal(m, b)
}
func (m *SendFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendFile.Marshal(b, m, deterministic)
}
func (m *SendFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendFile.Merge(m, src)
}
func (m *SendFile) XXX_Size() int {
	return xxx_messageInfo_SendFile.Size(m)
}
func (m *SendFile) XXX_DiscardUnknown() {
	xxx_messageInfo_SendFile.DiscardUnknown(m)
}

var xxx_messageInfo_SendFile proto.InternalMessageInfo

func (m *SendFile) GetIPaddress() string {
	if m != nil {
		return m.IPaddress
	}
	return ""
}

func (m *SendFile) GetOriginalSender() []byte {
	if m != nil {
		return m.OriginalSender
	}
	return nil
}

func (m *SendFile) GetKademliaID() []byte {
	if m != nil {
		return m.KademliaID
	}
	return nil
}

func (m *SendFile) GetFileHash() []byte {
	if m != nil {
		return m.FileHash
	}
	return nil
}

func (m *SendFile) GetFileSize() int64 {
	if m != nil {
		return m.FileSize
	}
	return 0
}

func init() {
	proto.RegisterType((*SendFile)(nil), "protocol.SendFile")
}

func init() { proto.RegisterFile("Protocol files/SEND_FILE.proto", fileDescriptor_17ee577262736f51) }

var fileDescriptor_17ee577262736f51 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x0b, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0xce, 0xcf, 0x51, 0x48, 0xcb, 0xcc, 0x49, 0x2d, 0xd6, 0x0f, 0x76, 0xf5, 0x73, 0x89,
	0x77, 0xf3, 0xf4, 0x71, 0xd5, 0x2b, 0x00, 0x49, 0x08, 0x71, 0x14, 0x40, 0xe5, 0x95, 0x56, 0x30,
	0x72, 0x71, 0x04, 0xa7, 0xe6, 0xa5, 0xb8, 0x65, 0xe6, 0xa4, 0x0a, 0xc9, 0x70, 0x71, 0x7a, 0x06,
	0x24, 0xa6, 0xa4, 0x14, 0xa5, 0x16, 0x17, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x21, 0x04,
	0x84, 0xd4, 0xb8, 0xf8, 0xf2, 0x8b, 0x32, 0xd3, 0x33, 0xf3, 0x12, 0x73, 0x40, 0x3a, 0x52, 0x8b,
	0x24, 0x98, 0x14, 0x18, 0x35, 0x78, 0x82, 0xd0, 0x44, 0x85, 0xe4, 0xb8, 0xb8, 0xbc, 0x13, 0x53,
	0x52, 0x73, 0x73, 0x32, 0x13, 0x3d, 0x5d, 0x24, 0x98, 0xc1, 0x6a, 0x90, 0x44, 0x84, 0xa4, 0xb8,
	0x38, 0x40, 0xae, 0xf2, 0x48, 0x2c, 0xce, 0x90, 0x60, 0x01, 0xcb, 0xc2, 0xf9, 0x30, 0xb9, 0xe0,
	0xcc, 0xaa, 0x54, 0x09, 0x56, 0x05, 0x46, 0x0d, 0xe6, 0x20, 0x38, 0x3f, 0x89, 0x0d, 0xec, 0x68,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x29, 0x08, 0xc4, 0x29, 0xdd, 0x00, 0x00, 0x00,
}