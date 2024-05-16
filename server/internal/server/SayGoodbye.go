package server

import (
	"context"
	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
)

// SayGoodBye implements the MyGreeterServer interface for the Server.
func (s *Server) SayGoodBye(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	logger := ctxlogger.GetLogger(ctx)
	logger.Info("SayGoodBye request received")

	// If the server is configured to forward requests to the demoserver, do so
	if s.client != nil {
		return s.client.SayGoodBye(ctx, req)
	}

	// Otherwise, construct the goodbye message using the name from the request
	message := "Goodbye, " + req.GetName() + "!"
	return &pb.HelloReply{Message: message}, nil
}
