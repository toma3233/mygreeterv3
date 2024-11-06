param location string = resourceGroup().location

module test1 './servicehub-generated-module.bicep' = {
  name: 'servicehubtest-clusterDeploy'
  params: {
    name: 'servicehubtest-cluster'
    location: location
    dnsPrefix: 'servicehubtest'
    sharedResource: false
  }
}
