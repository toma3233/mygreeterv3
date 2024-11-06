package server

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1"
	"<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1/client"
	"<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1/restsdk"
	"<<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/internal/logattrs"
	"github.com/Azure/aks-middleware/interceptor"

	log "log/slog"

	"github.com/Azure/aks-middleware/restlogger"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func sayHello(buf *bytes.Buffer, port int, req *pb.HelloRequest) {
	logger := log.New(log.NewTextHandler(buf, nil))
	host := fmt.Sprintf("localhost:%d", port)
	options := interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs())
	options.APIOutput = buf
	client, err := client.NewClient(host, options)
	// logging the error for transparency, but not failing the test since retry interceptor will handle it
	if err != nil {
		log.Error("did not connect: " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client.SayHello(ctx, req)
}

func AllocatePort() int {
	port := GetFreePort()
	Expect(port).ToNot(Equal(-1)) // must dynamically allocate port
	return port
}

func AllocateDistinctPorts() (int, int, int) {
	serverPort := AllocatePort()
	httpPort := AllocatePort()
	Expect(httpPort).ToNot(Equal(serverPort))
	demoserverPort := AllocatePort()
	Expect(demoserverPort).ToNot(Equal(serverPort))
	Expect(demoserverPort).ToNot(Equal(httpPort))

	return serverPort, httpPort, demoserverPort
}

// TODO: Allow user to filter out serial test cases
// These tests cannot be run in parallel
// - TestRetryWhenUnavailable() may send req to port server is running on
var _ = Describe("Interceptor test", func() {

	var serverPort int
	// var demoserverPort int
	var httpPort int
	var in *pb.HelloRequest

	BeforeEach(func() {
		serverPort, httpPort, _ = AllocateDistinctPorts()
		addr := &pb.Address{
			Street:  "123 Main St",
			City:    "Seattle",
			State:   "WA",
			Zipcode: 98012,
		}
		in = &pb.HelloRequest{Name: "Bob", Age: 53, Email: "test@test.com", Address: addr}
	})

	Context("when initializing the server", func() {
		It("should correctly initialize the server based on the provided options", func() {
			// Create a new server
			s := NewServer()

			// Set the options
			options := Options{
				JsonLog:    true,
				RemoteAddr: "localhost:50151",
			}

			s.init(options)
			Expect(s.client).ToNot(BeNil())
		})
	})

	Context("when the server is available", func() {
		BeforeEach(func() {

			// StartDemoServer(demoserverPort)
			// timeout := time.NewTimer(10 * time.Second)
			// for {
			// 	if IsServerRunning(demoserverPort) {
			// 		break
			// 	}
			// 	time.Sleep(1 * time.Second)
			// 	if !timeout.Stop() {
			// 		<-timeout.C
			// 		log.Error("Server startup check timed out")
			// 		return
			// 	}
			// }
			StartServer(serverPort, httpPort, -1)
			// Explicitly testing state of server
			// Continue with tests once server and demoserver and grpc-gateway are up and running
			Eventually(func() bool {
				return IsServerRunning(serverPort) && IsServerRunning(httpPort) // && IsServerRunning(demoserverPort)
			}, 10*time.Second).Should(BeTrue())
		})

		It("should not retry the request", func() {
			var buf bytes.Buffer
			sayHello(&buf, serverPort, in)
			Expect(strings.Count(buf.String(), "request-id")).To(Equal(1))
			Expect(strings.Count(buf.String(), "OK")).To(Equal(1)) // must not retry
		})

		It("should validate the name length", func() {
			var buf bytes.Buffer
			in.Name = "Z"
			sayHello(&buf, serverPort, in)
			Expect(strings.Count(buf.String(), "request-id")).To(Equal(1))
			Expect(buf.String()).To(ContainSubstring("value length must be at least 2 characters")) // must return error b/c name < 2 letters
		})

		It("should validate the age range", func() {
			var buf bytes.Buffer
			in.Age = 353
			sayHello(&buf, serverPort, in)
			Expect(strings.Count(buf.String(), "request-id")).To(Equal(1))
			Expect(buf.String()).To(ContainSubstring("value must be greater than or equal to 1 and less than 150")) // must return error b/c age > 150
		})

		It("should validate the email format", func() {
			var buf bytes.Buffer
			in.Email = "test"
			sayHello(&buf, serverPort, in)
			Expect(strings.Count(buf.String(), "request-id")).To(Equal(1))
			Expect(buf.String()).To(ContainSubstring("value does not match regex pattern")) // must return error b/c invalid email
		})

		It("should recover from panic", func() {
			os.Setenv("AKS_BIN_VERSION_GITBRANCH", "tomabraham/service")
			var buf bytes.Buffer
			in.Name = "TestPanic"
			sayHello(&buf, serverPort, in) // "TestPanic" is a special name that triggers panic
			Expect(buf.String()).To(ContainSubstring("SayHello.go, line:"))
			Expect(strings.Count(buf.String(), "code=Unknown")).To(Equal(1)) // must handle panic by returning gRPC code unknown and output filename/line # of panic
		})
	})

	Context("when the server is unavailable", func() {
		It("should retry the request", func() {
			var buf bytes.Buffer
			sayHello(&buf, serverPort, in)
			Expect(strings.Count(buf.String(), "request-id")).To(BeNumerically(">", 1)) // must retry
		})
	})
})

var _ = Describe("REST call test", func() {

	var serverPort int
	// var demoserverPort int
	var httpPort int

	BeforeEach(func() {
		serverPort, httpPort, _ = AllocateDistinctPorts()
		// StartDemoServer(demoserverPort)
		// timeout := time.NewTimer(10 * time.Second)
		// for {
		// 	if IsServerRunning(demoserverPort) {
		// 		break
		// 	}
		// 	time.Sleep(1 * time.Second)
		// 	if !timeout.Stop() {
		// 		<-timeout.C
		// 		log.Error("Server startup check timed out")
		// 		return
		// 	}
		// }
		StartServer(serverPort, httpPort, -1)
		Eventually(func() bool {
			return IsServerRunning(serverPort) && IsServerRunning(httpPort) //&& IsServerRunning(demoserverPort)
		}, 10*time.Second).Should(BeTrue())
	})

	Context("when sending a REST request", func() {
		It("should return successfully when making SayHello call", func() {
			logger := log.New(log.NewTextHandler(os.Stdout, nil))
			// Create a new Configuration instance
			cfg := &restsdk.Configuration{
				BasePath:      fmt.Sprintf("http://0.0.0.0:%d", httpPort),
				DefaultHeader: make(map[string]string),
				UserAgent:     "Swagger-Codegen/1.0.0/go",
				HTTPClient:    restlogger.NewLoggingClient(logger),
			}

			apiClient := restsdk.NewAPIClient(cfg)

			service := apiClient.<<.serviceInput.serviceName>>Api

			helloRequestBody := restsdk.HelloRequest{
				Name:  "MyName",
				Age:   53,
				Email: "test@test.com",
			}
			resp, _, err := service.<<.serviceInput.serviceName>>SayHello(context.Background(), helloRequestBody)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Message).To(ContainSubstring("Echo back what you sent me (SayHello): MyName 53 test@test.com"))
		})
	})
})
