// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: api/api.proto

package api

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

var File_api_api_proto protoreflect.FileDescriptor

var file_api_api_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x61, 0x63, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x38, 0x0a, 0x04, 0x52,
	0x53, 0x43, 0x41, 0x12, 0x30, 0x0a, 0x04, 0x50, 0x69, 0x70, 0x65, 0x12, 0x11, 0x2e, 0x61, 0x63,
	0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x11,
	0x2e, 0x61, 0x63, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34, 0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_api_proto_goTypes = []interface{}{
	(*Message)(nil), // 0: acre.api.Message
}
var file_api_api_proto_depIdxs = []int32{
	0, // 0: acre.api.RSCA.Pipe:input_type -> acre.api.Message
	0, // 1: acre.api.RSCA.Pipe:output_type -> acre.api.Message
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_api_proto_init() }
func file_api_api_proto_init() {
	if File_api_api_proto != nil {
		return
	}
	file_api_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_api_proto_goTypes,
		DependencyIndexes: file_api_api_proto_depIdxs,
	}.Build()
	File_api_api_proto = out.File
	file_api_api_proto_rawDesc = nil
	file_api_api_proto_goTypes = nil
	file_api_api_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RSCAClient is the client API for RSCA service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RSCAClient interface {
	Pipe(ctx context.Context, opts ...grpc.CallOption) (RSCA_PipeClient, error)
}

type rSCAClient struct {
	cc grpc.ClientConnInterface
}

func NewRSCAClient(cc grpc.ClientConnInterface) RSCAClient {
	return &rSCAClient{cc}
}

func (c *rSCAClient) Pipe(ctx context.Context, opts ...grpc.CallOption) (RSCA_PipeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RSCA_serviceDesc.Streams[0], "/acre.api.RSCA/Pipe", opts...)
	if err != nil {
		return nil, err
	}
	x := &rSCAPipeClient{stream}
	return x, nil
}

type RSCA_PipeClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type rSCAPipeClient struct {
	grpc.ClientStream
}

func (x *rSCAPipeClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *rSCAPipeClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RSCAServer is the server API for RSCA service.
type RSCAServer interface {
	Pipe(RSCA_PipeServer) error
}

// UnimplementedRSCAServer can be embedded to have forward compatible implementations.
type UnimplementedRSCAServer struct {
}

func (*UnimplementedRSCAServer) Pipe(RSCA_PipeServer) error {
	return status.Errorf(codes.Unimplemented, "method Pipe not implemented")
}

func RegisterRSCAServer(s *grpc.Server, srv RSCAServer) {
	s.RegisterService(&_RSCA_serviceDesc, srv)
}

func _RSCA_Pipe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RSCAServer).Pipe(&rSCAPipeServer{stream})
}

type RSCA_PipeServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type rSCAPipeServer struct {
	grpc.ServerStream
}

func (x *rSCAPipeServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *rSCAPipeServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _RSCA_serviceDesc = grpc.ServiceDesc{
	ServiceName: "acre.api.RSCA",
	HandlerType: (*RSCAServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Pipe",
			Handler:       _RSCA_Pipe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/api.proto",
}
