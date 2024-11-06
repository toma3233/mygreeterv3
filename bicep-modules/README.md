## Bicep Modules

IMPT: We are migrating to use bicep registry modules (AVM) rather than our internal modules. We only offer modules that are not offered or do not offer the log that we need. 

Note: To reference a module, please use the following syntax `br:servicehubregistry.azurecr.io/bicep/modules/{module-name}:v1`.

Here are the modules we offer. You can use the following names for {module-name}:
- aks-managed-cluster
- subscription-role-assignment

Example:
```bash
param location string = resourceGroup().location
param resourceGroupName string = 'example-rg'

module aks 'br:servicehubregistry.azurecr.io/bicep/modules/aks-managed-cluster:v1' = {
  name: 'servicehub-example-clusterDeploy'
  scope: resourceGroup(resourceGroupName)
  params: {
    name: 'servicehub-example-cluster'
    location: location
    dnsPrefix: 'example'
    sharedResource: false
  }
}
```