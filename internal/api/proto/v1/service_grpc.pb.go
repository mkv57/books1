// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: internal/api/proto/v1/service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BookAPI_AddBook_FullMethodName = "/api.proto.v1.BookAPI/AddBook"
)

// BookAPIClient is the client API for BookAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookAPIClient interface {
	AddBook(ctx context.Context, in *AddBookRequest, opts ...grpc.CallOption) (*AddBookResponse, error)
}

type bookAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewBookAPIClient(cc grpc.ClientConnInterface) BookAPIClient {
	return &bookAPIClient{cc}
}

func (c *bookAPIClient) AddBook(ctx context.Context, in *AddBookRequest, opts ...grpc.CallOption) (*AddBookResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddBookResponse)
	err := c.cc.Invoke(ctx, BookAPI_AddBook_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookAPIServer is the server API for BookAPI service.
// All implementations should embed UnimplementedBookAPIServer
// for forward compatibility.
type BookAPIServer interface {
	AddBook(context.Context, *AddBookRequest) (*AddBookResponse, error)
}

// UnimplementedBookAPIServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBookAPIServer struct{}

func (UnimplementedBookAPIServer) AddBook(context.Context, *AddBookRequest) (*AddBookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBook not implemented")
}
func (UnimplementedBookAPIServer) testEmbeddedByValue() {}

// UnsafeBookAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookAPIServer will
// result in compilation errors.
type UnsafeBookAPIServer interface {
	mustEmbedUnimplementedBookAPIServer()
}

func RegisterBookAPIServer(s grpc.ServiceRegistrar, srv BookAPIServer) {
	// If the following call pancis, it indicates UnimplementedBookAPIServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BookAPI_ServiceDesc, srv)
}

func _BookAPI_AddBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookAPIServer).AddBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookAPI_AddBook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookAPIServer).AddBook(ctx, req.(*AddBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BookAPI_ServiceDesc is the grpc.ServiceDesc for BookAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.proto.v1.BookAPI",
	HandlerType: (*BookAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddBook",
			Handler:    _BookAPI_AddBook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/api/proto/v1/service.proto",
}
