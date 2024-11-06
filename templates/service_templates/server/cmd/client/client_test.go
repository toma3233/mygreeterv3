package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"<<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/internal/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Cobra Cmd test", func() {
	var serverPort int
	var httpPort int
	var cmd *cobra.Command

	BeforeEach(func() {
		serverPort = server.GetFreePort()
		httpPort = server.GetFreePort()
		server.StartServer(serverPort, httpPort, -1)
		Eventually(func() bool {
			return server.IsServerRunning(serverPort) && server.IsServerRunning(httpPort)
		}, 10*time.Second).Should(BeTrue())

		cmd = &cobra.Command{
			Use:   "hello",
			Short: "Call SayHello",
			Run:   hello,
		}
	})

	AfterEach(func() {
		SetOutput(os.Stdout)
	})

	It("should call Execute() and log the response message", func() {
		var buf bytes.Buffer
		SetOutput(&buf)

		host := fmt.Sprintf("localhost:%d", serverPort)
		options.RemoteAddr = host
		options.JsonLog = true

		hello(cmd, nil)
		Expect(buf.String()).To(ContainSubstring("Echo back what you sent me (SayHello)"))
	})

	It("should call Execute() and log error", func() {
		var buf bytes.Buffer
		SetOutput(&buf)

		options.RemoteAddr = "fakeaddr"
		options.JsonLog = true

		hello(cmd, nil)
		Expect(buf.String()).To(ContainSubstring("connect: connection refused"))
	})
})