package server

import (
	"context"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateBlobContainer updates the metadata of a specified blob storage container.
func (s *Server) UpdateBlobContainer(ctx context.Context, in *pb.UpdateBlobContainerRequest) (*pb.UpdateBlobContainerResponse, error) {
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

	// Get the container client
	containerURL := serviceURL.NewContainerClient(in.GetContainerName())
	_, err = containerURL.SetMetadata(ctx, in.GetMetadata(), nil)
	if err != nil {
		logger.Error("Failed to update metadata for container: ", in.GetContainerName(), " Error: ", err)
		return nil, status.Errorf(codes.Internal, "Failed to update metadata for container: %s", in.GetContainerName())
	}

	logger.Info("Blob container metadata updated successfully: ", in.GetContainerName())
	return &pb.UpdateBlobContainerResponse{
		BlobContainer: &pb.BlobContainer{
			Id:       containerURL.URL(),
			Name:     in.GetContainerName(),
			Metadata: in.GetMetadata(),
		},
	}, nil
}
