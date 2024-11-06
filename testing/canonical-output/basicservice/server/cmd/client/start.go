package main

import (
	"context"
	"io"
	"os"
	"time"

	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1/client"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/internal/logattrs"

	"strconv"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1/restsdk"
	"github.com/Azure/aks-middleware/interceptor"
	"github.com/Azure/aks-middleware/restlogger"

	"strings"

	log "log/slog"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "hello",
	Short: "Call SayHello",
	Run:   hello,
}

var output io.Writer = os.Stdout

type Options struct {
	RemoteAddr       string
	HttpAddr         string
	JsonLog          bool
	Name             string
	Age              int32
	Email            string
	Address          string
	IntervalMilliSec int64
}

var options = newOptions()

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&options.RemoteAddr, "remote-addr", options.RemoteAddr, "the remote server's addr for this client to connect to")
	startCmd.Flags().StringVar(&options.HttpAddr, "http-addr", options.HttpAddr, "the remote http gateway addr")
	startCmd.Flags().BoolVar(&options.JsonLog, "json-log", options.JsonLog, "The format of the log is json or user friendly key-value pairs")
	startCmd.Flags().StringVar(&options.Name, "name", options.Name, "The name to send in Hello request")
	startCmd.Flags().Int32Var(&options.Age, "age", options.Age, "The age to send in Hello request")
	startCmd.Flags().StringVar(&options.Email, "email", options.Email, "The email to send in Hello request")
	startCmd.Flags().StringVar(&options.Address, "address", options.Address, "The address to send in Hello request")
	startCmd.Flags().Int64Var(&options.IntervalMilliSec, "interval-milli-sec", options.IntervalMilliSec,
		"The interval between two requests. Negative numbers mean sending one request.")
}

func newOptions() Options {
	return Options{
		RemoteAddr:       "localhost:50151",
		HttpAddr:         "http://localhost:50161",
		JsonLog:          false,
		Name:             "MyName",
		Age:              53,
		Email:            "test@test.com",
		Address:          "123 Main St, Seattle, WA 98101",
		IntervalMilliSec: -1,
	}
}

func SetOutput(out io.Writer) {
	output = out
}

func hello(cmd *cobra.Command, args []string) {
	logger := log.New(log.NewTextHandler(output, nil).WithAttrs(logattrs.GetAttrs()))
	if options.JsonLog {
		logger = log.New(log.NewJSONHandler(output, nil).WithAttrs(logattrs.GetAttrs()))
	}

	log.SetDefault(logger)

	client, err := client.NewClient(options.RemoteAddr, interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs()))
	// logging the error for transparency, retry interceptor will handle it
	if err != nil {
		log.Error("did not connect: " + err.Error())
	}

	if options.IntervalMilliSec < 0 {
		SayHello(client, options.Name, options.Age, options.Email, options.Address, logger)
		return
	}
	for {
		SayHello(client, options.Name, options.Age, options.Email, options.Address, logger)
		time.Sleep(time.Duration(options.IntervalMilliSec) * time.Millisecond)
	}
}

func SayHello(client pb.BasicServiceClient, name string, age int32, email string, address string, logger *log.Logger) {
	ctx := context.Background()

	addressParts := strings.Split(address, ",")
	street := addressParts[0]
	city := addressParts[1]
	stateAndZip := strings.Split(addressParts[2], " ")
	state := stateAndZip[1]
	parseZipcode, _ := strconv.ParseInt(stateAndZip[2], 10, 32)
	zipcode := int32(parseZipcode)

	addr := &pb.Address{
		Street:  street,
		City:    city,
		State:   state,
		Zipcode: zipcode,
	}
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name, Age: age, Email: email, Address: addr})
	if err != nil {
		log.Info("SayHello error: " + err.Error())
	} else {
		log.Info("Response message: " + r.GetMessage())
	}

	// Create a new Configuration instance
	cfg := &restsdk.Configuration{
		BasePath:      options.HttpAddr,
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
		HTTPClient:    restlogger.NewLoggingClient(logger),
	}

	apiClient := restsdk.NewAPIClient(cfg)

	service := apiClient.BasicServiceApi

	restAddr := &restsdk.Address{
		Street:  street,
		City:    city,
		State:   state,
		Zipcode: zipcode,
	}

	sayHelloRequest := restsdk.HelloRequest{
		Name:    name,
		Age:     age,
		Email:   email,
		Address: restAddr,
	}
	_, _, err = service.BasicServiceSayHello(context.TODO(), sayHelloRequest)
	if err != nil {
		log.Error("Error calling SayHello with restsdk" + err.Error())
	}
}
