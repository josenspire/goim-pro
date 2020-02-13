// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sms.proto

package com_salty_protos

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type ObtainSMSCodeReq_CodeType int32

const (
	ObtainSMSCodeReq_REGISTER      ObtainSMSCodeReq_CodeType = 0
	ObtainSMSCodeReq_LOGIN         ObtainSMSCodeReq_CodeType = 1
	ObtainSMSCodeReq_ResetPassword ObtainSMSCodeReq_CodeType = 2
)

var ObtainSMSCodeReq_CodeType_name = map[int32]string{
	0: "REGISTER",
	1: "LOGIN",
	2: "ResetPassword",
}

var ObtainSMSCodeReq_CodeType_value = map[string]int32{
	"REGISTER":      0,
	"LOGIN":         1,
	"ResetPassword": 2,
}

func (x ObtainSMSCodeReq_CodeType) String() string {
	return proto.EnumName(ObtainSMSCodeReq_CodeType_name, int32(x))
}

func (ObtainSMSCodeReq_CodeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c8d8bdc537111860, []int{0, 0}
}

//sms服务就是短信服务，邮件相关的接口要独立出来
type ObtainSMSCodeReq struct {
	CodeType             ObtainSMSCodeReq_CodeType `protobuf:"varint,1,opt,name=codeType,proto3,enum=com.salty.protos.ObtainSMSCodeReq_CodeType" json:"codeType,omitempty"`
	Telephone            string                    `protobuf:"bytes,2,opt,name=telephone,proto3" json:"telephone,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ObtainSMSCodeReq) Reset()         { *m = ObtainSMSCodeReq{} }
func (m *ObtainSMSCodeReq) String() string { return proto.CompactTextString(m) }
func (*ObtainSMSCodeReq) ProtoMessage()    {}
func (*ObtainSMSCodeReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8d8bdc537111860, []int{0}
}

func (m *ObtainSMSCodeReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ObtainSMSCodeReq.Unmarshal(m, b)
}
func (m *ObtainSMSCodeReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ObtainSMSCodeReq.Marshal(b, m, deterministic)
}
func (m *ObtainSMSCodeReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObtainSMSCodeReq.Merge(m, src)
}
func (m *ObtainSMSCodeReq) XXX_Size() int {
	return xxx_messageInfo_ObtainSMSCodeReq.Size(m)
}
func (m *ObtainSMSCodeReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ObtainSMSCodeReq.DiscardUnknown(m)
}

var xxx_messageInfo_ObtainSMSCodeReq proto.InternalMessageInfo

func (m *ObtainSMSCodeReq) GetCodeType() ObtainSMSCodeReq_CodeType {
	if m != nil {
		return m.CodeType
	}
	return ObtainSMSCodeReq_REGISTER
}

func (m *ObtainSMSCodeReq) GetTelephone() string {
	if m != nil {
		return m.Telephone
	}
	return ""
}

type ObtainSMSCodeResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ObtainSMSCodeResp) Reset()         { *m = ObtainSMSCodeResp{} }
func (m *ObtainSMSCodeResp) String() string { return proto.CompactTextString(m) }
func (*ObtainSMSCodeResp) ProtoMessage()    {}
func (*ObtainSMSCodeResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8d8bdc537111860, []int{1}
}

func (m *ObtainSMSCodeResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ObtainSMSCodeResp.Unmarshal(m, b)
}
func (m *ObtainSMSCodeResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ObtainSMSCodeResp.Marshal(b, m, deterministic)
}
func (m *ObtainSMSCodeResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObtainSMSCodeResp.Merge(m, src)
}
func (m *ObtainSMSCodeResp) XXX_Size() int {
	return xxx_messageInfo_ObtainSMSCodeResp.Size(m)
}
func (m *ObtainSMSCodeResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ObtainSMSCodeResp.DiscardUnknown(m)
}

var xxx_messageInfo_ObtainSMSCodeResp proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("com.salty.protos.ObtainSMSCodeReq_CodeType", ObtainSMSCodeReq_CodeType_name, ObtainSMSCodeReq_CodeType_value)
	proto.RegisterType((*ObtainSMSCodeReq)(nil), "com.salty.protos.ObtainSMSCodeReq")
	proto.RegisterType((*ObtainSMSCodeResp)(nil), "com.salty.protos.ObtainSMSCodeResp")
}

func init() { proto.RegisterFile("sms.proto", fileDescriptor_c8d8bdc537111860) }

var fileDescriptor_c8d8bdc537111860 = []byte{
	// 237 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0xbb, 0x01, 0x25, 0x19, 0x5a, 0xd9, 0x8e, 0x97, 0x1a, 0x3c, 0x94, 0x9c, 0x0a, 0xc2,
	0x1e, 0x2a, 0xf8, 0x00, 0x96, 0x12, 0x0b, 0xd6, 0x96, 0xdd, 0xe2, 0x3d, 0x4d, 0x06, 0x0c, 0xb4,
	0xd9, 0xed, 0xce, 0xa2, 0xf4, 0xb1, 0x7c, 0x43, 0x41, 0x8d, 0x62, 0xa4, 0xc7, 0xf9, 0x66, 0xbe,
	0xf9, 0xf9, 0x21, 0xe1, 0x3d, 0x2b, 0xe7, 0x6d, 0xb0, 0x28, 0x4b, 0xbb, 0x57, 0x5c, 0xec, 0xc2,
	0xf1, 0x0b, 0x70, 0xda, 0xa7, 0x26, 0xd4, 0xed, 0x98, 0xbd, 0x0b, 0x90, 0xab, 0x6d, 0x28, 0xea,
	0xc6, 0x2c, 0xcd, 0xcc, 0x56, 0xa4, 0xe9, 0x80, 0x39, 0xc4, 0xa5, 0xad, 0x68, 0x73, 0x74, 0x34,
	0x12, 0x63, 0x31, 0xb9, 0x98, 0xde, 0xa8, 0xee, 0x1f, 0xd5, 0xb5, 0xd4, 0xec, 0x5b, 0xd1, 0x3f,
	0x32, 0x5e, 0x43, 0x12, 0x68, 0x47, 0xee, 0xc5, 0x36, 0x34, 0x8a, 0xc6, 0x62, 0x92, 0xe8, 0x5f,
	0x90, 0xdd, 0x41, 0xdc, 0x3a, 0xd8, 0x87, 0x58, 0xcf, 0xf3, 0x85, 0xd9, 0xcc, 0xb5, 0xec, 0x61,
	0x02, 0x67, 0x8f, 0xab, 0x7c, 0xf1, 0x24, 0x05, 0x0e, 0x61, 0xa0, 0x89, 0x29, 0xac, 0x0b, 0xe6,
	0x37, 0xeb, 0x2b, 0x19, 0x65, 0x97, 0x30, 0xec, 0x84, 0xb3, 0x9b, 0x3e, 0x03, 0x98, 0xa5, 0x31,
	0xe4, 0x5f, 0xeb, 0x92, 0xf0, 0x01, 0x06, 0x7f, 0x4e, 0xf0, 0xea, 0x7f, 0x81, 0xdc, 0xbb, 0x52,
	0xd3, 0x21, 0x4d, 0x4f, 0xad, 0xd8, 0x65, 0xbd, 0xfb, 0x68, 0x2d, 0xb6, 0xe7, 0x9f, 0xf8, 0xf6,
	0x23, 0x00, 0x00, 0xff, 0xff, 0x58, 0xf3, 0x93, 0xa2, 0x58, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SMSServiceClient is the client API for SMSService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SMSServiceClient interface {
	//获取验证码
	ObtainSMSCode(ctx context.Context, in *GrpcReq, opts ...grpc.CallOption) (*GrpcResp, error)
}

type sMSServiceClient struct {
	cc *grpc.ClientConn
}

func NewSMSServiceClient(cc *grpc.ClientConn) SMSServiceClient {
	return &sMSServiceClient{cc}
}

func (c *sMSServiceClient) ObtainSMSCode(ctx context.Context, in *GrpcReq, opts ...grpc.CallOption) (*GrpcResp, error) {
	out := new(GrpcResp)
	err := c.cc.Invoke(ctx, "/com.salty.protos.SMSService/ObtainSMSCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SMSServiceServer is the server API for SMSService service.
type SMSServiceServer interface {
	//获取验证码
	ObtainSMSCode(context.Context, *GrpcReq) (*GrpcResp, error)
}

// UnimplementedSMSServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSMSServiceServer struct {
}

func (*UnimplementedSMSServiceServer) ObtainSMSCode(ctx context.Context, req *GrpcReq) (*GrpcResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ObtainSMSCode not implemented")
}

func RegisterSMSServiceServer(s *grpc.Server, srv SMSServiceServer) {
	s.RegisterService(&_SMSService_serviceDesc, srv)
}

func _SMSService_ObtainSMSCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SMSServiceServer).ObtainSMSCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.salty.protos.SMSService/ObtainSMSCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SMSServiceServer).ObtainSMSCode(ctx, req.(*GrpcReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _SMSService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.salty.protos.SMSService",
	HandlerType: (*SMSServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ObtainSMSCode",
			Handler:    _SMSService_ObtainSMSCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sms.proto",
}
