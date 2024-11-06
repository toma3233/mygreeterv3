targetScope = 'subscription'

@sys.description('Name of the role assignment')
param name string

@sys.description('Required. Provide the role definition id (i.e. c2f4ef07-c644-48eb-af81-4b1b4947fb11).')
param roleDefinitionId string

@sys.description('Required. The Principal or Object ID of the Security Principal (User, Group, Service Principal, Managed Identity).')
param principalId string

@sys.description('Optional. The description of the role assignment.')
param description string = ''

@sys.description('Optional. ID of the delegated managed identity resource.')
param delegatedManagedIdentityResourceId string = ''

@sys.description('Optional. The conditions on the role assignment. This limits the resources it can be assigned to.')
param condition string = ''

@sys.description('Optional. Location deployment metadata.')
param location string = deployment().location

@sys.description('Optional. Version of the condition. Currently accepted value is "2.0".')
@allowed([
  '2.0'
])
param conditionVersion string = '2.0'

@sys.description('Optional. The principal type of the assigned principal ID.')
@allowed([
  'ServicePrincipal'
  'Group'
  'User'
  'ForeignGroup'
  'Device'
  ''
])
param principalType string = ''

@sys.description('Specify true if the resource is a shared resource, or false if the resource is not shared.')
param sharedResource bool = false

@sys.description('The resource ID of the user-assigned managed identity.')
param userAssignedIdentityResourceId string = ''

resource roleAssignmentShared 'Microsoft.Authorization/roleAssignments@2022-04-01' existing = if (sharedResource) {
  name: name
}

resource roleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = if (!(sharedResource)) {
  name: name
  properties: {
    roleDefinitionId: '/providers/Microsoft.Authorization/roleDefinitions/${roleDefinitionId}'
    principalId: principalId
    description: !empty(description) ? description : null
    principalType: !empty(principalType) ? any(principalType) : null
    delegatedManagedIdentityResourceId: !empty(delegatedManagedIdentityResourceId) ? delegatedManagedIdentityResourceId : null
    conditionVersion: !empty(conditionVersion) && !empty(condition) ? conditionVersion : null
    condition: !empty(condition) ? condition : null
  }
}

resource userAssignedIdentityRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(subscription().id, userAssignedIdentityResourceId, roleDefinitionId)
  properties: {
    roleDefinitionId: roleDefinitionId
    principalId: userAssignedIdentityResourceId
  }
}

@sys.description('The GUID of the Role Assignment.')
output name string = name

@sys.description('The resource ID of the Role Assignment.')
output resourceId string = sharedResource ? roleAssignmentShared.id : roleAssignment.id

@sys.description('The scope this Role Assignment applies to.')
output scope string = subscription().id
