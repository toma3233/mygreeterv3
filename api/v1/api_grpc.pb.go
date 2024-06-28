// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: api.proto

package v1

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
	MyGreeter_SayHello_FullMethodName             = "/MyGreeter/SayHello"
	MyGreeter_CreateResourceGroup_FullMethodName  = "/MyGreeter/CreateResourceGroup"
	MyGreeter_ReadResourceGroup_FullMethodName    = "/MyGreeter/ReadResourceGroup"
	MyGreeter_DeleteResourceGroup_FullMethodName  = "/MyGreeter/DeleteResourceGroup"
	MyGreeter_UpdateResourceGroup_FullMethodName  = "/MyGreeter/UpdateResourceGroup"
	MyGreeter_ListResourceGroups_FullMethodName   = "/MyGreeter/ListResourceGroups"
	MyGreeter_CreateStorageAccount_FullMethodName = "/MyGreeter/CreateStorageAccount"
	MyGreeter_ReadStorageAccount_FullMethodName   = "/MyGreeter/ReadStorageAccount"
	MyGreeter_DeleteStorageAccount_FullMethodName = "/MyGreeter/DeleteStorageAccount"
	MyGreeter_UpdateStorageAccount_FullMethodName = "/MyGreeter/UpdateStorageAccount"
	MyGreeter_ListStorageAccounts_FullMethodName  = "/MyGreeter/ListStorageAccounts"
	MyGreeter_CreateBlobContainer_FullMethodName  = "/MyGreeter/CreateBlobContainer"
	MyGreeter_ReadBlobContainer_FullMethodName    = "/MyGreeter/ReadBlobContainer"
	MyGreeter_DeleteBlobContainer_FullMethodName  = "/MyGreeter/DeleteBlobContainer"
	MyGreeter_UpdateBlobContainer_FullMethodName  = "/MyGreeter/UpdateBlobContainer"
	MyGreeter_ListBlobContainers_FullMethodName   = "/MyGreeter/ListBlobContainers"
)

