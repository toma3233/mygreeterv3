package server

import (
	"context"
	"strconv"
	"time"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
)

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logger := ctxlogger.GetLogger(ctx)
	logger.Info("API handler logger output. req: " + in.String())

	if in.GetName() == "TestPanic" {
		panic("testing panic")
	}

	time.Sleep(200 * time.Millisecond)

	var err error
	var out = &pb.HelloReply{}
	if s.client != nil {
		out, err = s.client.SayHello(ctx, in)
		if err != nil {
			return out, err
		}
		out.Message += "| appended by server"
	} else {
		out, err = &pb.HelloReply{Message: "Echo back what you sent me (SayHello): " + in.GetName() + " " + strconv.Itoa(int(in.GetAge())) + " " + in.GetEmail()}, nil
	}
	return out, err
}
