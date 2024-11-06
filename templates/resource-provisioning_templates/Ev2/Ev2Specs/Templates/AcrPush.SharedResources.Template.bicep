// These resources are Ev2 specific. We use the identity to push images to the ACR in our Ev2 pipeline.
@sys.description('The name for the resources.')
param resourcesName string

@sys.description('The location of the resource group the resources are deployed to.')
param location string

resource registry 'Microsoft.ContainerRegistry/registries@2023-06-01-preview' existing = {
  name: '<<.sharedInput.productShortName>>${resourcesName}${location}acr'
}

var identityName = '<<.sharedInput.productShortName>>${resourcesName}pipelineidentity01'
resource Identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: identityName
  location: location
}

// The identity is used by all of Ev2 across services to push images to the ACR.
var msiResourceId = '/subscriptions/${subscription().subscriptionId}/resourceGroups/${resourceGroup().name}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/${identityName}'
resource acr_identityName_push 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  scope: registry
  name: guid('<<.sharedInput.productShortName>>${resourcesName}${location}acr', Identity.name, 'push')
  properties: {
    roleDefinitionId: subscriptionResourceId(
      'Microsoft.Authorization/roleDefinitions',
      '8311e382-0749-4cb8-b61a-304f252e45ec'
    )
    principalId: reference(msiResourceId, '2023-01-31').principalId
    principalType: 'ServicePrincipal'
  }
}
