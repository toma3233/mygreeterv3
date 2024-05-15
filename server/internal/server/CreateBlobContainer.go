package server

import (
	"context"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateBlobContainer creates a new blob storage container within a specified storage account.
func (s *Server) CreateBlobContainer(ctx context.Context, in *pb.CreateBlobContainerRequest) (*pb.CreateBlobContainerResponse, error) {
	logger := ctxlogger.GetLogger(ctx)

	// Retrieve the storage account key
	keys, err := s.AccountsClient.ListKeys(ctx, in.GetRgName(), in.GetSaName(), nil)
	if err != nil {
		logger.Error("Failed to list keys for storage account: ", in.GetSaName(), " Error: ", err)
		return nil, status.Errorf(codes.Internal, "Failed to list keys for storage account: %s", in.GetSaName())
	}
	accountKey := *(*keys.Keys)[0].Value

	// Create a service client
	serviceURL, err := azblob.NewServiceClientWithSharedKey("https://"+in.GetSaName()+".blob.core.windows.net/", accountKey, nil)
	if err != nil {
		logger.Error("Failed to create service client: ", err)
		return nil, status.Errorf(codes.Internal, "Failed to create service client")
	}

	// Create the container
	containerURL := serviceURL.NewContainerClient(in.GetContainerName())
	_, err = containerURL.Create(ctx, nil)
	if err != nil {
		logger.Error("Failed to create container: ", in.GetContainerName(), " Error: ", err)
		return nil, status.Errorf(codes.Internal, "Failed to create container: %s", in.GetContainerName())
	}

	logger.Info("Blob container created successfully: ", in.GetContainerName())
	return &pb.CreateBlobContainerResponse{Name: in.GetContainerName()}, nil
}
