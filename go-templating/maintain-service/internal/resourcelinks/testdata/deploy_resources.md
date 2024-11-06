# Resource Information

## Resource and Resource Type

| Name         | ResourceType |
|:-------------|:-------------|
| [testacr](https://portal.azure.com/#@microsoft.onmicrosoft.com/resource/subscriptions/1111/resourceGroups/test-rg/providers/Microsoft.ContainerRegistry/registries/testacr) |  Microsoft.ContainerRegistry/registries |
| [test-cluster](https://portal.azure.com/#@microsoft.onmicrosoft.com/resource/subscriptions/1111/resourceGroups/test-rg/providers/Microsoft.ContainerService/managedClusters/test-cluster) |  Microsoft.ContainerService/managedClusters |


## Deployments and Dependencies

| Name                  | Depends On  |
|:----------------------|:------------|
| test-clusterDeploy | <ul><li>test-data-collection-ruleDeploy</li><li>test-workspaceDeploy</li></ul>  |
| test-acrDeploy | <ul><li>test-rgDeploy</li><li>test-workspaceDeploy</li></ul>  |
