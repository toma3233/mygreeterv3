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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1/client"
	serviceHubPolicy "github.com/Azure/aks-middleware/policy"
)

type Server struct {
	pb.UnimplementedMyGreeterServer
	ResourceGroupClient *armresources.ResourceGroupsClient
	AccountsClient      *armstorage.AccountsClient
	BlobContainersClient *azblob.ContainerClient
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
		// Initialize BlobContainersClient
		blobServiceClient, err := azblob.NewServiceClientWithNoCredential(options.StorageAccountURL, nil)
		if err != nil {
			log.Error("Failed to create BlobServiceClient: ", err)
			return
		}
		s.BlobContainersClient = blobServiceClient.NewContainerClient(options.ContainerName)
	}

	if options.RemoteAddr != "" {
		s.client, err = client.NewClient(options.RemoteAddr, interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs()))
		// logging the error for transparency, retry interceptor will handle it
		if err != nil {
			log.Error("did not connect: " + err.Error())
		}
	}
}
