targetScope = 'subscription'

@sys.description('The name for the resources.')
param resourcesName string

@sys.description('The subscription the resources are deployed to.')
param subscriptionId string

@sys.description('The location of the resource group the resources are deployed to.')
param location string

@sys.description('The name of the resource group the resources are deployed to.')
param resourceGroupName string

module rg 'br/public:avm/res/resources/resource-group:0.2.3' = {
  name: '${resourceGroupName}Deploy'
  scope: subscription(subscriptionId)
  params: {
    name: resourceGroupName
    location: location
  }
}

module aks 'br:servicehubregistry.azurecr.io/bicep/modules/aks-managed-cluster:v5' = {
  name: 'servicehub-${resourcesName}-clusterDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-${resourcesName}-cluster'
    location: rg.outputs.location
    dnsPrefix: resourcesName
    monitoringWorkspaceId: workspace.outputs.resourceId
    omsAgentEnabled: true
    createDataCollectionRuleAssociation: true
    dataCollectionRuleId: dataCollectionRule.outputs.resourceId
    enableWorkloadIdentity: true
    enableIstioServiceMesh: true
    sharedResource: false
    serviceMeshRevision: 'asm-1-22'
    //TODO: This does not work for all regions. Please check the supported regions for the agent pool VM size.
    agentVMSize: 'standard_d2s_v3'
  }
}

// TODO: potentially make unique to cloud
module acr 'br/public:avm/res/container-registry/registry:0.1.1' = {
  name: 'servicehub-${resourcesName}-${location}acrDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub${resourcesName}${location}acr'
    location: rg.outputs.location
    roleAssignments: [
      {
        principalId: aks.outputs.kubeletidentityObjectId
        principalType: 'ServicePrincipal'
        roleDefinitionIdOrName: 'AcrPull'
      }
    ]
  }
}

module workspace 'br/public:avm/res/operational-insights/workspace:0.3.4' = {
  name: 'servicehub-${resourcesName}-workspaceDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-${resourcesName}-workspace'
    location: rg.outputs.location
  }
}

var streams = ['Microsoft-ContainerLogV2']
module dataCollectionRule 'br/public:avm/res/insights/data-collection-rule:0.1.2' = {
  name: 'servicehub-${resourcesName}-data-collection-ruleDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-${resourcesName}-data-collection-rule'
    location: rg.outputs.location
    dataFlows: [
      {
        streams: streams
        destinations: [
          'ciworkspace'
        ]
      }
    ]
    dataSources: {
      extensions: [
        {
          name: 'ContainerInsightsExtension'
          streams: streams
          extensionSettings: {
            dataCollectionSettings: {
              enableContainerLogV2: true
              interval: '1m'
              namespaceFilteringMode: 'Exclude'
            }
          }
          extensionName: 'ContainerInsights'
        }
      ]
    }
    destinations: {
      logAnalytics: [
        {
          workspaceResourceId: workspace.outputs.resourceId
          name: 'ciworkspace'
        }
      ]
    }
  }
}

module serviceBusNamespace 'br/public:avm/res/service-bus/namespace:0.9.0' = {
  name: 'servicehub-${resourcesName}-servicebus-namespaceDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-${resourcesName}-servicebus-namespace'
    location: rg.outputs.location
    queues: [
      {
        name: 'servicehub-${resourcesName}-queue'
      }
    ]
    skuObject: {
      name: 'Basic'
    }
    zoneRedundant: false
  }
}
