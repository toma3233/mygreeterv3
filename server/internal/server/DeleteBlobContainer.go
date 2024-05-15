package server

import (
	"context"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteBlobContainer deletes a specified blob storage container.
func (s *Server) DeleteBlobContainer(ctx context.Context, in *pb.DeleteBlobContainerRequest) (*emptypb.Empty, error) {
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
	_, err = containerURL.Delete(ctx, nil)
	if err != nil {
		logger.Error("Failed to delete container: ", in.GetContainerName(), " Error: ", err)
		return nil, status.Errorf(codes.Internal, "Failed to delete container: %s", in.GetContainerName())
	}

	logger.Info("Blob container deleted successfully: ", in.GetContainerName())
	return &emptypb.Empty{}, nil
}
