// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package accountproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AccountClient is the client API for Account service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountClient interface {
	/// --- Accounts --- ///
	ListAccounts(ctx context.Context, in *ListAccountsRequest, opts ...grpc.CallOption) (*ListAccountsResponse, error)
	ReadAccount(ctx context.Context, in *ReadAccountRequest, opts ...grpc.CallOption) (*ReadAccountResponse, error)
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*UpdateAccountResponse, error)
	PageAccount(ctx context.Context, in *PageAccountRequest, opts ...grpc.CallOption) (*PageAccountResponse, error)
	/// --- Exchanges --- ///
	AddExchange(ctx context.Context, in *AddExchangeRequest, opts ...grpc.CallOption) (*AddExchangeResponse, error)
	ListExchanges(ctx context.Context, in *ListExchangesRequest, opts ...grpc.CallOption) (*ListExchangesResponse, error)
	ReadExchange(ctx context.Context, in *ReadExchangeRequest, opts ...grpc.CallOption) (*ReadExchangeResponse, error)
	ReadPrimaryExchangeByUserID(ctx context.Context, in *ReadPrimaryExchangeByUserIDRequest, opts ...grpc.CallOption) (*ReadPrimaryExchangeByUserIDResponse, error)
}

type accountClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountClient(cc grpc.ClientConnInterface) AccountClient {
	return &accountClient{cc}
}

