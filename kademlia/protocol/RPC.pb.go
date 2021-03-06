// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Protocol files/RPC.proto

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

// Specifies the message type contained in message
type RPCMessageTYPE int32

const (
	RPC_PING         RPCMessageTYPE = 0
	RPC_STORE        RPCMessageTYPE = 1
	RPC_FIND_NODE    RPCMessageTYPE = 2
	RPC_FIND_VALUE   RPCMessageTYPE = 3
	RPC_PIN          RPCMessageTYPE = 4
	RPC_UNPIN        RPCMessageTYPE = 5
	RPC_STORE_ANSWER RPCMessageTYPE = 6
	RPC_SEND_FILE    RPCMessageTYPE = 7
)

var RPCMessageTYPE_name = map[int32]string{
	0: "PING",
	1: "STORE",
	2: "FIND_NODE",
	3: "FIND_VALUE",
	4: "PIN",
	5: "UNPIN",
	6: "STORE_ANSWER",
	7: "SEND_FILE",
}

var RPCMessageTYPE_value = map[string]int32{
	"PING":         0,
	"STORE":        1,
	"FIND_NODE":    2,
	"FIND_VALUE":   3,
	"PIN":          4,
	"UNPIN":        5,
	"STORE_ANSWER": 6,
	"SEND_FILE":    7,
}

func (x RPCMessageTYPE) String() string {
	return proto.EnumName(RPCMessageTYPE_name, int32(x))
}

func (RPCMessageTYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_638d2753ed677b27, []int{0, 0}
}

