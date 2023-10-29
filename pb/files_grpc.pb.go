// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: files.proto

package pb

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

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileServiceClient interface {
	GetChuncks(ctx context.Context, opts ...grpc.CallOption) (FileService_GetChuncksClient, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) GetChuncks(ctx context.Context, opts ...grpc.CallOption) (FileService_GetChuncksClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[0], "/FileService/GetChuncks", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceGetChuncksClient{stream}
	return x, nil
}

type FileService_GetChuncksClient interface {
	Send(*FileRequest) error
	CloseAndRecv() (*FileResponse, error)
	grpc.ClientStream
}

type fileServiceGetChuncksClient struct {
	grpc.ClientStream
}

func (x *fileServiceGetChuncksClient) Send(m *FileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileServiceGetChuncksClient) CloseAndRecv() (*FileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileServiceServer is the server API for FileService service.
// All implementations must embed UnimplementedFileServiceServer
// for forward compatibility
type FileServiceServer interface {
	GetChuncks(FileService_GetChuncksServer) error
	mustEmbedUnimplementedFileServiceServer()
}

// UnimplementedFileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileServiceServer struct {
}

func (UnimplementedFileServiceServer) GetChuncks(FileService_GetChuncksServer) error {
	return status.Errorf(codes.Unimplemented, "method GetChuncks not implemented")
}
func (UnimplementedFileServiceServer) mustEmbedUnimplementedFileServiceServer() {}

// UnsafeFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServiceServer will
// result in compilation errors.
type UnsafeFileServiceServer interface {
	mustEmbedUnimplementedFileServiceServer()
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func _FileService_GetChuncks_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).GetChuncks(&fileServiceGetChuncksServer{stream})
}

type FileService_GetChuncksServer interface {
	SendAndClose(*FileResponse) error
	Recv() (*FileRequest, error)
	grpc.ServerStream
}

type fileServiceGetChuncksServer struct {
	grpc.ServerStream
}

func (x *fileServiceGetChuncksServer) SendAndClose(m *FileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileServiceGetChuncksServer) Recv() (*FileRequest, error) {
	m := new(FileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileService_ServiceDesc is the grpc.ServiceDesc for FileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetChuncks",
			Handler:       _FileService_GetChuncks_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "files.proto",
}
