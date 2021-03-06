// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hellogrpc.proto

package com_salty_protos

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// defined Req message struct
type Req struct {
	// type
	JsonStr              string   `protobuf:"bytes,1,opt,name=jsonStr,proto3" json:"jsonStr,omitempty"`
	Data                 *any.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Req) Reset()         { *m = Req{} }
func (m *Req) String() string { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()    {}
func (*Req) Descriptor() ([]byte, []int) {
	return fileDescriptor_94272a0853831a50, []int{0}
}

func (m *Req) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Req.Unmarshal(m, b)
}
func (m *Req) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Req.Marshal(b, m, deterministic)
}
func (m *Req) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Req.Merge(m, src)
}
func (m *Req) XXX_Size() int {
	return xxx_messageInfo_Req.Size(m)
}
func (m *Req) XXX_DiscardUnknown() {
	xxx_messageInfo_Req.DiscardUnknown(m)
}

var xxx_messageInfo_Req proto.InternalMessageInfo

func (m *Req) GetJsonStr() string {
	if m != nil {
		return m.JsonStr
	}
	return ""
}

func (m *Req) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

//defined Res message struct
type Res struct {
	BackJson             string   `protobuf:"bytes,1,opt,name=backJson,proto3" json:"backJson,omitempty"`
	Data                 *any.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Res) Reset()         { *m = Res{} }
func (m *Res) String() string { return proto.CompactTextString(m) }
func (*Res) ProtoMessage()    {}
func (*Res) Descriptor() ([]byte, []int) {
	return fileDescriptor_94272a0853831a50, []int{1}
}

func (m *Res) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Res.Unmarshal(m, b)
}
func (m *Res) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Res.Marshal(b, m, deterministic)
}
func (m *Res) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Res.Merge(m, src)
}
func (m *Res) XXX_Size() int {
	return xxx_messageInfo_Res.Size(m)
}
func (m *Res) XXX_DiscardUnknown() {
	xxx_messageInfo_Res.DiscardUnknown(m)
}

var xxx_messageInfo_Res proto.InternalMessageInfo

func (m *Res) GetBackJson() string {
	if m != nil {
		return m.BackJson
	}
	return ""
}

func (m *Res) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Req)(nil), "com.salty.protos.Req")
	proto.RegisterType((*Res)(nil), "com.salty.protos.Res")
}

func init() { proto.RegisterFile("hellogrpc.proto", fileDescriptor_94272a0853831a50) }

var fileDescriptor_94272a0853831a50 = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0x48, 0xcd, 0xc9,
	0xc9, 0x4f, 0x2f, 0x2a, 0x48, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x48, 0xce, 0xcf,
	0xd5, 0x2b, 0x4e, 0xcc, 0x29, 0xa9, 0x84, 0x08, 0x14, 0x4b, 0x49, 0xa6, 0xe7, 0xe7, 0xa7, 0xe7,
	0xa4, 0xea, 0x83, 0xb9, 0x49, 0xa5, 0x69, 0xfa, 0x89, 0x79, 0x50, 0x39, 0x25, 0x4f, 0x2e, 0xe6,
	0xa0, 0xd4, 0x42, 0x21, 0x09, 0x2e, 0xf6, 0xac, 0xe2, 0xfc, 0xbc, 0xe0, 0x92, 0x22, 0x09, 0x46,
	0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x57, 0x48, 0x83, 0x8b, 0x25, 0x25, 0xb1, 0x24, 0x51, 0x82,
	0x49, 0x81, 0x51, 0x83, 0xdb, 0x48, 0x44, 0x0f, 0x62, 0x94, 0x1e, 0xcc, 0x28, 0x3d, 0xc7, 0xbc,
	0xca, 0x20, 0xb0, 0x0a, 0x25, 0x6f, 0x90, 0x51, 0xc5, 0x42, 0x52, 0x5c, 0x1c, 0x49, 0x89, 0xc9,
	0xd9, 0x5e, 0xc5, 0xf9, 0x79, 0x50, 0xb3, 0xe0, 0x7c, 0xe2, 0x0d, 0x33, 0x72, 0xe4, 0x62, 0x0b,
	0x4f, 0xcc, 0x2c, 0x49, 0x2d, 0x12, 0x32, 0xe7, 0x62, 0x75, 0xc9, 0xf7, 0x75, 0x31, 0x15, 0x12,
	0xd5, 0x43, 0xf7, 0x98, 0x5e, 0x50, 0x6a, 0xa1, 0x14, 0x56, 0xe1, 0x62, 0x25, 0x86, 0x24, 0x36,
	0x30, 0xcf, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x41, 0x81, 0xa7, 0x12, 0x21, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WaiterClient is the client API for Waiter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WaiterClient interface {
	// defined interface (struct can be reuse)
	// methods
	DoMD5(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Res, error)
}

type waiterClient struct {
	cc *grpc.ClientConn
}

func NewWaiterClient(cc *grpc.ClientConn) WaiterClient {
	return &waiterClient{cc}
}

func (c *waiterClient) DoMD5(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Res, error) {
	out := new(Res)
	err := c.cc.Invoke(ctx, "/com.salty.protos.Waiter/DoMD5", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WaiterServer is the server API for Waiter service.
type WaiterServer interface {
	// defined interface (struct can be reuse)
	// methods
	DoMD5(context.Context, *Req) (*Res, error)
}

// UnimplementedWaiterServer can be embedded to have forward compatible implementations.
type UnimplementedWaiterServer struct {
}

func (*UnimplementedWaiterServer) DoMD5(ctx context.Context, req *Req) (*Res, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoMD5 not implemented")
}

func RegisterWaiterServer(s *grpc.Server, srv WaiterServer) {
	s.RegisterService(&_Waiter_serviceDesc, srv)
}

func _Waiter_DoMD5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WaiterServer).DoMD5(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.salty.protos.Waiter/DoMD5",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WaiterServer).DoMD5(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

var _Waiter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.salty.protos.Waiter",
	HandlerType: (*WaiterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoMD5",
			Handler:    _Waiter_DoMD5_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hellogrpc.proto",
}
