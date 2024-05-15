package server

import (
	log "log/slog"
	"os"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/server/internal/logattrs"
	"github.com/Azure/aks-middleware/interceptor"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1/client"
	serviceHubPolicy "github.com/Azure/aks-middleware/policy"
)

type Server struct {
	// When the UnimplementedMyGreeterServer struct is embedded,
	// the generated method/implementation in .pb file will be associated with this struct.
	// If this struct doesn't implment some methods,
	// the .pb ones will be used. If this struct implement the methods, it will override the .pb ones.
	// The reason is that anonymous field's methods are promoted to the struct.
	//
	// When this struct is NOT embedded,, all methods have to be implemented to meet the interface requirement.
	// See https://go.dev/ref/spec#Struct_types.
	pb.UnimplementedMyGreeterServer
	ResourceGroupClient *armresources.ResourceGroupsClient
	AccountsClient      *armstorage.AccountsClient
	client              pb.MyGreeterClient
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) init(options Options) {
	var err error
	var cred azcore.TokenCredential

	logger := log.New(log.NewTextHandler(os.Stdout, nil).WithAttrs(logattrs.GetAttrs()))
	if options.JsonLog {
		logger = log.New(log.NewJSONHandler(os.Stdout, nil).WithAttrs(logattrs.GetAttrs()))
	}

	log.SetDefault(logger)
	if options.EnableAzureSDKCalls {
		armClientOptions := serviceHubPolicy.GetDefaultArmClientOptions(logger)
		// Use MSI in Standalone E2E env for credential
		if options.IdentityResourceID != "" {
			resourceID := azidentity.ResourceID(options.IdentityResourceID)
			opts := azidentity.ManagedIdentityCredentialOptions{ID: resourceID}
			cred, err = azidentity.NewManagedIdentityCredential(&opts)
		} else {
			cred, err = azidentity.NewDefaultAzureCredential(nil)
		}
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		resourcesClientFactory, err := armresources.NewClientFactory(options.SubscriptionID, cred, armClientOptions)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		s.ResourceGroupClient = resourcesClientFactory.NewResourceGroupsClient()
		s.AccountsClient, err = armstorage.NewAccountsClient(options.SubscriptionID, cred, armClientOptions)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	if options.RemoteAddr != "" {
		s.client, err = client.NewClient(options.RemoteAddr, interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs()))
		// logging the error for transparency, retry interceptor will handle it
		if err != nil {
			log.Error("did not connect: " + err.Error())
		}
	}
}

// Implementing the interface methods for the new blob container operations
func (s *Server) CreateBlobContainer(ctx context.Context, in *pb.CreateBlobContainerRequest) (*pb.CreateBlobContainerResponse, error) {
	// Implementation logic for creating a blob container
	return &pb.CreateBlobContainerResponse{}, nil
}

func (s *Server) ReadBlobContainer(ctx context.Context, in *pb.ReadBlobContainerRequest) (*pb.ReadBlobContainerResponse, error) {
	// Implementation logic for reading a blob container
	return &pb.ReadBlobContainerResponse{}, nil
}

func (s *Server) UpdateBlobContainer(ctx context.Context, in *pb.UpdateBlobContainerRequest) (*pb.UpdateBlobContainerResponse, error) {
	// Implementation logic for updating a blob container
	return &pb.UpdateBlobContainerResponse{}, nil
}

func (s *Server) DeleteBlobContainer(ctx context.Context, in *pb.DeleteBlobContainerRequest) (*emptypb.Empty, error) {
	// Implementation logic for deleting a blob container
	return &emptypb.Empty{}, nil
}

func (s *Server) ListBlobContainers(ctx context.Context, in *pb.ListBlobContainersRequest) (*pb.ListBlobContainersResponse, error) {
	// Implementation logic for listing blob containers
	return &pb.ListBlobContainersResponse{}, nil
}
