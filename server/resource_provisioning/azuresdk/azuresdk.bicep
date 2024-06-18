targetScope = 'subscription'

@sys.description('The name for the resources.')
param resourcesName string

@sys.description('The subscription the resources are deployed to.')
param subscriptionId string

var resourceGroupName = 'servicehub-${resourcesName}-rg'

// This resource is shared and defined in main.bicep in shared-resources directory; we only reference it here. Do not remove `existing` syntax.
resource rg 'Microsoft.Resources/resourceGroups@2021-04-01' existing = {
  name: resourceGroupName
  scope: subscription(subscriptionId)
}

// This resource is shared and defined in main.bicep in shared-resources directory; we only reference it here. Do not remove `existing` syntax.
// TODO: If we keep this for a long time, change it to be consistent with the `resource` & `existing` syntax.
module aks 'br:servicehubregistry.azurecr.io/bicep/modules/aks-managed-cluster:v5' = {
  name: 'servicehub-${resourcesName}-clusterDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-${resourcesName}-cluster'
    sharedResource: true // this indicates that the resource is shared and should not be modified.
  }
}

var serviceAccountNamespace = 'servicehub-mygreeterv3-server'
var serviceAccountName = 'servicehub-mygreeterv3-server'
module managedIdentity 'br/public:avm/res/managed-identity/user-assigned-identity:0.2.1' = {
  name: 'servicehub-mygreeterv3-managed-identityDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-mygreeterv3-managedIdentity'
    location: rg.location
    federatedIdentityCredentials: [
      {
        name: 'servicehub-mygreeterv3-fedIdentity'
        issuer: aks.outputs.oidcIssuerUrl
        subject: 'system:serviceaccount:${serviceAccountNamespace}:${serviceAccountName}'
        audiences: [ 'api://azureadtokenexchange' ]
      }
    ]
  }
}

module ownerManagedIdentity 'br/public:avm/res/managed-identity/user-assigned-identity:0.2.1' = {
  name: 'servicehub-mygreeterv3-owner-managed-identity'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-mygreeterv3-owner-managed-identity'
    location: rg.location
  }
}

// TODO: migrate to use bicep module registry once it's available
module azureSdkRoleAssignment 'br:servicehubregistry.azurecr.io/bicep/modules/subscription-role-assignment:v3' = {
  name: 'servicehub-mygreeterv3-azuresdk-role-assignmentDeploy'
  scope: subscription(subscriptionId)
  params: {
    name: guid('servicehub-mygreeterv3-azureSdkRoleAssignment', subscriptionId, 'Contributor', managedIdentity.outputs.principalId)
    location: rg.location
    roleDefinitionId: 'b24988ac-6180-42a0-ab88-20f7382dd24c' // Contributor
    principalId: managedIdentity.outputs.principalId
    principalType: 'ServicePrincipal'
    sharedResource: false
  }
}

module ownerRoleAssignment 'br:servicehubregistry.azurecr.io/bicep/modules/subscription-role-assignment:v3' = {
  name: 'servicehub-mygreeterv3-owner-role-assignmentDeploy'
  scope: subscription(subscriptionId)
  params: {
    name: guid('servicehub-mygreeterv3-ownerRoleAssignment', subscriptionId, 'Owner', ownerManagedIdentity.outputs.principalId)
    location: rg.location
    roleDefinitionId: '8e3af657-a8ff-443c-a75c-2fe8c4bcb635' // Owner
    principalId: ownerManagedIdentity.outputs.principalId
    principalType: 'ServicePrincipal'
    sharedResource: false
  }
}

@sys.description('Client Id of the managed identity.')
output clientId string = managedIdentity.outputs.clientId