type RPC struct {
	// Unique 160-bit message ID used for routing purposes
	MessageID   []byte         `protobuf:"bytes,1,opt,name=messageID,proto3" json:"messageID,omitempty"`
	MessageType RPCMessageTYPE `protobuf:"varint,2,opt,name=messageType,proto3,enum=protocol.RPCMessageTYPE" json:"messageType,omitempty"`
	// Contains the serialized message
	Message              []byte   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	IPaddress            string   `protobuf:"bytes,4,opt,name=IPaddress,proto3" json:"IPaddress,omitempty"`
	OriginalSender       []byte   `protobuf:"bytes,5,opt,name=originalSender,proto3" json:"originalSender,omitempty"`
	KademliaID           []byte   `protobuf:"bytes,6,opt,name=KademliaID,proto3" json:"KademliaID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RPC) Reset()         { *m = RPC{} }
func (m *RPC) String() string { return proto.CompactTextString(m) }
func (*RPC) ProtoMessage()    {}
func (*RPC) Descriptor() ([]byte, []int) {
	return fileDescriptor_638d2753ed677b27, []int{0}
}

func (m *RPC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPC.Unmarshal(m, b)
}
func (m *RPC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPC.Marshal(b, m, deterministic)
}
func (m *RPC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPC.Merge(m, src)
}
func (m *RPC) XXX_Size() int {
	return xxx_messageInfo_RPC.Size(m)
}
func (m *RPC) XXX_DiscardUnknown() {
	xxx_messageInfo_RPC.DiscardUnknown(m)
}

var xxx_messageInfo_RPC proto.InternalMessageInfo

func (m *RPC) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *RPC) GetMessageType() RPCMessageTYPE {
	if m != nil {
		return m.MessageType
	}
	return RPC_PING
}

func (m *RPC) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *RPC) GetIPaddress() string {
	if m != nil {
		return m.IPaddress
	}
	return ""
}

func (m *RPC) GetOriginalSender() []byte {
	if m != nil {
		return m.OriginalSender
	}
	return nil
}

func (m *RPC) GetKademliaID() []byte {
	if m != nil {
		return m.KademliaID
	}
	return nil
}

func init() {
	proto.RegisterEnum("protocol.RPCMessageTYPE", RPCMessageTYPE_name, RPCMessageTYPE_value)
	proto.RegisterType((*RPC)(nil), "protocol.RPC")
}

func init() { proto.RegisterFile("Protocol files/RPC.proto", fileDescriptor_638d2753ed677b27) }

var fileDescriptor_638d2753ed677b27 = []byte{
	// 285 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x5f, 0x4b, 0xc3, 0x30,
	0x14, 0xc5, 0xed, 0xba, 0x7f, 0xbd, 0xce, 0x11, 0xee, 0x53, 0x04, 0x91, 0xb2, 0x07, 0xd9, 0x53,
	0x05, 0x7d, 0xf4, 0x69, 0xac, 0x99, 0x04, 0x47, 0x16, 0xd2, 0x4d, 0xd9, 0x53, 0xa9, 0x36, 0x8e,
	0x42, 0xb7, 0x96, 0x56, 0x04, 0xbf, 0xa9, 0x1f, 0x47, 0x12, 0x5a, 0x3a, 0x7c, 0xcb, 0xf9, 0xe5,
	0x9c, 0x73, 0x2f, 0x17, 0xa8, 0xac, 0x8a, 0xaf, 0xe2, 0xa3, 0xc8, 0xfd, 0xcf, 0x2c, 0xd7, 0xf5,
	0xbd, 0x92, 0xcb, 0xa0, 0x34, 0x08, 0xc7, 0x65, 0xf3, 0x33, 0xfb, 0xed, 0x81, 0xab, 0xe4, 0x12,
	0x6f, 0xc0, 0x3b, 0xea, 0xba, 0x4e, 0x0e, 0x9a, 0x87, 0xd4, 0xf1, 0x9d, 0xf9, 0x44, 0x75, 0x00,
	0x9f, 0xe0, 0xb2, 0x11, 0xdb, 0x9f, 0x52, 0xd3, 0x9e, 0xef, 0xcc, 0xa7, 0x0f, 0xd7, 0x41, 0xdb,
	0x12, 0x98, 0xe6, 0xd6, 0xb0, 0x97, 0x4c, 0x9d, 0xbb, 0x91, 0xc2, 0xa8, 0x91, 0xd4, 0xb5, 0xc5,
	0xad, 0x34, 0x43, 0xb9, 0x4c, 0xd2, 0xb4, 0xd2, 0x75, 0x4d, 0xfb, 0xbe, 0x33, 0xf7, 0x54, 0x07,
	0xf0, 0x0e, 0xa6, 0x45, 0x95, 0x1d, 0xb2, 0x53, 0x92, 0x47, 0xfa, 0x94, 0xea, 0x8a, 0x0e, 0x6c,
	0xfc, 0x1f, 0xc5, 0x5b, 0x80, 0x97, 0x24, 0xd5, 0xc7, 0x3c, 0x4b, 0x78, 0x48, 0x87, 0xd6, 0x73,
	0x46, 0x66, 0xdf, 0xdd, 0xf2, 0x7b, 0xc9, 0x70, 0x0c, 0x7d, 0xc9, 0xc5, 0x33, 0xb9, 0x40, 0x0f,
	0x06, 0xd1, 0x76, 0xa3, 0x18, 0x71, 0xf0, 0x0a, 0xbc, 0x15, 0x17, 0x61, 0x2c, 0x36, 0x21, 0x23,
	0x3d, 0x9c, 0x02, 0x58, 0xf9, 0xba, 0x58, 0xef, 0x18, 0x71, 0x71, 0x04, 0xae, 0xe4, 0x82, 0xf4,
	0x4d, 0x64, 0x27, 0xcc, 0x73, 0x80, 0x04, 0x26, 0x36, 0x1d, 0x2f, 0x44, 0xf4, 0xc6, 0x14, 0x19,
	0x9a, 0x92, 0x88, 0x89, 0x30, 0x5e, 0xf1, 0x35, 0x23, 0xa3, 0xf7, 0xa1, 0x3d, 0xcf, 0xe3, 0x5f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x59, 0x88, 0x28, 0x87, 0x01, 0x00, 0x00,
}
