package server

import (
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1/client"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/internal/logattrs"
	"github.com/Azure/aks-middleware/interceptor"
)

type Server struct {
	// When the UnimplementedBasicServiceServer struct is embedded,
	// the generated method/implementation in .pb file will be associated with this struct.
	// If this struct doesn't implment some methods,
	// the .pb ones will be used. If this struct implement the methods, it will override the .pb ones.
	// The reason is that anonymous field's methods are promoted to the struct.
	//
	// When this struct is NOT embedded,, all methods have to be implemented to meet the interface requirement.
	// See https://go.dev/ref/spec#Struct_types.
	pb.UnimplementedBasicServiceServer
	client pb.BasicServiceClient
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) init(options Options) {
	var err error

	logger := log.New(log.NewTextHandler(os.Stdout, nil).WithAttrs(logattrs.GetAttrs()))
	if options.JsonLog {
		logger = log.New(log.NewJSONHandler(os.Stdout, nil).WithAttrs(logattrs.GetAttrs()))
	}

	log.SetDefault(logger)

	if options.RemoteAddr != "" {
		s.client, err = client.NewClient(options.RemoteAddr, interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs()))
		// logging the error for transparency, retry interceptor will handle it
		if err != nil {
			log.Error("did not connect: " + err.Error())
		}
	}

	s.setupShutdown()
}

func (s *Server) setupShutdown() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine to handle shutdown
	go func() {
		<-stopChan
		log.Info("Shutting down the server")

		// Any future connections can be added here.

		log.Info("Server stopped.")
		os.Exit(0)
	}()
}
