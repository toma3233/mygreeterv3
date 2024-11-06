package server

import (
	"context"
	"encoding/json"

	"github.com/Azure/aks-async/operationsbus"
	"github.com/Azure/aks-middleware/ctxlogger"
	"github.com/gofrs/uuid"

	pb "<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1"
	//TODO(mheberling): Connect to the operationContainer when it exists
	// pt "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/operationContainer/api/v1"
)

func (s *Server) StartLongRunningOperation(ctx context.Context, in *pb.StartLongRunningOperationRequest) (*pb.StartLongRunningOperationResponse, error) {
	logger := ctxlogger.GetLogger(ctx)
	logger.Info("Starting async operation.")

	operationId, err := uuid.NewV4()
	if err != nil {
		logger.Error("Failed to generate UUID: " + err.Error())
		return nil, err
	}

	operation := &operationsbus.OperationRequest{
		OperationName:       LroName,
		APIVersion:          "",
		OperationId:         operationId.String(),
		Body:                nil,
		HttpMethod:          "",
		RetryCount:          0,
		EntityId:            in.GetEntityId(),
		EntityType:          in.GetEntityType(),
		ExpirationTimestamp: in.GetExpirationTimestamp(),
	}

	marshalledOperation, err := json.Marshal(operation)
	if err != nil {
		logger.Error("Error marshalling operation: " + err.Error())
		return nil, err
	}

	logger.Info("Sending message to Service Bus")
	err = s.serviceBusSender.SendMessage(ctx, marshalledOperation)
	if err != nil {
		logger.Error("Error sending message to service bus: " + err.Error())
		return nil, err
	}

	//TODO(mheberling): Uncomment once operationContainer is implemented.
	// createOperationStatusRequest := &pt.CreateOperationStatusRequest{
	// 	OperationName:  in.GetOperationName(),
	// 	EntityId:       in.GetEntityId(),
	// 	EntityType:     in.GetEntityType(),
	// 	ExpirationDate: in.GetExpirationDate(),
	// 	OperationId:    operationId.String(),
	// }
	//
	// // Add the operation to the db.
	// logger.Info("Adding operation to db.")
	// createOperationStatusResponse, err := s.containerClient.CreateOperationStatus(ctx, createOperationStatusRequest)
	// if err != nil {
	// 	logger.Error("Error creating operation status: " + err.Error())
	// }
	//
	// startOperationResponse := &pb.StartOperationResponse{OperationId: createOperationStatusResponse.GetOperationId()}
	startOperationResponse := &pb.StartLongRunningOperationResponse{OperationId: operationId.String()}
	return startOperationResponse, nil
}
