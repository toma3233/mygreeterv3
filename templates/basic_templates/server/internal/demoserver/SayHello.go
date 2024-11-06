package demoserver

import (
	"context"
	"strconv"
	"time"

	pb "<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
)

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.GetName() == "TestPanic" {
		panic("testing panic")
	}
	logger := ctxlogger.GetLogger(ctx)
	logger.Info("API handler logger output. req: " + in.String())

	time.Sleep(400 * time.Millisecond)
	return &pb.HelloReply{Message: "Echo back what you sent me (SayHello): " + in.GetName() + " " + strconv.Itoa(int(in.GetAge())) + " " + in.GetEmail()}, nil
}