func (c *accountClient) ListAccounts(ctx context.Context, in *ListAccountsRequest, opts ...grpc.CallOption) (*ListAccountsResponse, error) {
	out := new(ListAccountsResponse)
	err := c.cc.Invoke(ctx, "/account/ListAccounts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) ReadAccount(ctx context.Context, in *ReadAccountRequest, opts ...grpc.CallOption) (*ReadAccountResponse, error) {
	out := new(ReadAccountResponse)
	err := c.cc.Invoke(ctx, "/account/ReadAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, "/account/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*UpdateAccountResponse, error) {
	out := new(UpdateAccountResponse)
	err := c.cc.Invoke(ctx, "/account/UpdateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) PageAccount(ctx context.Context, in *PageAccountRequest, opts ...grpc.CallOption) (*PageAccountResponse, error) {
	out := new(PageAccountResponse)
	err := c.cc.Invoke(ctx, "/account/PageAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) AddExchange(ctx context.Context, in *AddExchangeRequest, opts ...grpc.CallOption) (*AddExchangeResponse, error) {
	out := new(AddExchangeResponse)
	err := c.cc.Invoke(ctx, "/account/AddExchange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) ListExchanges(ctx context.Context, in *ListExchangesRequest, opts ...grpc.CallOption) (*ListExchangesResponse, error) {
	out := new(ListExchangesResponse)
	err := c.cc.Invoke(ctx, "/account/ListExchanges", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) ReadExchange(ctx context.Context, in *ReadExchangeRequest, opts ...grpc.CallOption) (*ReadExchangeResponse, error) {
	out := new(ReadExchangeResponse)
	err := c.cc.Invoke(ctx, "/account/ReadExchange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) ReadPrimaryExchangeByUserID(ctx context.Context, in *ReadPrimaryExchangeByUserIDRequest, opts ...grpc.CallOption) (*ReadPrimaryExchangeByUserIDResponse, error) {
	out := new(ReadPrimaryExchangeByUserIDResponse)
	err := c.cc.Invoke(ctx, "/account/ReadPrimaryExchangeByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServer is the server API for Account service.
// All implementations must embed UnimplementedAccountServer
// for forward compatibility
type AccountServer interface {
	/// --- Accounts --- ///
	ListAccounts(context.Context, *ListAccountsRequest) (*ListAccountsResponse, error)
	ReadAccount(context.Context, *ReadAccountRequest) (*ReadAccountResponse, error)
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	UpdateAccount(context.Context, *UpdateAccountRequest) (*UpdateAccountResponse, error)
	PageAccount(context.Context, *PageAccountRequest) (*PageAccountResponse, error)
	/// --- Exchanges --- ///
	AddExchange(context.Context, *AddExchangeRequest) (*AddExchangeResponse, error)
	ListExchanges(context.Context, *ListExchangesRequest) (*ListExchangesResponse, error)
	ReadExchange(context.Context, *ReadExchangeRequest) (*ReadExchangeResponse, error)
	ReadPrimaryExchangeByUserID(context.Context, *ReadPrimaryExchangeByUserIDRequest) (*ReadPrimaryExchangeByUserIDResponse, error)
	mustEmbedUnimplementedAccountServer()
}

// UnimplementedAccountServer must be embedded to have forward compatible implementations.
type UnimplementedAccountServer struct {
}

func (*UnimplementedAccountServer) ListAccounts(context.Context, *ListAccountsRequest) (*ListAccountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccounts not implemented")
}
func (*UnimplementedAccountServer) ReadAccount(context.Context, *ReadAccountRequest) (*ReadAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAccount not implemented")
}
func (*UnimplementedAccountServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (*UnimplementedAccountServer) UpdateAccount(context.Context, *UpdateAccountRequest) (*UpdateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccount not implemented")
}
func (*UnimplementedAccountServer) PageAccount(context.Context, *PageAccountRequest) (*PageAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PageAccount not implemented")
}
func (*UnimplementedAccountServer) AddExchange(context.Context, *AddExchangeRequest) (*AddExchangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddExchange not implemented")
}
func (*UnimplementedAccountServer) ListExchanges(context.Context, *ListExchangesRequest) (*ListExchangesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListExchanges not implemented")
}
func (*UnimplementedAccountServer) ReadExchange(context.Context, *ReadExchangeRequest) (*ReadExchangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadExchange not implemented")
}
func (*UnimplementedAccountServer) ReadPrimaryExchangeByUserID(context.Context, *ReadPrimaryExchangeByUserIDRequest) (*ReadPrimaryExchangeByUserIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadPrimaryExchangeByUserID not implemented")
}
func (*UnimplementedAccountServer) mustEmbedUnimplementedAccountServer() {}

func RegisterAccountServer(s *grpc.Server, srv AccountServer) {
	s.RegisterService(&_Account_serviceDesc, srv)
}

func _Account_ListAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).ListAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/ListAccounts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).ListAccounts(ctx, req.(*ListAccountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_ReadAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).ReadAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/ReadAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).ReadAccount(ctx, req.(*ReadAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_UpdateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).UpdateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/UpdateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).UpdateAccount(ctx, req.(*UpdateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_PageAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).PageAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/PageAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).PageAccount(ctx, req.(*PageAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_AddExchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddExchangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).AddExchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/AddExchange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).AddExchange(ctx, req.(*AddExchangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_ListExchanges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListExchangesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).ListExchanges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/ListExchanges",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).ListExchanges(ctx, req.(*ListExchangesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_ReadExchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadExchangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).ReadExchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/ReadExchange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).ReadExchange(ctx, req.(*ReadExchangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_ReadPrimaryExchangeByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadPrimaryExchangeByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).ReadPrimaryExchangeByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/account/ReadPrimaryExchangeByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).ReadPrimaryExchangeByUserID(ctx, req.(*ReadPrimaryExchangeByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Account_serviceDesc = grpc.ServiceDesc{
	ServiceName: "account",
	HandlerType: (*AccountServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAccounts",
			Handler:    _Account_ListAccounts_Handler,
		},
		{
			MethodName: "ReadAccount",
			Handler:    _Account_ReadAccount_Handler,
		},
		{
			MethodName: "CreateAccount",
			Handler:    _Account_CreateAccount_Handler,
		},
		{
			MethodName: "UpdateAccount",
			Handler:    _Account_UpdateAccount_Handler,
		},
		{
			MethodName: "PageAccount",
			Handler:    _Account_PageAccount_Handler,
		},
		{
			MethodName: "AddExchange",
			Handler:    _Account_AddExchange_Handler,
		},
		{
			MethodName: "ListExchanges",
			Handler:    _Account_ListExchanges_Handler,
		},
		{
			MethodName: "ReadExchange",
			Handler:    _Account_ReadExchange_Handler,
		},
		{
			MethodName: "ReadPrimaryExchangeByUserID",
			Handler:    _Account_ReadPrimaryExchangeByUserID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "s.account/proto/account.proto",
}
