// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

/*
Package rpcconfig is a generated protocol buffer package.

It is generated from these files:
	config.proto

It has these top-level messages:
	Content
	Nop
*/
package rpcconfig

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type Content struct {
	Content []byte `protobuf:"bytes,1,opt,name=Content,json=content,proto3" json:"Content,omitempty"`
}

func (m *Content) Reset()                    { *m = Content{} }
func (m *Content) String() string            { return proto.CompactTextString(m) }
func (*Content) ProtoMessage()               {}
func (*Content) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Content) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type Nop struct {
}

func (m *Nop) Reset()                    { *m = Nop{} }
func (m *Nop) String() string            { return proto.CompactTextString(m) }
func (*Nop) ProtoMessage()               {}
func (*Nop) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Content)(nil), "rpcconfig.Content")
	proto.RegisterType((*Nop)(nil), "rpcconfig.Nop")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ConfigurationServer service

type ConfigurationServerClient interface {
	Set(ctx context.Context, in *Content, opts ...grpc.CallOption) (*Content, error)
	Get(ctx context.Context, in *Content, opts ...grpc.CallOption) (*Content, error)
	Save(ctx context.Context, in *Nop, opts ...grpc.CallOption) (*Nop, error)
}

type configurationServerClient struct {
	cc *grpc.ClientConn
}

func NewConfigurationServerClient(cc *grpc.ClientConn) ConfigurationServerClient {
	return &configurationServerClient{cc}
}

func (c *configurationServerClient) Set(ctx context.Context, in *Content, opts ...grpc.CallOption) (*Content, error) {
	out := new(Content)
	err := grpc.Invoke(ctx, "/rpcconfig.ConfigurationServer/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configurationServerClient) Get(ctx context.Context, in *Content, opts ...grpc.CallOption) (*Content, error) {
	out := new(Content)
	err := grpc.Invoke(ctx, "/rpcconfig.ConfigurationServer/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configurationServerClient) Save(ctx context.Context, in *Nop, opts ...grpc.CallOption) (*Nop, error) {
	out := new(Nop)
	err := grpc.Invoke(ctx, "/rpcconfig.ConfigurationServer/Save", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ConfigurationServer service

type ConfigurationServerServer interface {
	Set(context.Context, *Content) (*Content, error)
	Get(context.Context, *Content) (*Content, error)
	Save(context.Context, *Nop) (*Nop, error)
}

func RegisterConfigurationServerServer(s *grpc.Server, srv ConfigurationServerServer) {
	s.RegisterService(&_ConfigurationServer_serviceDesc, srv)
}

func _ConfigurationServer_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Content)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationServerServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcconfig.ConfigurationServer/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationServerServer).Set(ctx, req.(*Content))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigurationServer_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Content)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationServerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcconfig.ConfigurationServer/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationServerServer).Get(ctx, req.(*Content))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigurationServer_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nop)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationServerServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcconfig.ConfigurationServer/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationServerServer).Save(ctx, req.(*Nop))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConfigurationServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpcconfig.ConfigurationServer",
	HandlerType: (*ConfigurationServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _ConfigurationServer_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ConfigurationServer_Get_Handler,
		},
		{
			MethodName: "Save",
			Handler:    _ConfigurationServer_Save_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "config.proto",
}

func init() { proto.RegisterFile("config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 178 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2c, 0x2a, 0x48, 0x86, 0x08, 0x28,
	0x29, 0x73, 0xb1, 0x3b, 0xe7, 0xe7, 0x95, 0xa4, 0xe6, 0x95, 0x08, 0x49, 0xc0, 0x99, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0xec, 0xc9, 0x10, 0xae, 0x12, 0x2b, 0x17, 0xb3, 0x5f, 0x7e, 0x81,
	0xd1, 0x42, 0x46, 0x2e, 0x61, 0x67, 0xb0, 0xb6, 0xd2, 0xa2, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0xe0,
	0xd4, 0xa2, 0xb2, 0xd4, 0x22, 0x21, 0x7d, 0x2e, 0xe6, 0xe0, 0xd4, 0x12, 0x21, 0x21, 0x3d, 0xb8,
	0xb1, 0x7a, 0x50, 0x83, 0xa4, 0xb0, 0x88, 0x29, 0x31, 0x80, 0x34, 0xb8, 0x93, 0xa4, 0x41, 0x83,
	0x8b, 0x25, 0x38, 0xb1, 0x2c, 0x55, 0x88, 0x0f, 0x49, 0xd6, 0x2f, 0xbf, 0x40, 0x0a, 0x8d, 0xaf,
	0xc4, 0xe0, 0x64, 0xc8, 0xa5, 0x90, 0x8c, 0xec, 0x44, 0xbd, 0xf4, 0xa2, 0x82, 0x64, 0xbd, 0x62,
	0xb0, 0x43, 0xf5, 0x20, 0x32, 0x4e, 0xdc, 0x10, 0x4f, 0x04, 0x80, 0xc2, 0x22, 0x80, 0x31, 0x89,
	0x0d, 0x1c, 0x28, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf1, 0x3e, 0xda, 0x7e, 0x24, 0x01,
	0x00, 0x00,
}
