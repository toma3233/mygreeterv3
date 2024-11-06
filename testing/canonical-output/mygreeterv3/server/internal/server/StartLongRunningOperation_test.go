package server

import (
	"context"
	"encoding/json"

	"time"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	opbus "github.com/Azure/aks-async/operationsbus"
	"github.com/Azure/aks-async/servicebus"

	asyncMocks "github.com/Azure/aks-async/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("Mock Testing", func() {
	var (
		ctrl       *gomock.Controller
		s          *Server
		mockSender *asyncMocks.MockSenderInterface
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSender = asyncMocks.NewMockSenderInterface(ctrl)
		s = &Server{
			serviceBusSender: mockSender,
		}
	})

	Context("async operations", func() {
		It("should return operationId", func() {
			protoExpirationTime := timestamppb.New(time.Now().Add(1 * time.Hour))
			in := &pb.StartLongRunningOperationRequest{
				EntityId:            "20",
				EntityType:          "Cluster",
				ExpirationTimestamp: protoExpirationTime,
			}

			mockSender.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			out, err := s.StartLongRunningOperation(context.Background(), in)
			Expect(err).To(BeNil())
			Expect(out.OperationId).NotTo(BeNil())
		})
	})
})

var _ = Describe("Fakes Testing", func() {
	var (
		s *Server
	)

	BeforeEach(func() {
		sbClient := servicebus.NewFakeServiceBusClient()
		sbSender, _ := sbClient.NewServiceBusSender(nil, "", nil)
		s = &Server{ResourceGroupClient: nil, serviceBusClient: sbClient, serviceBusSender: sbSender}
	})

	Context("Message should exist in the service bus", func() {
		It("should send the message successfully", func() {
			protoExpirationTime := timestamppb.New(time.Now().Add(1 * time.Hour))
			in := &pb.StartLongRunningOperationRequest{
				EntityId:            "20",
				EntityType:          "Cluster",
				ExpirationTimestamp: protoExpirationTime,
			}
			_, err := s.StartLongRunningOperation(context.Background(), in)
			Expect(err).ToNot(HaveOccurred())

			sbReceiver, _ := s.serviceBusClient.NewServiceBusReceiver(nil, "", nil)
			msg, err := sbReceiver.ReceiveMessage(nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(msg).NotTo(BeNil())

			opRequestExpected := opbus.NewOperationRequest("LongRunningOperation", "", "", "20", "Cluster", 0, protoExpirationTime, nil, "", nil)

			var opRequestReceived opbus.OperationRequest
			err = json.Unmarshal(msg, &opRequestReceived)
			Expect(err).ToNot(HaveOccurred())

			Expect(opRequestReceived.OperationName).To(Equal(opRequestExpected.OperationName))
			Expect(opRequestReceived.OperationId).NotTo(BeNil())
			Expect(opRequestReceived.RetryCount).To(Equal(opRequestExpected.RetryCount))
			Expect(opRequestReceived.EntityType).To(Equal(opRequestExpected.EntityType))
			Expect(opRequestReceived.EntityId).To(Equal(opRequestExpected.EntityId))
			Expect(opRequestReceived.ExpirationTimestamp).To(Equal(opRequestExpected.ExpirationTimestamp))
			Expect(opRequestReceived.APIVersion).To(Equal(opRequestExpected.APIVersion))
		})
	})
})
