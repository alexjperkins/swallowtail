// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package bitfinexproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BitfinexClient is the client API for Bitfinex service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BitfinexClient interface {
	GetBitfinexStatus(ctx context.Context, in *GetBitfinexStatusRequest, opts ...grpc.CallOption) (*GetBitfinexStatusResponse, error)
	GetBitfinexFundingRates(ctx context.Context, in *GetBitfinexFundingRatesRequest, opts ...grpc.CallOption) (*GetBitfinexFundingRatesResponse, error)
}

type bitfinexClient struct {
	cc grpc.ClientConnInterface
}

func NewBitfinexClient(cc grpc.ClientConnInterface) BitfinexClient {
	return &bitfinexClient{cc}
}

func (c *bitfinexClient) GetBitfinexStatus(ctx context.Context, in *GetBitfinexStatusRequest, opts ...grpc.CallOption) (*GetBitfinexStatusResponse, error) {
	out := new(GetBitfinexStatusResponse)
	err := c.cc.Invoke(ctx, "/bitfinex/GetBitfinexStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bitfinexClient) GetBitfinexFundingRates(ctx context.Context, in *GetBitfinexFundingRatesRequest, opts ...grpc.CallOption) (*GetBitfinexFundingRatesResponse, error) {
	out := new(GetBitfinexFundingRatesResponse)
	err := c.cc.Invoke(ctx, "/bitfinex/GetBitfinexFundingRates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BitfinexServer is the server API for Bitfinex service.
// All implementations must embed UnimplementedBitfinexServer
// for forward compatibility
type BitfinexServer interface {
	GetBitfinexStatus(context.Context, *GetBitfinexStatusRequest) (*GetBitfinexStatusResponse, error)
	GetBitfinexFundingRates(context.Context, *GetBitfinexFundingRatesRequest) (*GetBitfinexFundingRatesResponse, error)
	mustEmbedUnimplementedBitfinexServer()
}

// UnimplementedBitfinexServer must be embedded to have forward compatible implementations.
type UnimplementedBitfinexServer struct {
}

func (*UnimplementedBitfinexServer) GetBitfinexStatus(context.Context, *GetBitfinexStatusRequest) (*GetBitfinexStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBitfinexStatus not implemented")
}
func (*UnimplementedBitfinexServer) GetBitfinexFundingRates(context.Context, *GetBitfinexFundingRatesRequest) (*GetBitfinexFundingRatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBitfinexFundingRates not implemented")
}
func (*UnimplementedBitfinexServer) mustEmbedUnimplementedBitfinexServer() {}

func RegisterBitfinexServer(s *grpc.Server, srv BitfinexServer) {
	s.RegisterService(&_Bitfinex_serviceDesc, srv)
}

func _Bitfinex_GetBitfinexStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBitfinexStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitfinexServer).GetBitfinexStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bitfinex/GetBitfinexStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitfinexServer).GetBitfinexStatus(ctx, req.(*GetBitfinexStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bitfinex_GetBitfinexFundingRates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBitfinexFundingRatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitfinexServer).GetBitfinexFundingRates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bitfinex/GetBitfinexFundingRates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitfinexServer).GetBitfinexFundingRates(ctx, req.(*GetBitfinexFundingRatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Bitfinex_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bitfinex",
	HandlerType: (*BitfinexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBitfinexStatus",
			Handler:    _Bitfinex_GetBitfinexStatus_Handler,
		},
		{
			MethodName: "GetBitfinexFundingRates",
			Handler:    _Bitfinex_GetBitfinexFundingRates_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "s.bitfinex/proto/bitfinex.proto",
}
