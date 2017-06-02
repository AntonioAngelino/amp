// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/tests/integration/etcd/store_test.proto
// DO NOT EDIT!

/*
Package etcd is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/tests/integration/etcd/store_test.proto

It has these top-level messages:
	TestMessage
*/
package etcd

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TestMessage struct {
	Id    string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name  string   `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Names []string `protobuf:"bytes,3,rep,name=names" json:"names,omitempty"`
}

func (m *TestMessage) Reset()                    { *m = TestMessage{} }
func (m *TestMessage) String() string            { return proto.CompactTextString(m) }
func (*TestMessage) ProtoMessage()               {}
func (*TestMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TestMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TestMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TestMessage) GetNames() []string {
	if m != nil {
		return m.Names
	}
	return nil
}

func init() {
	proto.RegisterType((*TestMessage)(nil), "etcd.TestMessage")
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/tests/integration/etcd/store_test.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 157 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x24, 0x8c, 0xbb, 0xaa, 0xc3, 0x30,
	0x10, 0x44, 0xf1, 0xe3, 0x5e, 0xb0, 0x02, 0x29, 0x44, 0x0a, 0x97, 0x26, 0x95, 0x2b, 0x6f, 0x91,
	0x4f, 0x48, 0x91, 0x2a, 0x8d, 0x49, 0x1f, 0xd6, 0xf6, 0xe2, 0x08, 0x22, 0xad, 0xd0, 0x6e, 0xfe,
	0x3f, 0x48, 0xa9, 0x66, 0xe6, 0x0c, 0x1c, 0x73, 0xdd, 0x9d, 0xbe, 0x3e, 0xcb, 0xb4, 0xb2, 0x07,
	0x8c, 0x71, 0xa5, 0x37, 0x25, 0x54, 0x4e, 0x80, 0x3e, 0x82, 0x92, 0xa8, 0x80, 0x0b, 0x4a, 0x7b,
	0x42, 0x75, 0x1c, 0x80, 0x74, 0xdd, 0x40, 0x94, 0x13, 0x3d, 0xf3, 0x39, 0xc5, 0xc4, 0xca, 0xb6,
	0xcd, 0xf8, 0x7c, 0x33, 0x87, 0x07, 0x89, 0xde, 0x49, 0x04, 0x77, 0xb2, 0x47, 0x53, 0xbb, 0xad,
	0xaf, 0x86, 0x6a, 0xec, 0xe6, 0xda, 0x6d, 0xd6, 0x9a, 0x36, 0xa0, 0xa7, 0xbe, 0x2e, 0xa4, 0x74,
	0x7b, 0x32, 0x7f, 0x39, 0xa5, 0x6f, 0x86, 0x66, 0xec, 0xe6, 0xdf, 0x58, 0xfe, 0x8b, 0xf5, 0xf2,
	0x0d, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x73, 0xe5, 0x72, 0x9c, 0x00, 0x00, 0x00,
}