// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package marketdataproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MarketdataClient is the client API for Marketdata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MarketdataClient interface {
	PublishLatestPriceInformation(ctx context.Context, in *PublishLatestPriceInformationRequest, opts ...grpc.CallOption) (*PublishLatestPriceInformationResponse, error)
	PublishVolatilityInformation(ctx context.Context, in *PublishVolatilityInformationRequest, opts ...grpc.CallOption) (*PublishVolatilityInformationResponse, error)
	PublishATHInformation(ctx context.Context, in *PublishATHInformationRequest, opts ...grpc.CallOption) (*PublishATHInformationResponse, error)
	PublishFundingRatesInformation(ctx context.Context, in *PublishFundingRatesInformationRequest, opts ...grpc.CallOption) (*PublishFundingRatesInformationResponse, error)
	PublishSolanaNFTPriceInformation(ctx context.Context, in *PublishSolanaNFTPriceInformationRequest, opts ...grpc.CallOption) (*PublishSolanaNFTPriceInformationResponse, error)
}

type marketdataClient struct {
	cc grpc.ClientConnInterface
}

func NewMarketdataClient(cc grpc.ClientConnInterface) MarketdataClient {
	return &marketdataClient{cc}
}

func (c *marketdataClient) PublishLatestPriceInformation(ctx context.Context, in *PublishLatestPriceInformationRequest, opts ...grpc.CallOption) (*PublishLatestPriceInformationResponse, error) {
	out := new(PublishLatestPriceInformationResponse)
	err := c.cc.Invoke(ctx, "/marketdata/PublishLatestPriceInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketdataClient) PublishVolatilityInformation(ctx context.Context, in *PublishVolatilityInformationRequest, opts ...grpc.CallOption) (*PublishVolatilityInformationResponse, error) {
	out := new(PublishVolatilityInformationResponse)
	err := c.cc.Invoke(ctx, "/marketdata/PublishVolatilityInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketdataClient) PublishATHInformation(ctx context.Context, in *PublishATHInformationRequest, opts ...grpc.CallOption) (*PublishATHInformationResponse, error) {
	out := new(PublishATHInformationResponse)
	err := c.cc.Invoke(ctx, "/marketdata/PublishATHInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketdataClient) PublishFundingRatesInformation(ctx context.Context, in *PublishFundingRatesInformationRequest, opts ...grpc.CallOption) (*PublishFundingRatesInformationResponse, error) {
	out := new(PublishFundingRatesInformationResponse)
	err := c.cc.Invoke(ctx, "/marketdata/PublishFundingRatesInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketdataClient) PublishSolanaNFTPriceInformation(ctx context.Context, in *PublishSolanaNFTPriceInformationRequest, opts ...grpc.CallOption) (*PublishSolanaNFTPriceInformationResponse, error) {
	out := new(PublishSolanaNFTPriceInformationResponse)
	err := c.cc.Invoke(ctx, "/marketdata/PublishSolanaNFTPriceInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MarketdataServer is the server API for Marketdata service.
// All implementations must embed UnimplementedMarketdataServer
// for forward compatibility
type MarketdataServer interface {
	PublishLatestPriceInformation(context.Context, *PublishLatestPriceInformationRequest) (*PublishLatestPriceInformationResponse, error)
	PublishVolatilityInformation(context.Context, *PublishVolatilityInformationRequest) (*PublishVolatilityInformationResponse, error)
	PublishATHInformation(context.Context, *PublishATHInformationRequest) (*PublishATHInformationResponse, error)
	PublishFundingRatesInformation(context.Context, *PublishFundingRatesInformationRequest) (*PublishFundingRatesInformationResponse, error)
	PublishSolanaNFTPriceInformation(context.Context, *PublishSolanaNFTPriceInformationRequest) (*PublishSolanaNFTPriceInformationResponse, error)
	mustEmbedUnimplementedMarketdataServer()
}

// UnimplementedMarketdataServer must be embedded to have forward compatible implementations.
type UnimplementedMarketdataServer struct {
}

func (*UnimplementedMarketdataServer) PublishLatestPriceInformation(context.Context, *PublishLatestPriceInformationRequest) (*PublishLatestPriceInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishLatestPriceInformation not implemented")
}
func (*UnimplementedMarketdataServer) PublishVolatilityInformation(context.Context, *PublishVolatilityInformationRequest) (*PublishVolatilityInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishVolatilityInformation not implemented")
}
func (*UnimplementedMarketdataServer) PublishATHInformation(context.Context, *PublishATHInformationRequest) (*PublishATHInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishATHInformation not implemented")
}
func (*UnimplementedMarketdataServer) PublishFundingRatesInformation(context.Context, *PublishFundingRatesInformationRequest) (*PublishFundingRatesInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishFundingRatesInformation not implemented")
}
func (*UnimplementedMarketdataServer) PublishSolanaNFTPriceInformation(context.Context, *PublishSolanaNFTPriceInformationRequest) (*PublishSolanaNFTPriceInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishSolanaNFTPriceInformation not implemented")
}
func (*UnimplementedMarketdataServer) mustEmbedUnimplementedMarketdataServer() {}

func RegisterMarketdataServer(s *grpc.Server, srv MarketdataServer) {
	s.RegisterService(&_Marketdata_serviceDesc, srv)
}

func _Marketdata_PublishLatestPriceInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishLatestPriceInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketdataServer).PublishLatestPriceInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketdata/PublishLatestPriceInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketdataServer).PublishLatestPriceInformation(ctx, req.(*PublishLatestPriceInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketdata_PublishVolatilityInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishVolatilityInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketdataServer).PublishVolatilityInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketdata/PublishVolatilityInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketdataServer).PublishVolatilityInformation(ctx, req.(*PublishVolatilityInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketdata_PublishATHInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishATHInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketdataServer).PublishATHInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketdata/PublishATHInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketdataServer).PublishATHInformation(ctx, req.(*PublishATHInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketdata_PublishFundingRatesInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishFundingRatesInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketdataServer).PublishFundingRatesInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketdata/PublishFundingRatesInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketdataServer).PublishFundingRatesInformation(ctx, req.(*PublishFundingRatesInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketdata_PublishSolanaNFTPriceInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishSolanaNFTPriceInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketdataServer).PublishSolanaNFTPriceInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketdata/PublishSolanaNFTPriceInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketdataServer).PublishSolanaNFTPriceInformation(ctx, req.(*PublishSolanaNFTPriceInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Marketdata_serviceDesc = grpc.ServiceDesc{
	ServiceName: "marketdata",
	HandlerType: (*MarketdataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishLatestPriceInformation",
			Handler:    _Marketdata_PublishLatestPriceInformation_Handler,
		},
		{
			MethodName: "PublishVolatilityInformation",
			Handler:    _Marketdata_PublishVolatilityInformation_Handler,
		},
		{
			MethodName: "PublishATHInformation",
			Handler:    _Marketdata_PublishATHInformation_Handler,
		},
		{
			MethodName: "PublishFundingRatesInformation",
			Handler:    _Marketdata_PublishFundingRatesInformation_Handler,
		},
		{
			MethodName: "PublishSolanaNFTPriceInformation",
			Handler:    _Marketdata_PublishSolanaNFTPriceInformation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "s.market-data/proto/marketdata.proto",
}
