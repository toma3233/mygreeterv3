# basicservice





## Setup and Development

Note that we use the remote aks middleware. This middleware is responsible for features such as logging, retry, and input validation. To learn more, please visit the [repo](https://github.com/Azure/aks-middleware/tree/main).

### Initialize service

```bash
./init.sh
# Follow instructions from the scripts to create the api module, etc.
```

### Run Service Locally

There is a simple way to run the BasicService service, after everything has been properly generated. Inside the BasicService directory, you can run the client, demoserver and server.

#### Server

To run the server:

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/server start 
```

By default the server starts in port `localhost:50151`

By default, the sayHello calls are served directly by the server. In order to forward the call to the demoserver:

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/server start --remote-addr <remote_addr>
```

#### Client

To run the client:

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/client hello
```

By default the client sends messages to port `localhost:50151`. This can be changed by running

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/client hello --remote-addr <remote_addr>
```

#### Demoserver

To run the demoserver, you must use a different port than the server is already using, so you can send messages to the demoserver from the server.

To run the demoserver:

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/demoserver start
```

To run the demoserver in a particular port:

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/demoserver start --port <local_port>
```

#### Help

You can run help on every command in order to get more information on how to use them.

Examples:

```bash
go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/client help

go run dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/server/cmd/demoserver start -h
```

### Resource Provisioning

Before creating any service resource, please check [Create or Update Shared Resources](../shared-resources/README.md).

#### Create or Update Service Resources

We have bicep templates set up for you. To deploy the resources:

```bash
cd server
# (optional) modify bicep files in resource provisioning directory
make deploy-resources
```

[Optional] Should you want to modify the parameter values for the bicep resources, follow the instructions in the README section [Making changes to Bicep Resources](../README.md) at the root of the repo.

##### View All Resources and Dependencies

See [resources.md](server/resources.md). This provides a high-level overview of all your deployments.

This file will only exist after you have run `make deploy-resources`. To see the resources you have created and their dependencies, click the different links in this file. Each link is a different markdown file that is associated with a bicep deployment. Each bicep deployment associated file has:

- list of resources you have created via bicep file
- links to the resources in Azure portal
- the dependencies of each resource

| Resource Created | Feature | Shared vs Service |
|----------|----------|----------|
| Resource Group | N/A | Shared |
| AKS Cluster | N/A | Shared |
| Azure Container Registry | N/A | Shared |
| Resource Role Assignment | N/A | Shared |
| Log Analytics Workspace | Monitoring | Shared |
| Data Collection Rule | Monitoring | Shared |
| Alert Rules | Monitoring | Service |
| Managed Identity | AzureSDK | Service |
| Subscription Role Assignment | AzureSDK | Service |

### Run Service on AKS Cluster

Before deploying service to cluster:

- Ensure your [service is running locally](#run-service-locally).
- Complete all steps in [Resource Provisioning](#resource-provisioning).

#### Deploy Service to Cluster
[Dockerfile used to build service image](server/Dockerfile)

```bash
# Assumption: in current service directory
cd server
#Templates env-config.yaml values into all the required files. We assume env-config.yaml exists in your generated folder. (i.e. the folder that stores the generated directories)
make template-files
# Tidys up module dependencies, runs tests, and builds executables
make all
# Build image
# Make sure your api module is tagged to the right version as the go.work file is not used in server/Dockerfile (linked above)
make build-image
# Push image to acr
make push-image
# (If svc running on aks cluster) Upgrade service on AKS cluster
make upgrade
# (If svc not running on aks cluster) Deploy service to AKS cluster
make install
```

#### Check if Service Deployment Successful

You may need wait a few minutes before pods are created and logs show up.

If you do not have kubectl installed you can run these commands to set up the docker container with an environment that will allow you to run the kubectl commands.
```bash
#Assuming you are at the root of the generated directory (the one that contains basicservice)
export src=$(pwd)
docker run -it --mount src=$src,target=/app/binded-data,type=bind servicehubregistry.azurecr.io/service_hub_environment:$20240912 /bin/bash
#Once you are in the container
export KUBECONFIG=app/binded-data/basicservice/server/.kube/config
```

Once inside the container or on your local machine that has kubectl installed

Server:
```bash
# check if pod is running
kubectl get pods -n servicehub-basicservice-server

# check logs
export SERVER_POD=$(kubectl get pod -n servicehub-basicservice-server -o jsonpath="{.items[0].metadata.name}")
kubectl logs $SERVER_POD -n servicehub-basicservice-server
```

Demoserver:
```bash
# check if pod is running
kubectl get pods -n servicehub-basicservice-demoserver

# check logs
export DEMOSERVER_POD=$(kubectl get pod -n servicehub-basicservice-demoserver -o jsonpath="{.items[0].metadata.name}")
kubectl logs $DEMOSERVER_POD -n servicehub-basicservice-demoserver
```

Client
```bash
# check if pod is running
kubectl get pods -n servicehub-basicservice-client

# check logs
export CLIENT_POD=$(kubectl get pod -n servicehub-basicservice-client -o jsonpath="{.items[0].metadata.name}")
kubectl logs $CLIENT_POD -n servicehub-basicservice-client
```


## Debugging and Common Failures


## Monitoring

To view your service in Azure Data Explorer (ADX), follow [ADX dashboard creation/update instructions](server/monitoring/README.md).
