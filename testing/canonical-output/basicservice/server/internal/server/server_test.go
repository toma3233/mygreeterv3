package server

import (
	"context"
	"errors"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

var _ = Describe("Mock Testing", func() {
	var (
		ctrl       *gomock.Controller
		mockClient *mock.MockBasicServiceClient
		s          *Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockClient = mock.NewMockBasicServiceClient(ctrl)
		s = &Server{client: mockClient}
	})

	Context("when the server is available", func() {
		It("should not retry the request", func() {
			addr := &pb.Address{
				Street:  "123 Main St",
				City:    "Seattle",
				State:   "WA",
				Zipcode: 98012,
			}
			in := &pb.HelloRequest{Name: "Bob", Age: 53, Email: "test@test.com", Address: addr}
			mockClient.EXPECT().SayHello(gomock.Any(), gomock.Any()).Return(&pb.HelloReply{Message: "Hello, this is a successful mock response"}, nil).Times(1)
			out, err := s.SayHello(context.Background(), in)
			Expect(err).ToNot(HaveOccurred())
			Expect(out.Message).To(Equal("Hello, this is a successful mock response| appended by server"))
		})
	})

	Context("when the server is unavailable", func() {
		It("should retry the request", func() {
			addr := &pb.Address{
				Street:  "123 Main St",
				City:    "Seattle",
				State:   "WA",
				Zipcode: 98012,
			}
			in := &pb.HelloRequest{Name: "Bob", Age: 53, Email: "test@test.com", Address: addr}

			mockClient.EXPECT().SayHello(gomock.Any(), gomock.Any()).Return(nil, errors.New("server unavailable")).Times(2)
			_, err := s.SayHello(context.Background(), in)
			Expect(err).To(HaveOccurred())
			_, err = s.SayHello(context.Background(), in)
			Expect(err).To(HaveOccurred())
		})
	})
})
