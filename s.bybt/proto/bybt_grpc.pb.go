// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package bybtproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BybtClient is the client API for Bybt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BybtClient interface {
	GetExchangeFundingRates(ctx context.Context, in *GetExchangeFundingRatesRequest, opts ...grpc.CallOption) (*GetExchangeFundingRatesResponse, error)
}

type bybtClient struct {
	cc grpc.ClientConnInterface
}

func NewBybtClient(cc grpc.ClientConnInterface) BybtClient {
	return &bybtClient{cc}
}

func (c *bybtClient) GetExchangeFundingRates(ctx context.Context, in *GetExchangeFundingRatesRequest, opts ...grpc.CallOption) (*GetExchangeFundingRatesResponse, error) {
	out := new(GetExchangeFundingRatesResponse)
	err := c.cc.Invoke(ctx, "/bybt/GetExchangeFundingRates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BybtServer is the server API for Bybt service.
// All implementations must embed UnimplementedBybtServer
// for forward compatibility
type BybtServer interface {
	GetExchangeFundingRates(context.Context, *GetExchangeFundingRatesRequest) (*GetExchangeFundingRatesResponse, error)
	mustEmbedUnimplementedBybtServer()
}

// UnimplementedBybtServer must be embedded to have forward compatible implementations.
type UnimplementedBybtServer struct {
}

func (*UnimplementedBybtServer) GetExchangeFundingRates(context.Context, *GetExchangeFundingRatesRequest) (*GetExchangeFundingRatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExchangeFundingRates not implemented")
}
func (*UnimplementedBybtServer) mustEmbedUnimplementedBybtServer() {}

func RegisterBybtServer(s *grpc.Server, srv BybtServer) {
	s.RegisterService(&_Bybt_serviceDesc, srv)
}

func _Bybt_GetExchangeFundingRates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExchangeFundingRatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BybtServer).GetExchangeFundingRates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bybt/GetExchangeFundingRates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BybtServer).GetExchangeFundingRates(ctx, req.(*GetExchangeFundingRatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Bybt_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bybt",
	HandlerType: (*BybtServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetExchangeFundingRates",
			Handler:    _Bybt_GetExchangeFundingRates_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "s.bybt/proto/bybt.proto",
}
