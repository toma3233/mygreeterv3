package client

import (
	"github.com/Azure/aks-middleware/interceptor"
	pb "<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1"

	log "log/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient returns a client that has all the interceptors registered.
func NewClient(remoteAddr string, options interceptor.ClientInterceptorLogOptions) (pb.<<.serviceInput.serviceName>>Client, error) {
	conn, err := grpc.Dial(
		remoteAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.DefaultClientInterceptors(options)...,
		),
	)
	if err != nil {
		// logging for transparency, error handled by retry interceptor
		log.Error("did not connect: " + err.Error())
	}

	return pb.New<<.serviceInput.serviceName>>Client(conn), err
	// TODO: Figure out how to close the connection when the program exits.
	// defer conn.Close()
}



