// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: internal/plugin/action.proto

package plugin

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ActionService_Execute_FullMethodName  = "/plugin.ActionService/Execute"
	ActionService_Shutdown_FullMethodName = "/plugin.ActionService/Shutdown"
)

// ActionServiceClient is the client API for ActionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ActionServiceClient interface {
	Execute(ctx context.Context, in *ActionInput, opts ...grpc.CallOption) (*ActionOutput, error)
	Shutdown(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type actionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewActionServiceClient(cc grpc.ClientConnInterface) ActionServiceClient {
	return &actionServiceClient{cc}
}

func (c *actionServiceClient) Execute(ctx context.Context, in *ActionInput, opts ...grpc.CallOption) (*ActionOutput, error) {
	out := new(ActionOutput)
	err := c.cc.Invoke(ctx, ActionService_Execute_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actionServiceClient) Shutdown(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ActionService_Shutdown_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActionServiceServer is the server API for ActionService service.
// All implementations should embed UnimplementedActionServiceServer
// for forward compatibility
type ActionServiceServer interface {
	Execute(context.Context, *ActionInput) (*ActionOutput, error)
	Shutdown(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

// UnimplementedActionServiceServer should be embedded to have forward compatible implementations.
type UnimplementedActionServiceServer struct {
}

func (UnimplementedActionServiceServer) Execute(context.Context, *ActionInput) (*ActionOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedActionServiceServer) Shutdown(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}

// UnsafeActionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ActionServiceServer will
// result in compilation errors.
type UnsafeActionServiceServer interface {
	mustEmbedUnimplementedActionServiceServer()
}

func RegisterActionServiceServer(s grpc.ServiceRegistrar, srv ActionServiceServer) {
	s.RegisterService(&ActionService_ServiceDesc, srv)
}

func _ActionService_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActionInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActionServiceServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ActionService_Execute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActionServiceServer).Execute(ctx, req.(*ActionInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActionService_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActionServiceServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ActionService_Shutdown_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActionServiceServer).Shutdown(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ActionService_ServiceDesc is the grpc.ServiceDesc for ActionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ActionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.ActionService",
	HandlerType: (*ActionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _ActionService_Execute_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _ActionService_Shutdown_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/plugin/action.proto",
}