// MyGreeterClient is the client API for MyGreeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MyGreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	// Creates a resource group
	CreateResourceGroup(ctx context.Context, in *CreateResourceGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Reads a resource group
	ReadResourceGroup(ctx context.Context, in *ReadResourceGroupRequest, opts ...grpc.CallOption) (*ReadResourceGroupResponse, error)
	// Deletes a resource group
	DeleteResourceGroup(ctx context.Context, in *DeleteResourceGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Updates a resource group
	UpdateResourceGroup(ctx context.Context, in *UpdateResourceGroupRequest, opts ...grpc.CallOption) (*UpdateResourceGroupResponse, error)
	// Lists all resource groups
	ListResourceGroups(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListResourceGroupResponse, error)
	// Creates a storage account
	CreateStorageAccount(ctx context.Context, in *CreateStorageAccountRequest, opts ...grpc.CallOption) (*CreateStorageAccountResponse, error)
	// Reads a storage account
	ReadStorageAccount(ctx context.Context, in *ReadStorageAccountRequest, opts ...grpc.CallOption) (*ReadStorageAccountResponse, error)
	// Deletes a storage account
	DeleteStorageAccount(ctx context.Context, in *DeleteStorageAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Updates a storage account
	UpdateStorageAccount(ctx context.Context, in *UpdateStorageAccountRequest, opts ...grpc.CallOption) (*UpdateStorageAccountResponse, error)
	// Lists all storage accounts
	ListStorageAccounts(ctx context.Context, in *ListStorageAccountRequest, opts ...grpc.CallOption) (*ListStorageAccountResponse, error)
	// Creates a blob storage container
	CreateBlobContainer(ctx context.Context, in *CreateBlobContainerRequest, opts ...grpc.CallOption) (*CreateBlobContainerResponse, error)
	// Reads a blob storage container
	ReadBlobContainer(ctx context.Context, in *ReadBlobContainerRequest, opts ...grpc.CallOption) (*ReadBlobContainerResponse, error)
	// Deletes a blob storage container
	DeleteBlobContainer(ctx context.Context, in *DeleteBlobContainerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Updates a blob storage container
	UpdateBlobContainer(ctx context.Context, in *UpdateBlobContainerRequest, opts ...grpc.CallOption) (*UpdateBlobContainerResponse, error)
	// Lists all blob storage containers
	ListBlobContainers(ctx context.Context, in *ListBlobContainersRequest, opts ...grpc.CallOption) (*ListBlobContainersResponse, error)
}

type myGreeterClient struct {
	cc grpc.ClientConnInterface
}

func NewMyGreeterClient(cc grpc.ClientConnInterface) MyGreeterClient {
	return &myGreeterClient{cc}
}

func (c *myGreeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, MyGreeter_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) CreateResourceGroup(ctx context.Context, in *CreateResourceGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MyGreeter_CreateResourceGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) ReadResourceGroup(ctx context.Context, in *ReadResourceGroupRequest, opts ...grpc.CallOption) (*ReadResourceGroupResponse, error) {
	out := new(ReadResourceGroupResponse)
	err := c.cc.Invoke(ctx, MyGreeter_ReadResourceGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) DeleteResourceGroup(ctx context.Context, in *DeleteResourceGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MyGreeter_DeleteResourceGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) UpdateResourceGroup(ctx context.Context, in *UpdateResourceGroupRequest, opts ...grpc.CallOption) (*UpdateResourceGroupResponse, error) {
	out := new(UpdateResourceGroupResponse)
	err := c.cc.Invoke(ctx, MyGreeter_UpdateResourceGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) ListResourceGroups(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListResourceGroupResponse, error) {
	out := new(ListResourceGroupResponse)
	err := c.cc.Invoke(ctx, MyGreeter_ListResourceGroups_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) CreateStorageAccount(ctx context.Context, in *CreateStorageAccountRequest, opts ...grpc.CallOption) (*CreateStorageAccountResponse, error) {
	out := new(CreateStorageAccountResponse)
	err := c.cc.Invoke(ctx, MyGreeter_CreateStorageAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) ReadStorageAccount(ctx context.Context, in *ReadStorageAccountRequest, opts ...grpc.CallOption) (*ReadStorageAccountResponse, error) {
	out := new(ReadStorageAccountResponse)
	err := c.cc.Invoke(ctx, MyGreeter_ReadStorageAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) DeleteStorageAccount(ctx context.Context, in *DeleteStorageAccountRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MyGreeter_DeleteStorageAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) UpdateStorageAccount(ctx context.Context, in *UpdateStorageAccountRequest, opts ...grpc.CallOption) (*UpdateStorageAccountResponse, error) {
	out := new(UpdateStorageAccountResponse)
	err := c.cc.Invoke(ctx, MyGreeter_UpdateStorageAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) ListStorageAccounts(ctx context.Context, in *ListStorageAccountRequest, opts ...grpc.CallOption) (*ListStorageAccountResponse, error) {
	out := new(ListStorageAccountResponse)
	err := c.cc.Invoke(ctx, MyGreeter_ListStorageAccounts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) CreateBlobContainer(ctx context.Context, in *CreateBlobContainerRequest, opts ...grpc.CallOption) (*CreateBlobContainerResponse, error) {
	out := new(CreateBlobContainerResponse)
	err := c.cc.Invoke(ctx, MyGreeter_CreateBlobContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) ReadBlobContainer(ctx context.Context, in *ReadBlobContainerRequest, opts ...grpc.CallOption) (*ReadBlobContainerResponse, error) {
	out := new(ReadBlobContainerResponse)
	err := c.cc.Invoke(ctx, MyGreeter_ReadBlobContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) DeleteBlobContainer(ctx context.Context, in *DeleteBlobContainerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MyGreeter_DeleteBlobContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) UpdateBlobContainer(ctx context.Context, in *UpdateBlobContainerRequest, opts ...grpc.CallOption) (*UpdateBlobContainerResponse, error) {
	out := new(UpdateBlobContainerResponse)
	err := c.cc.Invoke(ctx, MyGreeter_UpdateBlobContainer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *myGreeterClient) ListBlobContainers(ctx context.Context, in *ListBlobContainersRequest, opts ...grpc.CallOption) (*ListBlobContainersResponse, error) {
	out := new(ListBlobContainersResponse)
	err := c.cc.Invoke(ctx, MyGreeter_ListBlobContainers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MyGreeterServer is the server API for MyGreeter service.
// All implementations must embed UnimplementedMyGreeterServer
// for forward compatibility
type MyGreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	// Creates a resource group
	CreateResourceGroup(context.Context, *CreateResourceGroupRequest) (*emptypb.Empty, error)
	// Reads a resource group
	ReadResourceGroup(context.Context, *ReadResourceGroupRequest) (*ReadResourceGroupResponse, error)
	// Deletes a resource group
	DeleteResourceGroup(context.Context, *DeleteResourceGroupRequest) (*emptypb.Empty, error)
	// Updates a resource group
	UpdateResourceGroup(context.Context, *UpdateResourceGroupRequest) (*UpdateResourceGroupResponse, error)
	// Lists all resource groups
	ListResourceGroups(context.Context, *emptypb.Empty) (*ListResourceGroupResponse, error)
	// Creates a storage account
	CreateStorageAccount(context.Context, *CreateStorageAccountRequest) (*CreateStorageAccountResponse, error)
	// Reads a storage account
	ReadStorageAccount(context.Context, *ReadStorageAccountRequest) (*ReadStorageAccountResponse, error)
	// Deletes a storage account
	DeleteStorageAccount(context.Context, *DeleteStorageAccountRequest) (*emptypb.Empty, error)
	// Updates a storage account
	UpdateStorageAccount(context.Context, *UpdateStorageAccountRequest) (*UpdateStorageAccountResponse, error)
	// Lists all storage accounts
	ListStorageAccounts(context.Context, *ListStorageAccountRequest) (*ListStorageAccountResponse, error)
	// Creates a blob storage container
	CreateBlobContainer(context.Context, *CreateBlobContainerRequest) (*CreateBlobContainerResponse, error)
	// Reads a blob storage container
	ReadBlobContainer(context.Context, *ReadBlobContainerRequest) (*ReadBlobContainerResponse, error)
	// Deletes a blob storage container
	DeleteBlobContainer(context.Context, *DeleteBlobContainerRequest) (*emptypb.Empty, error)
	// Updates a blob storage container
	UpdateBlobContainer(context.Context, *UpdateBlobContainerRequest) (*UpdateBlobContainerResponse, error)
	// Lists all blob storage containers
	ListBlobContainers(context.Context, *ListBlobContainersRequest) (*ListBlobContainersResponse, error)
	mustEmbedUnimplementedMyGreeterServer()
}

// UnimplementedMyGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedMyGreeterServer struct {
}

func (UnimplementedMyGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedMyGreeterServer) CreateResourceGroup(context.Context, *CreateResourceGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateResourceGroup not implemented")
}
func (UnimplementedMyGreeterServer) ReadResourceGroup(context.Context, *ReadResourceGroupRequest) (*ReadResourceGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadResourceGroup not implemented")
}
func (UnimplementedMyGreeterServer) DeleteResourceGroup(context.Context, *DeleteResourceGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteResourceGroup not implemented")
}
func (UnimplementedMyGreeterServer) UpdateResourceGroup(context.Context, *UpdateResourceGroupRequest) (*UpdateResourceGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResourceGroup not implemented")
}
func (UnimplementedMyGreeterServer) ListResourceGroups(context.Context, *emptypb.Empty) (*ListResourceGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListResourceGroups not implemented")
}
func (UnimplementedMyGreeterServer) CreateStorageAccount(context.Context, *CreateStorageAccountRequest) (*CreateStorageAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStorageAccount not implemented")
}
func (UnimplementedMyGreeterServer) ReadStorageAccount(context.Context, *ReadStorageAccountRequest) (*ReadStorageAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadStorageAccount not implemented")
}
func (UnimplementedMyGreeterServer) DeleteStorageAccount(context.Context, *DeleteStorageAccountRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStorageAccount not implemented")
}
func (UnimplementedMyGreeterServer) UpdateStorageAccount(context.Context, *UpdateStorageAccountRequest) (*UpdateStorageAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStorageAccount not implemented")
}
func (UnimplementedMyGreeterServer) ListStorageAccounts(context.Context, *ListStorageAccountRequest) (*ListStorageAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStorageAccounts not implemented")
}
func (UnimplementedMyGreeterServer) CreateBlobContainer(context.Context, *CreateBlobContainerRequest) (*CreateBlobContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBlobContainer not implemented")
}
func (UnimplementedMyGreeterServer) ReadBlobContainer(context.Context, *ReadBlobContainerRequest) (*ReadBlobContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadBlobContainer not implemented")
}
func (UnimplementedMyGreeterServer) DeleteBlobContainer(context.Context, *DeleteBlobContainerRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBlobContainer not implemented")
}
func (UnimplementedMyGreeterServer) UpdateBlobContainer(context.Context, *UpdateBlobContainerRequest) (*UpdateBlobContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBlobContainer not implemented")
}
func (UnimplementedMyGreeterServer) ListBlobContainers(context.Context, *ListBlobContainersRequest) (*ListBlobContainersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBlobContainers not implemented")
}
func (UnimplementedMyGreeterServer) mustEmbedUnimplementedMyGreeterServer() {}

// UnsafeMyGreeterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MyGreeterServer will
// result in compilation errors.
type UnsafeMyGreeterServer interface {
	mustEmbedUnimplementedMyGreeterServer()
}

func RegisterMyGreeterServer(s grpc.ServiceRegistrar, srv MyGreeterServer) {
	s.RegisterService(&MyGreeter_ServiceDesc, srv)
}

func _MyGreeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_CreateResourceGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateResourceGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).CreateResourceGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_CreateResourceGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).CreateResourceGroup(ctx, req.(*CreateResourceGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_ReadResourceGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadResourceGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).ReadResourceGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_ReadResourceGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).ReadResourceGroup(ctx, req.(*ReadResourceGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_DeleteResourceGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteResourceGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).DeleteResourceGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_DeleteResourceGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).DeleteResourceGroup(ctx, req.(*DeleteResourceGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_UpdateResourceGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).UpdateResourceGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_UpdateResourceGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).UpdateResourceGroup(ctx, req.(*UpdateResourceGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_ListResourceGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).ListResourceGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_ListResourceGroups_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).ListResourceGroups(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_CreateStorageAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStorageAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).CreateStorageAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_CreateStorageAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MyGreeterServer).CreateStorageAccount(ctx, req.(*CreateStorageAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MyGreeter_ReadStorageAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadStorageAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MyGreeterServer).ReadStorageAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MyGreeter_ReadStorageAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error)