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
	List(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ProgramDescList, error)
	Reload(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	ProgramStatus(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error)
	ProgramStart(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error)
	ProgramStop(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error)
	ProgramRestart(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error)
}

type taskmasterClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskmasterClient(cc grpc.ClientConnInterface) TaskmasterClient {
	return &taskmasterClient{cc}
}

func (c *taskmasterClient) List(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ProgramDescList, error) {
	out := new(ProgramDescList)
	err := c.cc.Invoke(ctx, "/Taskmaster/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) Reload(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Taskmaster/Reload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Taskmaster/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) ProgramStatus(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error) {
	out := new(ProgramDesc)
	err := c.cc.Invoke(ctx, "/Taskmaster/ProgramStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) ProgramStart(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error) {
	out := new(ProgramDesc)
	err := c.cc.Invoke(ctx, "/Taskmaster/ProgramStart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) ProgramStop(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error) {
	out := new(ProgramDesc)
	err := c.cc.Invoke(ctx, "/Taskmaster/ProgramStop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) ProgramRestart(ctx context.Context, in *Program, opts ...grpc.CallOption) (*ProgramDesc, error) {
	out := new(ProgramDesc)
	err := c.cc.Invoke(ctx, "/Taskmaster/ProgramRestart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskmasterServer is the server API for Taskmaster service.
// All implementations must embed UnimplementedTaskmasterServer
// for forward compatibility
type TaskmasterServer interface {
	List(context.Context, *Empty) (*ProgramDescList, error)
	Reload(context.Context, *Empty) (*Empty, error)
	Stop(context.Context, *Empty) (*Empty, error)
	ProgramStatus(context.Context, *Program) (*ProgramDesc, error)
	ProgramStart(context.Context, *Program) (*ProgramDesc, error)
	ProgramStop(context.Context, *Program) (*ProgramDesc, error)
	ProgramRestart(context.Context, *Program) (*ProgramDesc, error)
	mustEmbedUnimplementedTaskmasterServer()
}

// UnimplementedTaskmasterServer must be embedded to have forward compatible implementations.
type UnimplementedTaskmasterServer struct {
}

func (UnimplementedTaskmasterServer) List(context.Context, *Empty) (*ProgramDescList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedTaskmasterServer) Reload(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reload not implemented")
}
func (UnimplementedTaskmasterServer) Stop(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedTaskmasterServer) ProgramStatus(context.Context, *Program) (*ProgramDesc, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProgramStatus not implemented")
}
func (UnimplementedTaskmasterServer) ProgramStart(context.Context, *Program) (*ProgramDesc, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProgramStart not implemented")
}
func (UnimplementedTaskmasterServer) ProgramStop(context.Context, *Program) (*ProgramDesc, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProgramStop not implemented")
}
func (UnimplementedTaskmasterServer) ProgramRestart(context.Context, *Program) (*ProgramDesc, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProgramRestart not implemented")
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

func _Taskmaster_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
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
		return srv.(TaskmasterServer).Stop(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_ProgramStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Program)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).ProgramStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/ProgramStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).ProgramStatus(ctx, req.(*Program))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_ProgramStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Program)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).ProgramStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/ProgramStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).ProgramStart(ctx, req.(*Program))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_ProgramStop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Program)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).ProgramStop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/ProgramStop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).ProgramStop(ctx, req.(*Program))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_ProgramRestart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Program)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).ProgramRestart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Taskmaster/ProgramRestart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).ProgramRestart(ctx, req.(*Program))
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
			MethodName: "List",
			Handler:    _Taskmaster_List_Handler,
		},
		{
			MethodName: "Reload",
			Handler:    _Taskmaster_Reload_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Taskmaster_Stop_Handler,
		},
		{
			MethodName: "ProgramStatus",
			Handler:    _Taskmaster_ProgramStatus_Handler,
		},
		{
			MethodName: "ProgramStart",
			Handler:    _Taskmaster_ProgramStart_Handler,
		},
		{
			MethodName: "ProgramStop",
			Handler:    _Taskmaster_ProgramStop_Handler,
		},
		{
			MethodName: "ProgramRestart",
			Handler:    _Taskmaster_ProgramRestart_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/taskmaster.proto",
}
