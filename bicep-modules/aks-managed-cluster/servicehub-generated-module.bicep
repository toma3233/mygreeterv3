@sys.description('The name of the Managed Cluster resource.')
param name string

@sys.description('The location of the Managed Cluster resource.')
param location string = resourceGroup().location

@sys.description('Optional DNS prefix to use with hosted Kubernetes API server FQDN.')
param dnsPrefix string = ''

@sys.description('Disk size (in GB) to provision for each of the agent pool nodes. This value ranges from 0 to 1023. Specifying 0 will apply the default disk size for that agentVMSize.')
@minValue(0)
@maxValue(1023)
param osDiskSizeGB int = 0

@sys.description('The number of nodes for the cluster.')
@minValue(1)
@maxValue(50)
param agentCount int = 3

@sys.description('The size of the Virtual Machine.')
param agentVMSize string = 'standard_d2s_v3'

@sys.description('Optional. Specifies whether the OMS agent is enabled.')
param omsAgentEnabled bool = true

@sys.description('Optional. Resource ID of the monitoring log analytics workspace.')
param monitoringWorkspaceId string = ''

@sys.description('Optional. Specifies whether to create the cluster as a private cluster or not.')
param enablePrivateCluster bool = false

@sys.description('Optional. Whether to enable Workload Identity. Requires OIDC issuer profile to be enabled.')
param enableWorkloadIdentity bool = false

// not in public cluster module
@sys.description('Specify true to create a data collection rule association with the cluster, or false to not create a data collection rule association.')
param createDataCollectionRuleAssociation bool = false

// not in public cluster module
@sys.description('Full Resource ID of data collection rule that is used to specify data collection from AKS Cluster. Only needed if createDataCollectionRuleAssociation is true.  For example /subscriptions/00000000-0000-0000-0000-0000-00000000/resourceGroups/ResourceGroupName/providers/Microsoft.Insights/dataCollectionRules/DataCollectionRuleName.')
param dataCollectionRuleId string = ''

@sys.description('Specify true if the resource is a shared resource, or false if the resource is not shared.')
param sharedResource bool = false

@sys.description('Optional. Specifies whether the Istio-based service mesh add-on is enabled or not. Cannot be enabled when osm is enabled.')
param enableIstioServiceMesh bool = false

@description('Optional. Version of Kubernetes specified when creating the managed cluster.')
param kubernetesVersion string = ''

@description('Optional. The default service mesh revision. Used when enableIstioServiceMesh is true.')
param serviceMeshRevision string = 'asm-1-19'

resource managedClusterShared 'Microsoft.ContainerService/managedClusters@2023-07-02-preview' existing = if (sharedResource) {
  name: name
}

resource managedCluster 'Microsoft.ContainerService/managedClusters@2023-07-02-preview' = if (!(sharedResource)) {
  name: name
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    dnsPrefix: dnsPrefix
    kubernetesVersion: (empty(kubernetesVersion) ? null : kubernetesVersion)
    agentPoolProfiles: [
      {
        name: 'agentpool'
        osDiskSizeGB: osDiskSizeGB
        count: agentCount
        vmSize: agentVMSize
        osType: 'Linux'
        mode: 'System'
      }
    ]
    addonProfiles: {
      omsagent: {
        enabled: omsAgentEnabled && !empty(monitoringWorkspaceId)
        config: omsAgentEnabled && !empty(monitoringWorkspaceId) ? {
          logAnalyticsWorkspaceResourceID: !empty(monitoringWorkspaceId) ? any(monitoringWorkspaceId) : null
          useAADauth: 'true'
        } : null
      }
    }
    oidcIssuerProfile: enableWorkloadIdentity ? {
      enabled: enableWorkloadIdentity
    } : null
    securityProfile: {
      workloadIdentity: enableWorkloadIdentity ? {
        enabled: enableWorkloadIdentity
      } : null
    }
    serviceMeshProfile: enableIstioServiceMesh ? {
      istio: {
        certificateAuthority: null
        components: {
          egressGateways: null
          ingressGateways: null
        }
        revisions: [
          serviceMeshRevision
        ]
      }
      mode: 'Istio'
    } : null
  }
}

resource dataCollectionRuleAssociation 'Microsoft.Insights/dataCollectionRuleAssociations@2022-06-01' = if (createDataCollectionRuleAssociation) {
  name: 'ContainerInsightsExtension'
  scope: managedCluster // requires association with an existing resource
  properties: {
    dataCollectionRuleId: dataCollectionRuleId
    description: '${name} Data Collection Rule Association'
  }
}

var selectedManagedCluster = sharedResource ? managedClusterShared : managedCluster

@sys.description('The resource ID of the managed cluster.')
output resourceId string = sharedResource ? managedClusterShared.id : managedCluster.id

@sys.description('The resource group the managed cluster was deployed into.')
output resourceGroupName string = resourceGroup().name

@sys.description('The name of the managed cluster.')
output name string = name

@sys.description('The control plane FQDN of the managed cluster.')
output controlPlaneFQDN string = enablePrivateCluster ? selectedManagedCluster.properties.privateFQDN : selectedManagedCluster.properties.fqdn

@sys.description('The Object ID of the AKS identity.')
output kubeletidentityObjectId string = contains(selectedManagedCluster.properties, 'identityProfile') ? contains(selectedManagedCluster.properties.identityProfile, 'kubeletidentity') ? selectedManagedCluster.properties.identityProfile.kubeletidentity.objectId : '' : ''

@sys.description('The Object ID of the OMS agent identity.')
output omsagentIdentityObjectId string = contains(selectedManagedCluster.properties, 'addonProfiles') ? contains(selectedManagedCluster.properties.addonProfiles, 'omsagent') ? contains(selectedManagedCluster.properties.addonProfiles.omsagent, 'identity') ? selectedManagedCluster.properties.addonProfiles.omsagent.identity.objectId : '' : '' : ''

@sys.description('The Object ID of the Key Vault Secrets Provider identity.')
output keyvaultIdentityObjectId string = contains(selectedManagedCluster.properties, 'addonProfiles') ? contains(selectedManagedCluster.properties.addonProfiles, 'azureKeyvaultSecretsProvider') ? contains(selectedManagedCluster.properties.addonProfiles.azureKeyvaultSecretsProvider, 'identity') ? selectedManagedCluster.properties.addonProfiles.azureKeyvaultSecretsProvider.identity.objectId : '' : '' : ''

@sys.description('The Client ID of the Key Vault Secrets Provider identity.')
output keyvaultIdentityClientId string = contains(selectedManagedCluster.properties, 'addonProfiles') ? contains(selectedManagedCluster.properties.addonProfiles, 'azureKeyvaultSecretsProvider') ? contains(selectedManagedCluster.properties.addonProfiles.azureKeyvaultSecretsProvider, 'identity') ? selectedManagedCluster.properties.addonProfiles.azureKeyvaultSecretsProvider.identity.clientId : '' : '' : ''

@sys.description('The location the resource was deployed into.')
output location string = selectedManagedCluster.location

@sys.description('The OIDC token issuer URL.')
output oidcIssuerUrl string = selectedManagedCluster.properties.oidcIssuerProfile.enabled ? selectedManagedCluster.properties.oidcIssuerProfile.issuerURL : ''

@sys.description('The addonProfiles of the Kubernetes cluster.')
output addonProfiles object = contains(selectedManagedCluster.properties, 'addonProfiles') ? selectedManagedCluster.properties.addonProfiles : {}
