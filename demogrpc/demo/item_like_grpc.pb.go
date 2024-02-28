// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: demogrpc/item_like.proto

package demo

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

// ItemLikeServiceClient is the client API for ItemLikeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ItemLikeServiceClient interface {
	GetItemLikes(ctx context.Context, in *GetItemLikesReq, opts ...grpc.CallOption) (*ItemLikeResp, error)
}

type itemLikeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewItemLikeServiceClient(cc grpc.ClientConnInterface) ItemLikeServiceClient {
	return &itemLikeServiceClient{cc}
}

func (c *itemLikeServiceClient) GetItemLikes(ctx context.Context, in *GetItemLikesReq, opts ...grpc.CallOption) (*ItemLikeResp, error) {
	out := new(ItemLikeResp)
	err := c.cc.Invoke(ctx, "/demo.ItemLikeService/GetItemLikes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemLikeServiceServer is the server API for ItemLikeService service.
// All implementations should embed UnimplementedItemLikeServiceServer
// for forward compatibility
type ItemLikeServiceServer interface {
	GetItemLikes(context.Context, *GetItemLikesReq) (*ItemLikeResp, error)
}

// UnimplementedItemLikeServiceServer should be embedded to have forward compatible implementations.
type UnimplementedItemLikeServiceServer struct {
}

func (UnimplementedItemLikeServiceServer) GetItemLikes(context.Context, *GetItemLikesReq) (*ItemLikeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemLikes not implemented")
}

// UnsafeItemLikeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ItemLikeServiceServer will
// result in compilation errors.
type UnsafeItemLikeServiceServer interface {
	mustEmbedUnimplementedItemLikeServiceServer()
}

func RegisterItemLikeServiceServer(s grpc.ServiceRegistrar, srv ItemLikeServiceServer) {
	s.RegisterService(&ItemLikeService_ServiceDesc, srv)
}

func _ItemLikeService_GetItemLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemLikesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemLikeServiceServer).GetItemLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.ItemLikeService/GetItemLikes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemLikeServiceServer).GetItemLikes(ctx, req.(*GetItemLikesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ItemLikeService_ServiceDesc is the grpc.ServiceDesc for ItemLikeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ItemLikeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "demo.ItemLikeService",
	HandlerType: (*ItemLikeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetItemLikes",
			Handler:    _ItemLikeService_GetItemLikes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "demogrpc/item_like.proto",
}