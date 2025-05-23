// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/dfs.proto

package __

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
	MasterTracker_GetUploadNode_FullMethodName     = "/dfs.MasterTracker/GetUploadNode"
	MasterTracker_RegisterUpload_FullMethodName    = "/dfs.MasterTracker/RegisterUpload"
	MasterTracker_GetDownloadNodes_FullMethodName  = "/dfs.MasterTracker/GetDownloadNodes"
	MasterTracker_Heartbeat_FullMethodName         = "/dfs.MasterTracker/Heartbeat"
	MasterTracker_NotifyReplication_FullMethodName = "/dfs.MasterTracker/NotifyReplication"
)

// MasterTrackerClient is the client API for MasterTracker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Master Tracker service definition
type MasterTrackerClient interface {
	// Get Data Keeper node for uploading a file
	GetUploadNode(ctx context.Context, in *GetUploadNodeRequest, opts ...grpc.CallOption) (*GetUploadNodeResponse, error)
	// Register a successful upload with the master tracker
	RegisterUpload(ctx context.Context, in *RegisterUploadRequest, opts ...grpc.CallOption) (*RegisterUploadResponse, error)
	// Get available nodes for downloading a file
	GetDownloadNodes(ctx context.Context, in *GetDownloadNodesRequest, opts ...grpc.CallOption) (*GetDownloadNodesResponse, error)
	// Heartbeat from Data Keeper to Master Tracker
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	// Notify Data Keepers to replicate a file
	NotifyReplication(ctx context.Context, in *ReplicationRequest, opts ...grpc.CallOption) (*ReplicationResponse, error)
}

type masterTrackerClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterTrackerClient(cc grpc.ClientConnInterface) MasterTrackerClient {
	return &masterTrackerClient{cc}
}

