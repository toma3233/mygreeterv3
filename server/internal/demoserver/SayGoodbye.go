package demoserver

import (
	"context"
	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
)

// SayGoodBye implements the MyGreeterServer interface for the DemoServer.
func (s *Server) SayGoodBye(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// Construct the goodbye message using the name from the request.
	message := "Goodbye, " + req.GetName() + "!"
	return &pb.HelloReply{Message: message}, nil
}
