// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: github.com/na4ma4/rsca/api/admin.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Admin_ListHosts_FullMethodName   = "/rsca.api.Admin/ListHosts"
	Admin_RemoveHost_FullMethodName  = "/rsca.api.Admin/RemoveHost"
	Admin_TriggerAll_FullMethodName  = "/rsca.api.Admin/TriggerAll"
	Admin_TriggerInfo_FullMethodName = "/rsca.api.Admin/TriggerInfo"
)

// AdminClient is the client API for Admin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminClient interface {
	ListHosts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Admin_ListHostsClient, error)
	RemoveHost(ctx context.Context, in *RemoveHostRequest, opts ...grpc.CallOption) (*RemoveHostResponse, error)
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
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Admin_ServiceDesc.Streams[0], Admin_ListHosts_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &adminListHostsClient{ClientStream: stream}
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

func (c *adminClient) RemoveHost(ctx context.Context, in *RemoveHostRequest, opts ...grpc.CallOption) (*RemoveHostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveHostResponse)
	err := c.cc.Invoke(ctx, Admin_RemoveHost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) TriggerAll(ctx context.Context, in *Members, opts ...grpc.CallOption) (*TriggerAllResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TriggerAllResponse)
	err := c.cc.Invoke(ctx, Admin_TriggerAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) TriggerInfo(ctx context.Context, in *Members, opts ...grpc.CallOption) (*TriggerInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TriggerInfoResponse)
	err := c.cc.Invoke(ctx, Admin_TriggerInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServer is the server API for Admin service.
// All implementations should embed UnimplementedAdminServer
// for forward compatibility
type AdminServer interface {
	ListHosts(*Empty, Admin_ListHostsServer) error
	RemoveHost(context.Context, *RemoveHostRequest) (*RemoveHostResponse, error)
	TriggerAll(context.Context, *Members) (*TriggerAllResponse, error)
	TriggerInfo(context.Context, *Members) (*TriggerInfoResponse, error)
}

// UnimplementedAdminServer should be embedded to have forward compatible implementations.
type UnimplementedAdminServer struct {
}

func (UnimplementedAdminServer) ListHosts(*Empty, Admin_ListHostsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListHosts not implemented")
}
func (UnimplementedAdminServer) RemoveHost(context.Context, *RemoveHostRequest) (*RemoveHostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveHost not implemented")
}
func (UnimplementedAdminServer) TriggerAll(context.Context, *Members) (*TriggerAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerAll not implemented")
}
func (UnimplementedAdminServer) TriggerInfo(context.Context, *Members) (*TriggerInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerInfo not implemented")
}

// UnsafeAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServer will
// result in compilation errors.
type UnsafeAdminServer interface {
	mustEmbedUnimplementedAdminServer()
}

func RegisterAdminServer(s grpc.ServiceRegistrar, srv AdminServer) {
	s.RegisterService(&Admin_ServiceDesc, srv)
}

func _Admin_ListHosts_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AdminServer).ListHosts(m, &adminListHostsServer{ServerStream: stream})
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

func _Admin_RemoveHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).RemoveHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Admin_RemoveHost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).RemoveHost(ctx, req.(*RemoveHostRequest))
	}
	return interceptor(ctx, in, info, handler)
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
		FullMethod: Admin_TriggerAll_FullMethodName,
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
		FullMethod: Admin_TriggerInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).TriggerInfo(ctx, req.(*Members))
	}
	return interceptor(ctx, in, info, handler)
}

// Admin_ServiceDesc is the grpc.ServiceDesc for Admin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Admin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rsca.api.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RemoveHost",
			Handler:    _Admin_RemoveHost_Handler,
		},
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
