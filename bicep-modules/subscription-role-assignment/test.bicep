
targetScope = 'subscription'

var resourceGroupName = 'servicehubtest-rg'

module rgTest1 '../resource-group/servicehub-generated-module.bicep' = {
  name: '${resourceGroupName}Deploy'
  params: {
    name: resourceGroupName
    location: 'eastus'
    sharedResource: false
  }
}

module aksTest1 '../aks-managed-cluster/servicehub-generated-module.bicep' = {
  name: 'servicehubtest-clusterDeploy'
  scope: resourceGroup(resourceGroupName)
  params: {
    name: 'servicehubtest-cluster'
    location: rgTest1.outputs.location
    dnsPrefix: 'servicehubtest'
    enableWorkloadIdentity: true
    sharedResource: false
  }
}

var serviceAccountNamespace = 'servicehubtest-server'
var serviceAccountName = 'servicehubtest-server'
module managedIdentityTest1 '../managed-identity/servicehub-generated-module.bicep' = {
  name: 'servicehubtest-managedIdentityDeploy'
  scope: resourceGroup(resourceGroupName)
  params: {
    name: 'servicehubtest-ManagedIdentity'
    location: rgTest1.outputs.location
    federatedCredentials: [
      {
        name: 'servicehubtest-FedIdentity'
        issuer: aksTest1.outputs.oidcIssuerUrl
        subject: 'system:serviceaccount:${serviceAccountNamespace}:${serviceAccountName}'
        audiences: ['api://AzureADTokenExchange']
      }
    ]
  }
}

module roleAssignmentTest1 './servicehub-generated-module.bicep' = {
  name: 'servicehubtest-SubscriptionRoleAssignmentDeploy'
  params: {
    name: guid('servicehubtest-subscriptionRoleAssignment', subscription().id, 'Contributor', managedIdentityTest1.outputs.principalId)
    location: rgTest1.outputs.location
    roleDefinitionId: 'b24988ac-6180-42a0-ab88-20f7382dd24c' // Contributor
    principalId: managedIdentityTest1.outputs.principalId
    principalType: 'ServicePrincipal'
  }
}
