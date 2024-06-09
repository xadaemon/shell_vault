// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: adminproto/admin.proto

package adminproto

import (
	context "context"
	common "github.com/stateprism/shell_vault/rpc/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminServiceClient interface {
	Authenticate(ctx context.Context, in *common.Empty, opts ...grpc.CallOption) (*common.AuthReply, error)
	ChangeRootCert(ctx context.Context, in *ChangeRootCertRequest, opts ...grpc.CallOption) (*ChangeRootCertResponse, error)
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*UserActionResponse, error)
	RestartServer(ctx context.Context, in *StopServerRequest, opts ...grpc.CallOption) (*common.Empty, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) Authenticate(ctx context.Context, in *common.Empty, opts ...grpc.CallOption) (*common.AuthReply, error) {
	out := new(common.AuthReply)
	err := c.cc.Invoke(ctx, "/AdminService/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) ChangeRootCert(ctx context.Context, in *ChangeRootCertRequest, opts ...grpc.CallOption) (*ChangeRootCertResponse, error) {
	out := new(ChangeRootCertResponse)
	err := c.cc.Invoke(ctx, "/AdminService/ChangeRootCert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*UserActionResponse, error) {
	out := new(UserActionResponse)
	err := c.cc.Invoke(ctx, "/AdminService/AddUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) RestartServer(ctx context.Context, in *StopServerRequest, opts ...grpc.CallOption) (*common.Empty, error) {
	out := new(common.Empty)
	err := c.cc.Invoke(ctx, "/AdminService/RestartServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
// All implementations must embed UnimplementedAdminServiceServer
// for forward compatibility
type AdminServiceServer interface {
	Authenticate(context.Context, *common.Empty) (*common.AuthReply, error)
	ChangeRootCert(context.Context, *ChangeRootCertRequest) (*ChangeRootCertResponse, error)
	AddUser(context.Context, *AddUserRequest) (*UserActionResponse, error)
	RestartServer(context.Context, *StopServerRequest) (*common.Empty, error)
	mustEmbedUnimplementedAdminServiceServer()
}

// UnimplementedAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (UnimplementedAdminServiceServer) Authenticate(context.Context, *common.Empty) (*common.AuthReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedAdminServiceServer) ChangeRootCert(context.Context, *ChangeRootCertRequest) (*ChangeRootCertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeRootCert not implemented")
}
func (UnimplementedAdminServiceServer) AddUser(context.Context, *AddUserRequest) (*UserActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedAdminServiceServer) RestartServer(context.Context, *StopServerRequest) (*common.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartServer not implemented")
}
func (UnimplementedAdminServiceServer) mustEmbedUnimplementedAdminServiceServer() {}

// UnsafeAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServiceServer will
// result in compilation errors.
type UnsafeAdminServiceServer interface {
	mustEmbedUnimplementedAdminServiceServer()
}

func RegisterAdminServiceServer(s grpc.ServiceRegistrar, srv AdminServiceServer) {
	s.RegisterService(&AdminService_ServiceDesc, srv)
}

func _AdminService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).Authenticate(ctx, req.(*common.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_ChangeRootCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeRootCertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ChangeRootCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/ChangeRootCert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ChangeRootCert(ctx, req.(*ChangeRootCertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/AddUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_RestartServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).RestartServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/RestartServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).RestartServer(ctx, req.(*StopServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminService_ServiceDesc is the grpc.ServiceDesc for AdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _AdminService_Authenticate_Handler,
		},
		{
			MethodName: "ChangeRootCert",
			Handler:    _AdminService_ChangeRootCert_Handler,
		},
		{
			MethodName: "AddUser",
			Handler:    _AdminService_AddUser_Handler,
		},
		{
			MethodName: "RestartServer",
			Handler:    _AdminService_RestartServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "adminproto/admin.proto",
}
