targetScope = 'subscription'

@sys.description('The name for the resources.')
param resourcesName string

@sys.description('The subscription the resources are deployed to.')
param subscriptionId string

@sys.description('The location of the resource group the resources are deployed to.')
param location string

@sys.description('The name of the resource group the resources are deployed to.')
param resourceGroupName string

// This resource is shared and defined in resources/Main.SharedResources.Template.bicep in shared-resources directory; we only reference it here. Do not remove `existing` syntax.
resource rg 'Microsoft.Resources/resourceGroups@2021-04-01' existing = {
  name: resourceGroupName
  scope: subscription(subscriptionId)
}

// This resource is shared and defined in resources/Main.SharedResources.Template.bicep in shared-resources directory; we only reference it here. Do not remove `existing` syntax.
resource serviceBusNamespace 'Microsoft.ServiceBus/namespaces@2022-10-01-preview' existing = {
  name:  '<<.sharedInput.productShortName>>-${resourcesName}-servicebus-namespace'
  scope: resourceGroup(subscriptionId, resourceGroupName)
}

// This resource is shared and defined in resources/Main.SharedResources.Template.bicep in shared-resources directory; we only reference it here. Do not remove `existing` syntax.
// TODO: If we keep this for a long time, change it to be consistent with the `resource` & `existing` syntax.
module aks 'br:servicehubregistry.azurecr.io/bicep/modules/aks-managed-cluster:v5' = {
  name: '<<.sharedInput.productShortName>>-${resourcesName}-clusterDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: '<<.sharedInput.productShortName>>-${resourcesName}-cluster'
    sharedResource: true // this indicates that the resource is shared and should not be modified.
  }
}

var serviceAccountNamespace = '<<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server'
var serviceAccountName = '<<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server'
module managedIdentity 'br/public:avm/res/managed-identity/user-assigned-identity:0.2.1' = {
  name: '<<.sharedInput.productShortName>>-${resourcesName}-<<.serviceInput.directoryName>>-managed-identityDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    // Name needs to be unique in the entire subscription, thus why we add the `${resourcesName}` to avoid conflicts from different developers.
    name: '<<.sharedInput.productShortName>>-${resourcesName}-<<.serviceInput.directoryName>>-managedIdentity'
    location: rg.location
    federatedIdentityCredentials: [
      {
        name: '<<.sharedInput.productShortName>>-${resourcesName}-<<.serviceInput.directoryName>>-fedIdentity'
        issuer: aks.outputs.oidcIssuerUrl
        subject: 'system:serviceaccount:${serviceAccountNamespace}:${serviceAccountName}'
        audiences: ['api://AzureADTokenExchange']
      }
    ]
  }
}

// TODO: migrate to use bicep module registry since it's available
module azureSdkRoleAssignment 'br:servicehubregistry.azurecr.io/bicep/modules/subscription-role-assignment:v3' = {
  name: '<<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>azuresdkra${location}Deploy'
  scope: subscription(subscriptionId)
  params: {
    name: guid(
      '<<.serviceInput.directoryName>>azuresdk',
      'Contributor',
      managedIdentity.outputs.principalId,
      resourcesName,
      location
    )
    location: rg.location
    roleDefinitionId: 'b24988ac-6180-42a0-ab88-20f7382dd24c' // Contributor
    principalId: managedIdentity.outputs.principalId
    principalType: 'ServicePrincipal'
    sharedResource: false
  }
}

module resourceRoleAssignment 'br/public:avm/ptn/authorization/resource-role-assignment:0.1.1' = {
  name: 'resourceRoleAssignmentDeployment'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    // Required parameters
    principalId: managedIdentity.outputs.principalId
    resourceId: serviceBusNamespace.id
    roleDefinitionId: '090c5cfd-751d-490a-894a-3ce6f1109419'
    // Non-required parameters
    description: 'Assign Service Bus Data Sender permissions to managed identity.'
    principalType: 'ServicePrincipal'
    roleName: 'Service Bus Data Sender'
  }
}

@sys.description('Client Id of the managed identity.')
output clientId string = managedIdentity.outputs.clientId
