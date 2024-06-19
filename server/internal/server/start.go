package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
)

// Server represents the gRPC server
type Server struct {
	pb.UnimplementedMyGreeterServer
}

// NewServer creates a new gRPC server instance
func NewServer() *Server {
	return &Server{}
}

// Start runs the gRPC server
func (s *Server) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMyGreeterServer(grpcServer, s)
	reflection.Register(grpcServer)

	return grpcServer.Serve(lis)
}
