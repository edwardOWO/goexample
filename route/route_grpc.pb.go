// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: route.proto

package route

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

// RouteGuideClient is the client API for RouteGuide service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouteGuideClient interface {
	// Unary
	GetFeature(ctx context.Context, in *Point, opts ...grpc.CallOption) (*Feature, error)
}

type routeGuideClient struct {
	cc grpc.ClientConnInterface
}

func NewRouteGuideClient(cc grpc.ClientConnInterface) RouteGuideClient {
	return &routeGuideClient{cc}
}

func (c *routeGuideClient) GetFeature(ctx context.Context, in *Point, opts ...grpc.CallOption) (*Feature, error) {
	out := new(Feature)
	err := c.cc.Invoke(ctx, "/route.RouteGuide/GetFeature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RouteGuideServer is the server API for RouteGuide service.
// All implementations must embed UnimplementedRouteGuideServer
// for forward compatibility
type RouteGuideServer interface {
	// Unary
	GetFeature(context.Context, *Point) (*Feature, error)
	mustEmbedUnimplementedRouteGuideServer()
}

// UnimplementedRouteGuideServer must be embedded to have forward compatible implementations.
type UnimplementedRouteGuideServer struct {
}

func (UnimplementedRouteGuideServer) GetFeature(context.Context, *Point) (*Feature, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeature not implemented")
}
func (UnimplementedRouteGuideServer) mustEmbedUnimplementedRouteGuideServer() {}

// UnsafeRouteGuideServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RouteGuideServer will
// result in compilation errors.
type UnsafeRouteGuideServer interface {
	mustEmbedUnimplementedRouteGuideServer()
}

func RegisterRouteGuideServer(s grpc.ServiceRegistrar, srv RouteGuideServer) {
	s.RegisterService(&RouteGuide_ServiceDesc, srv)
}

func _RouteGuide_GetFeature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Point)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouteGuideServer).GetFeature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route.RouteGuide/GetFeature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouteGuideServer).GetFeature(ctx, req.(*Point))
	}
	return interceptor(ctx, in, info, handler)
}

// RouteGuide_ServiceDesc is the grpc.ServiceDesc for RouteGuide service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RouteGuide_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "route.RouteGuide",
	HandlerType: (*RouteGuideServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFeature",
			Handler:    _RouteGuide_GetFeature_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "route.proto",
}

// GetTimeClient is the client API for GetTime service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GetTimeClient interface {
	// Sends a greeting
	GetTime(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*TimeReply, error)
}

type getTimeClient struct {
	cc grpc.ClientConnInterface
}

func NewGetTimeClient(cc grpc.ClientConnInterface) GetTimeClient {
	return &getTimeClient{cc}
}

func (c *getTimeClient) GetTime(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*TimeReply, error) {
	out := new(TimeReply)
	err := c.cc.Invoke(ctx, "/route.GetTime/GetTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetTimeServer is the server API for GetTime service.
// All implementations must embed UnimplementedGetTimeServer
// for forward compatibility
type GetTimeServer interface {
	// Sends a greeting
	GetTime(context.Context, *TimeRequest) (*TimeReply, error)
	mustEmbedUnimplementedGetTimeServer()
}

// UnimplementedGetTimeServer must be embedded to have forward compatible implementations.
type UnimplementedGetTimeServer struct {
}

func (UnimplementedGetTimeServer) GetTime(context.Context, *TimeRequest) (*TimeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTime not implemented")
}
func (UnimplementedGetTimeServer) mustEmbedUnimplementedGetTimeServer() {}

// UnsafeGetTimeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GetTimeServer will
// result in compilation errors.
type UnsafeGetTimeServer interface {
	mustEmbedUnimplementedGetTimeServer()
}

func RegisterGetTimeServer(s grpc.ServiceRegistrar, srv GetTimeServer) {
	s.RegisterService(&GetTime_ServiceDesc, srv)
}

func _GetTime_GetTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetTimeServer).GetTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route.GetTime/GetTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetTimeServer).GetTime(ctx, req.(*TimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GetTime_ServiceDesc is the grpc.ServiceDesc for GetTime service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GetTime_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "route.GetTime",
	HandlerType: (*GetTimeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTime",
			Handler:    _GetTime_GetTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "route.proto",
}

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamServiceClient interface {
	FetchResponse(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (StreamService_FetchResponseClient, error)
}

type streamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamServiceClient(cc grpc.ClientConnInterface) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) FetchResponse(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (StreamService_FetchResponseClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[0], "/route.StreamService/FetchResponse", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceFetchResponseClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_FetchResponseClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type streamServiceFetchResponseClient struct {
	grpc.ClientStream
}

func (x *streamServiceFetchResponseClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamServiceServer is the server API for StreamService service.
// All implementations must embed UnimplementedStreamServiceServer
// for forward compatibility
type StreamServiceServer interface {
	FetchResponse(*StreamRequest, StreamService_FetchResponseServer) error
	mustEmbedUnimplementedStreamServiceServer()
}

// UnimplementedStreamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServiceServer struct {
}

func (UnimplementedStreamServiceServer) FetchResponse(*StreamRequest, StreamService_FetchResponseServer) error {
	return status.Errorf(codes.Unimplemented, "method FetchResponse not implemented")
}
func (UnimplementedStreamServiceServer) mustEmbedUnimplementedStreamServiceServer() {}

// UnsafeStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServiceServer will
// result in compilation errors.
type UnsafeStreamServiceServer interface {
	mustEmbedUnimplementedStreamServiceServer()
}

func RegisterStreamServiceServer(s grpc.ServiceRegistrar, srv StreamServiceServer) {
	s.RegisterService(&StreamService_ServiceDesc, srv)
}

func _StreamService_FetchResponse_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).FetchResponse(m, &streamServiceFetchResponseServer{stream})
}

type StreamService_FetchResponseServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type streamServiceFetchResponseServer struct {
	grpc.ServerStream
}

func (x *streamServiceFetchResponseServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// StreamService_ServiceDesc is the grpc.ServiceDesc for StreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "route.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FetchResponse",
			Handler:       _StreamService_FetchResponse_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "route.proto",
}
