package server

import (
	"context"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
)

// SayGoodbye implements the SayGoodbye method of the MyGreeter service.
func (s *Server) SayGoodbye(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logger := ctxlogger.GetLogger(ctx)
	logger.Info("SayGoodbye called with request: ", in.String())

	var message string
	if s.client != nil {
		out, err := s.client.SayGoodbye(ctx, in)
		if err != nil {
			return nil, err
		}
		message = out.GetMessage() + "| appended by server"
	} else {
		message = "Goodbye, " + in.GetName() + "!"
	}

	return &pb.HelloReply{Message: message}, nil
}