func (c *masterTrackerClient) GetUploadNode(ctx context.Context, in *GetUploadNodeRequest, opts ...grpc.CallOption) (*GetUploadNodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUploadNodeResponse)
	err := c.cc.Invoke(ctx, MasterTracker_GetUploadNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) RegisterUpload(ctx context.Context, in *RegisterUploadRequest, opts ...grpc.CallOption) (*RegisterUploadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterUploadResponse)
	err := c.cc.Invoke(ctx, MasterTracker_RegisterUpload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) GetDownloadNodes(ctx context.Context, in *GetDownloadNodesRequest, opts ...grpc.CallOption) (*GetDownloadNodesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDownloadNodesResponse)
	err := c.cc.Invoke(ctx, MasterTracker_GetDownloadNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, MasterTracker_Heartbeat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterTrackerClient) NotifyReplication(ctx context.Context, in *ReplicationRequest, opts ...grpc.CallOption) (*ReplicationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReplicationResponse)
	err := c.cc.Invoke(ctx, MasterTracker_NotifyReplication_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterTrackerServer is the server API for MasterTracker service.
// All implementations must embed UnimplementedMasterTrackerServer
// for forward compatibility.
//
// Master Tracker service definition
type MasterTrackerServer interface {
	// Get Data Keeper node for uploading a file
	GetUploadNode(context.Context, *GetUploadNodeRequest) (*GetUploadNodeResponse, error)
	// Register a successful upload with the master tracker
	RegisterUpload(context.Context, *RegisterUploadRequest) (*RegisterUploadResponse, error)
	// Get available nodes for downloading a file
	GetDownloadNodes(context.Context, *GetDownloadNodesRequest) (*GetDownloadNodesResponse, error)
	// Heartbeat from Data Keeper to Master Tracker
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	// Notify Data Keepers to replicate a file
	NotifyReplication(context.Context, *ReplicationRequest) (*ReplicationResponse, error)
	mustEmbedUnimplementedMasterTrackerServer()
}

// UnimplementedMasterTrackerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMasterTrackerServer struct{}

func (UnimplementedMasterTrackerServer) GetUploadNode(context.Context, *GetUploadNodeRequest) (*GetUploadNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUploadNode not implemented")
}
func (UnimplementedMasterTrackerServer) RegisterUpload(context.Context, *RegisterUploadRequest) (*RegisterUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUpload not implemented")
}
func (UnimplementedMasterTrackerServer) GetDownloadNodes(context.Context, *GetDownloadNodesRequest) (*GetDownloadNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDownloadNodes not implemented")
}
func (UnimplementedMasterTrackerServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedMasterTrackerServer) NotifyReplication(context.Context, *ReplicationRequest) (*ReplicationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyReplication not implemented")
}
func (UnimplementedMasterTrackerServer) mustEmbedUnimplementedMasterTrackerServer() {}
func (UnimplementedMasterTrackerServer) testEmbeddedByValue()                       {}

// UnsafeMasterTrackerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterTrackerServer will
// result in compilation errors.
type UnsafeMasterTrackerServer interface {
	mustEmbedUnimplementedMasterTrackerServer()
}

func RegisterMasterTrackerServer(s grpc.ServiceRegistrar, srv MasterTrackerServer) {
	// If the following call pancis, it indicates UnimplementedMasterTrackerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MasterTracker_ServiceDesc, srv)
}

func _MasterTracker_GetUploadNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUploadNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).GetUploadNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_GetUploadNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).GetUploadNode(ctx, req.(*GetUploadNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_RegisterUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).RegisterUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_RegisterUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).RegisterUpload(ctx, req.(*RegisterUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_GetDownloadNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDownloadNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).GetDownloadNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_GetDownloadNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).GetDownloadNodes(ctx, req.(*GetDownloadNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_Heartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterTracker_NotifyReplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterTrackerServer).NotifyReplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterTracker_NotifyReplication_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterTrackerServer).NotifyReplication(ctx, req.(*ReplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MasterTracker_ServiceDesc is the grpc.ServiceDesc for MasterTracker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MasterTracker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dfs.MasterTracker",
	HandlerType: (*MasterTrackerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUploadNode",
			Handler:    _MasterTracker_GetUploadNode_Handler,
		},
		{
			MethodName: "RegisterUpload",
			Handler:    _MasterTracker_RegisterUpload_Handler,
		},
		{
			MethodName: "GetDownloadNodes",
			Handler:    _MasterTracker_GetDownloadNodes_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _MasterTracker_Heartbeat_Handler,
		},
		{
			MethodName: "NotifyReplication",
			Handler:    _MasterTracker_NotifyReplication_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/dfs.proto",
}

const (
	DataKeeper_InitiateUpload_FullMethodName = "/dfs.DataKeeper/InitiateUpload"
	DataKeeper_ReplicateFile_FullMethodName  = "/dfs.DataKeeper/ReplicateFile"
)

// DataKeeperClient is the client API for DataKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Data Keeper service definition
type DataKeeperClient interface {
	// Initiate file transfer from client to data keeper
	InitiateUpload(ctx context.Context, in *InitiateUploadRequest, opts ...grpc.CallOption) (*InitiateUploadResponse, error)
	// Replicate a file from another data keeper
	ReplicateFile(ctx context.Context, in *ReplicateFileRequest, opts ...grpc.CallOption) (*ReplicateFileResponse, error)
}

type dataKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewDataKeeperClient(cc grpc.ClientConnInterface) DataKeeperClient {
	return &dataKeeperClient{cc}
}

func (c *dataKeeperClient) InitiateUpload(ctx context.Context, in *InitiateUploadRequest, opts ...grpc.CallOption) (*InitiateUploadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InitiateUploadResponse)
	err := c.cc.Invoke(ctx, DataKeeper_InitiateUpload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataKeeperClient) ReplicateFile(ctx context.Context, in *ReplicateFileRequest, opts ...grpc.CallOption) (*ReplicateFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReplicateFileResponse)
	err := c.cc.Invoke(ctx, DataKeeper_ReplicateFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataKeeperServer is the server API for DataKeeper service.
// All implementations must embed UnimplementedDataKeeperServer
// for forward compatibility.
//
// Data Keeper service definition
type DataKeeperServer interface {
	// Initiate file transfer from client to data keeper
	InitiateUpload(context.Context, *InitiateUploadRequest) (*InitiateUploadResponse, error)
	// Replicate a file from another data keeper
	ReplicateFile(context.Context, *ReplicateFileRequest) (*ReplicateFileResponse, error)
	mustEmbedUnimplementedDataKeeperServer()
}

// UnimplementedDataKeeperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataKeeperServer struct{}

func (UnimplementedDataKeeperServer) InitiateUpload(context.Context, *InitiateUploadRequest) (*InitiateUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitiateUpload not implemented")
}
func (UnimplementedDataKeeperServer) ReplicateFile(context.Context, *ReplicateFileRequest) (*ReplicateFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicateFile not implemented")
}
func (UnimplementedDataKeeperServer) mustEmbedUnimplementedDataKeeperServer() {}
func (UnimplementedDataKeeperServer) testEmbeddedByValue()                    {}

// UnsafeDataKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataKeeperServer will
// result in compilation errors.
type UnsafeDataKeeperServer interface {
	mustEmbedUnimplementedDataKeeperServer()
}

func RegisterDataKeeperServer(s grpc.ServiceRegistrar, srv DataKeeperServer) {
	// If the following call pancis, it indicates UnimplementedDataKeeperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataKeeper_ServiceDesc, srv)
}

func _DataKeeper_InitiateUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitiateUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServer).InitiateUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeper_InitiateUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServer).InitiateUpload(ctx, req.(*InitiateUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataKeeper_ReplicateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicateFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServer).ReplicateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeper_ReplicateFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServer).ReplicateFile(ctx, req.(*ReplicateFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataKeeper_ServiceDesc is the grpc.ServiceDesc for DataKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dfs.DataKeeper",
	HandlerType: (*DataKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitiateUpload",
			Handler:    _DataKeeper_InitiateUpload_Handler,
		},
		{
			MethodName: "ReplicateFile",
			Handler:    _DataKeeper_ReplicateFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/dfs.proto",
}
