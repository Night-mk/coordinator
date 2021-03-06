// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.3
// source: stateTransfer.proto

package protocol

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

// CommitServiceClient is the client API for CommitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommitServiceClient interface {
	Propose(ctx context.Context, in *ProposeRequest, opts ...grpc.CallOption) (*ProposeResponse, error)
	Commit(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*Response, error)
	StateTransfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Value, error)
}

type commitServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommitServiceClient(cc grpc.ClientConnInterface) CommitServiceClient {
	return &commitServiceClient{cc}
}

func (c *commitServiceClient) Propose(ctx context.Context, in *ProposeRequest, opts ...grpc.CallOption) (*ProposeResponse, error) {
	out := new(ProposeResponse)
	err := c.cc.Invoke(ctx, "/protocol.CommitService/Propose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commitServiceClient) Commit(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/protocol.CommitService/Commit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commitServiceClient) StateTransfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/protocol.CommitService/StateTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commitServiceClient) Get(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Value, error) {
	out := new(Value)
	err := c.cc.Invoke(ctx, "/protocol.CommitService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommitServiceServer is the server API for CommitService service.
// All implementations must embed UnimplementedCommitServiceServer
// for forward compatibility
type CommitServiceServer interface {
	Propose(context.Context, *ProposeRequest) (*ProposeResponse, error)
	Commit(context.Context, *CommitRequest) (*Response, error)
	StateTransfer(context.Context, *TransferRequest) (*Response, error)
	Get(context.Context, *Msg) (*Value, error)
	mustEmbedUnimplementedCommitServiceServer()
}

// UnimplementedCommitServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommitServiceServer struct {
}

func (UnimplementedCommitServiceServer) Propose(context.Context, *ProposeRequest) (*ProposeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Propose not implemented")
}
func (UnimplementedCommitServiceServer) Commit(context.Context, *CommitRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedCommitServiceServer) StateTransfer(context.Context, *TransferRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StateTransfer not implemented")
}
func (UnimplementedCommitServiceServer) Get(context.Context, *Msg) (*Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCommitServiceServer) mustEmbedUnimplementedCommitServiceServer() {}

// UnsafeCommitServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommitServiceServer will
// result in compilation errors.
type UnsafeCommitServiceServer interface {
	mustEmbedUnimplementedCommitServiceServer()
}

func RegisterCommitServiceServer(s grpc.ServiceRegistrar, srv CommitServiceServer) {
	s.RegisterService(&CommitService_ServiceDesc, srv)
}

func _CommitService_Propose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProposeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommitServiceServer).Propose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.CommitService/Propose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommitServiceServer).Propose(ctx, req.(*ProposeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommitService_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommitServiceServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.CommitService/Commit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommitServiceServer).Commit(ctx, req.(*CommitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommitService_StateTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommitServiceServer).StateTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.CommitService/StateTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommitServiceServer).StateTransfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommitService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommitServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.CommitService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommitServiceServer).Get(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

// CommitService_ServiceDesc is the grpc.ServiceDesc for CommitService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommitService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.CommitService",
	HandlerType: (*CommitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Propose",
			Handler:    _CommitService_Propose_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _CommitService_Commit_Handler,
		},
		{
			MethodName: "StateTransfer",
			Handler:    _CommitService_StateTransfer_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _CommitService_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stateTransfer.proto",
}
