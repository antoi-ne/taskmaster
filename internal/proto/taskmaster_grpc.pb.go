// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: internal/proto/taskmaster.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TaskmasterClient is the client API for Taskmaster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskmasterClient interface {
	Reload(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	List(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ServiceStatusList, error)
	Status(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error)
	Start(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error)
	Stop(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error)
	Restart(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error)
}

type taskmasterClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskmasterClient(cc grpc.ClientConnInterface) TaskmasterClient {
	return &taskmasterClient{cc}
}

func (c *taskmasterClient) Reload(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Taskmaster/Reload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) List(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ServiceStatusList, error) {
	out := new(ServiceStatusList)
	err := c.cc.Invoke(ctx, "/Taskmaster/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) Status(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error) {
	out := new(ServiceStatus)
	err := c.cc.Invoke(ctx, "/Taskmaster/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) Start(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error) {
	out := new(ServiceStatus)
	err := c.cc.Invoke(ctx, "/Taskmaster/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) Stop(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error) {
	out := new(ServiceStatus)
	err := c.cc.Invoke(ctx, "/Taskmaster/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) Restart(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceStatus, error) {
	out := new(ServiceStatus)
	err := c.cc.Invoke(ctx, "/Taskmaster/Restart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskmasterServer is the server API for Taskmaster service.
// All implementations must embed UnimplementedTaskmasterServer
// for forward compatibility
type TaskmasterServer interface {
	Reload(context.Context, *Empty) (*Empty, error)
	List(context.Context, *Empty) (*ServiceStatusList, error)
	Status(context.Context, *Service) (*ServiceStatus, error)
	Start(context.Context, *Service) (*ServiceStatus, error)
	Stop(context.Context, *Service) (*ServiceStatus, error)
	Restart(context.Context, *Service) (*ServiceStatus, error)
	mustEmbedUnimplementedTaskmasterServer()
}

// UnimplementedTaskmasterServer must be embedded to have forward compatible implementations.
type UnimplementedTaskmasterServer struct {
}

func (UnimplementedTaskmasterServer) Reload(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reload not implemented")
}
func (UnimplementedTaskmasterServer) List(context.Context, *Empty) (*ServiceStatusList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedTaskmasterServer) Status(context.Context, *Service) (*ServiceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedTaskmasterServer) Start(context.Context, *Service) (*ServiceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedTaskmasterServer) Stop(context.Context, *Service) (*ServiceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedTaskmasterServer) Restart(context.Context, *Service) (*ServiceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Restart not implemented")
}
func (UnimplementedTaskmasterServer) mustEmbedUnimplementedTaskmasterServer() {}

// UnsafeTaskmasterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskmasterServer will
// result in compilation errors.
type UnsafeTaskmasterServer interface {
	mustEmbedUnimplementedTaskmasterServer()
}

func RegisterTaskmasterServer(s grpc.ServiceRegistrar, srv TaskmasterServer) {
	s.RegisterService(&Taskmaster_ServiceDesc, srv)
}

func _Taskmaster_Reload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).Reload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/Reload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).Reload(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).List(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).Status(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).Start(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).Stop(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_Restart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).Restart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/Restart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).Restart(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

// Taskmaster_ServiceDesc is the grpc.ServiceDesc for Taskmaster service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Taskmaster_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Taskmaster",
	HandlerType: (*TaskmasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reload",
			Handler:    _Taskmaster_Reload_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Taskmaster_List_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Taskmaster_Status_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _Taskmaster_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Taskmaster_Stop_Handler,
		},
		{
			MethodName: "Restart",
			Handler:    _Taskmaster_Restart_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/taskmaster.proto",
}
