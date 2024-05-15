package server

import (
	"context"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"github.com/Azure/aks-middleware/ctxlogger"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListBlobContainers lists all blob storage containers within a specified storage account.
func (s *Server) ListBlobContainers(ctx context.Context, in *pb.ListBlobContainersRequest) (*pb.ListBlobContainersResponse, error) {
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

	// List the containers
	pager := serviceURL.ListContainers(nil)
	var containersList []*pb.BlobContainer
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			logger.Error("Failed to list containers: ", err)
			return nil, status.Errorf(codes.Internal, "Failed to list containers")
		}
		for _, containerItem := range resp.ContainerItems {
			container := &pb.BlobContainer{
				Id:   containerItem.Name,
				Name: containerItem.Name,
				// Metadata is not directly available in the listing operation
				Metadata: map[string]string{},
			}
			containersList = append(containersList, container)
		}
	}

	logger.Info("Blob containers listed successfully")
	return &pb.ListBlobContainersResponse{ContainersList: containersList}, nil
}
