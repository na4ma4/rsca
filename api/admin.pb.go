// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.7
// source: github.com/na4ma4/rsca/api/admin.proto

package api

import (
	context "context"
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

var File_github_com_na4ma4_rsca_api_admin_proto protoreflect.FileDescriptor

var file_github_com_na4ma4_rsca_api_admin_proto_rawDesc = []byte{
	0x0a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34,
	0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61,
	0x70, 0x69, 0x1a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e,
	0x61, 0x34, 0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xb9, 0x01, 0x0a, 0x05,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x30, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x6f, 0x73,
	0x74, 0x73, 0x12, 0x0f, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x30, 0x01, 0x12, 0x3d, 0x0a, 0x0a, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x41, 0x6c, 0x6c, 0x12, 0x11, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x1c, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x41, 0x6c, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x11, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x1d, 0x2e, 0x72, 0x73, 0x63, 0x61, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x34, 0x6d, 0x61, 0x34, 0x2f, 0x72, 0x73, 0x63,
	0x61, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_github_com_na4ma4_rsca_api_admin_proto_goTypes = []interface{}{
	(*Empty)(nil),               // 0: rsca.api.Empty
	(*Members)(nil),             // 1: rsca.api.Members
	(*Member)(nil),              // 2: rsca.api.Member
	(*TriggerAllResponse)(nil),  // 3: rsca.api.TriggerAllResponse
	(*TriggerInfoResponse)(nil), // 4: rsca.api.TriggerInfoResponse
}
var file_github_com_na4ma4_rsca_api_admin_proto_depIdxs = []int32{
	0, // 0: rsca.api.Admin.ListHosts:input_type -> rsca.api.Empty
	1, // 1: rsca.api.Admin.TriggerAll:input_type -> rsca.api.Members
	1, // 2: rsca.api.Admin.TriggerInfo:input_type -> rsca.api.Members
	2, // 3: rsca.api.Admin.ListHosts:output_type -> rsca.api.Member
	3, // 4: rsca.api.Admin.TriggerAll:output_type -> rsca.api.TriggerAllResponse
	4, // 5: rsca.api.Admin.TriggerInfo:output_type -> rsca.api.TriggerInfoResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_na4ma4_rsca_api_admin_proto_init() }
func file_github_com_na4ma4_rsca_api_admin_proto_init() {
	if File_github_com_na4ma4_rsca_api_admin_proto != nil {
		return
	}
	file_github_com_na4ma4_rsca_api_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_na4ma4_rsca_api_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_na4ma4_rsca_api_admin_proto_goTypes,
		DependencyIndexes: file_github_com_na4ma4_rsca_api_admin_proto_depIdxs,
	}.Build()
	File_github_com_na4ma4_rsca_api_admin_proto = out.File
	file_github_com_na4ma4_rsca_api_admin_proto_rawDesc = nil
	file_github_com_na4ma4_rsca_api_admin_proto_goTypes = nil
	file_github_com_na4ma4_rsca_api_admin_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AdminClient is the client API for Admin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdminClient interface {
	ListHosts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Admin_ListHostsClient, error)
	TriggerAll(ctx context.Context, in *Members, opts ...grpc.CallOption) (*TriggerAllResponse, error)
	TriggerInfo(ctx context.Context, in *Members, opts ...grpc.CallOption) (*TriggerInfoResponse, error)
}

type adminClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminClient(cc grpc.ClientConnInterface) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) ListHosts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Admin_ListHostsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Admin_serviceDesc.Streams[0], "/rsca.api.Admin/ListHosts", opts...)
	if err != nil {
		return nil, err
	}
	x := &adminListHostsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Admin_ListHostsClient interface {
	Recv() (*Member, error)
	grpc.ClientStream
}

type adminListHostsClient struct {
	grpc.ClientStream
}

func (x *adminListHostsClient) Recv() (*Member, error) {
	m := new(Member)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *adminClient) TriggerAll(ctx context.Context, in *Members, opts ...grpc.CallOption) (*TriggerAllResponse, error) {
	out := new(TriggerAllResponse)
	err := c.cc.Invoke(ctx, "/rsca.api.Admin/TriggerAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) TriggerInfo(ctx context.Context, in *Members, opts ...grpc.CallOption) (*TriggerInfoResponse, error) {
	out := new(TriggerInfoResponse)
	err := c.cc.Invoke(ctx, "/rsca.api.Admin/TriggerInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServer is the server API for Admin service.
type AdminServer interface {
	ListHosts(*Empty, Admin_ListHostsServer) error
	TriggerAll(context.Context, *Members) (*TriggerAllResponse, error)
	TriggerInfo(context.Context, *Members) (*TriggerInfoResponse, error)
}

// UnimplementedAdminServer can be embedded to have forward compatible implementations.
type UnimplementedAdminServer struct {
}

func (*UnimplementedAdminServer) ListHosts(*Empty, Admin_ListHostsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListHosts not implemented")
}
func (*UnimplementedAdminServer) TriggerAll(context.Context, *Members) (*TriggerAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerAll not implemented")
}
func (*UnimplementedAdminServer) TriggerInfo(context.Context, *Members) (*TriggerInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerInfo not implemented")
}

func RegisterAdminServer(s *grpc.Server, srv AdminServer) {
	s.RegisterService(&_Admin_serviceDesc, srv)
}

func _Admin_ListHosts_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AdminServer).ListHosts(m, &adminListHostsServer{stream})
}

type Admin_ListHostsServer interface {
	Send(*Member) error
	grpc.ServerStream
}

type adminListHostsServer struct {
	grpc.ServerStream
}

func (x *adminListHostsServer) Send(m *Member) error {
	return x.ServerStream.SendMsg(m)
}

func _Admin_TriggerAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Members)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).TriggerAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsca.api.Admin/TriggerAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).TriggerAll(ctx, req.(*Members))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_TriggerInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Members)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).TriggerInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rsca.api.Admin/TriggerInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).TriggerInfo(ctx, req.(*Members))
	}
	return interceptor(ctx, in, info, handler)
}

var _Admin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rsca.api.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TriggerAll",
			Handler:    _Admin_TriggerAll_Handler,
		},
		{
			MethodName: "TriggerInfo",
			Handler:    _Admin_TriggerInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListHosts",
			Handler:       _Admin_ListHosts_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/na4ma4/rsca/api/admin.proto",
}
