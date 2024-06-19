package server

import (
	"context"
	"strconv"
	"time"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
)

func (s *Server) SayGoodbye(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logger := ctxlogger.GetLogger(ctx)
	logger.Info("API handler logger output. req: " + in.String())

	var err error
	var out = &pb.HelloReply{}

	time.Sleep(200 * time.Millisecond)
	if s.client != nil {
		out, err = s.client.SayGoodbye(ctx, in)
		if err != nil {
			return out, err
		}
		out.Message += "| appended by server"
	} else {
		if in.GetName() == "TestPanic" {
			panic("testing panic")
		}
		logger := ctxlogger.GetLogger(ctx)
		logger.Info("API handler logger output. req: " + in.String())

		time.Sleep(400 * time.Millisecond)
		out, err = &pb.HelloReply{Message: "Echo back what you sent me (SayGoodbye): " + in.GetName() + " " + strconv.Itoa(int(in.GetAge())) + " " + in.GetEmail()}, nil
	}

	return out, err
}
