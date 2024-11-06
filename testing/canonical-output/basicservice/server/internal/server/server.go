// Auto generated. Don't modify.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/internal/logattrs"

	"github.com/Azure/aks-middleware/interceptor"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	log "log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func Serve(options Options) {
	logger := log.New(log.NewTextHandler(os.Stdout, nil).WithAttrs(logattrs.GetAttrs()))
	if options.JsonLog {
		logger = log.New(log.NewJSONHandler(os.Stdout, nil).WithAttrs(logattrs.GetAttrs()))
	}

	log.SetDefault(logger)

	apiServer := NewServer()
	apiServer.init(options)

	gRpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.DefaultServerInterceptors(interceptor.GetServerInterceptorLogOptions(logger, logattrs.GetAttrs()))...,
	))
	pb.RegisterBasicServiceServer(gRpcServer, apiServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gRpcServer, healthServer)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", options.Port))
	if err != nil {
		logger.Error("failed to listen: " + err.Error())
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("server listening at %s", listener.Addr().String()))
	go func() {
		if err := gRpcServer.Serve(listener); err != nil {
			logger.Error("failed to serve: " + err.Error())
			os.Exit(1)
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("localhost:%d", options.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Error("Failed to dial server: " + err.Error())
		os.Exit(1)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterBasicServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		logger.Error("Failed to register gateway: " + err.Error())
		os.Exit(1)
	}

	gwServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", options.HTTPPort),
		Handler:           gwmux,
		ReadHeaderTimeout: 10 * time.Second, // TODO: determine what time is appropriate
	}

	logger.Info("serving gRPC-Gateway on [::]:" + strconv.Itoa(options.HTTPPort))
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	if err := gwServer.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func GetFreePort() int {
	var listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return -1
	}
	var port = listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	return port
}

func StartServer(serverPort int, httpPort int, demoserverPort int) {
	go func() {
		var serverOptions = Options{}
		serverOptions.Port = serverPort
		serverOptions.JsonLog = false
		serverOptions.HTTPPort = httpPort
		if demoserverPort == -1 {
			serverOptions.RemoteAddr = ""
		} else {
			serverOptions.RemoteAddr = fmt.Sprintf("localhost:%d", demoserverPort)
		}
		Serve(serverOptions)
	}()
}

// TODO: Uncomment the following code once demoserver is merged in
// func StartDemoServer(demoserverPort int) {
// 	go func() {
// 		var demoserverOptions = demoserver.Options{}
// 		demoserverOptions.Port = demoserverPort
// 		demoserverOptions.JsonLog = false
// 		demoserver.Serve(demoserverOptions)
// 	}()
// }

func IsServerRunning(port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}
