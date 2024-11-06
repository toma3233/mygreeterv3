package server

import (
	"context"
	log "log/slog"
	"os"

	pb "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/api/v1"
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/server/internal/logattrs"
	"github.com/Azure/aks-async/servicebus"
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
	// TODO(mheberling): uncomment once operation_container is implemented.
	// containerClient     pt.OperationContainerClient
	serviceBusClient servicebus.ServiceBusClientInterface
	serviceBusSender servicebus.SenderInterface
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

	if options.ServiceBusHostName != "" {
		s.serviceBusClient, err = servicebus.CreateServiceBusClient(context.Background(), options.ServiceBusHostName, nil, nil)
		if err != nil {
			logger.Error("Something went wrong creating the service bus client: " + err.Error())
			os.Exit(1)
		}
	}

	if options.ServiceBusQueueName != "" {
		s.serviceBusSender, err = s.serviceBusClient.NewServiceBusSender(context.Background(), options.ServiceBusQueueName, nil)
		if err != nil {
			logger.Error("Something went wrong creating the service bus sender: " + err.Error())
			os.Exit(1)
		}
	}

	//TODO(mheberling): Uncomment when operation_container is implemented.
	// if options.OperationContainerAddr != "" {
	// 	s.containerClient, err = containerClient.NewClient(options.OperationContainerAddr, interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs()))
	// 	// containerClient.NewOperationContainerClient(options.OperationContainerAddr, interceptor.GetClientInterceptorLogOptions(logger, logattrs.GetAttrs()))
	// 	if err != nil {
	// 		log.Error("did not connect to containerClient: " + err.Error())
	// 	}
	// }

}
