targetScope = 'subscription'

@sys.description('The name for the resources.')
param resourcesName string

@sys.description('The subscription the resources are deployed to.')
param subscriptionId string

var resourceGroupName = 'servicehub-${resourcesName}-rg'

// This resource is shared and defined in main.bicep in shared-resources directory; we only reference it here. Do not remove `existing` syntax.
resource logAnalyticsWorkspace 'Microsoft.OperationalInsights/workspaces@2022-10-01' existing = {
  name: 'servicehub-${resourcesName}-workspace'
  scope: resourceGroup(subscriptionId, resourceGroupName)
}

module qpsAlertRule 'br/public:avm/res/insights/scheduled-query-rule:0.1.2' = {
  name: 'mygreeterv3-CreateStorageAccount-query-per-second'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'mygreeterv3-CreateStorageAccount-query-per-second'
    alertDescription: 'This is the alert for mygreeterv3 for the calculated metric query per second.'
    criterias: {
      allOf: [
        {
          metricMeasureColumn:'QPS'
          operator: 'GreaterThan'
          query:'let method = "CreateStorageAccount";\nlet binSizeMinute = timespan(5m);\nlet binSizeSecond = binSizeMinute / 1s;\nContainerLogV2\n| where ContainerName == "servicehub-mygreeterv3-server"\n| where LogMessage["component"] == "server"\n| where LogMessage["method"] == method or isempty(method)\n| where LogMessage["msg"] == "finished call"\n| summarize QPS = count()/binSizeSecond by tostring(LogMessage["code"]), bin(todatetime(LogMessage["time"]), binSizeMinute)\n\n'
          threshold: '0.05'
          timeAggregation: 'Maximum'
        }
      ]
    }
    enabled: true
    windowSize: 'PT5M'
    evaluationFrequency: 'PT5M'
    severity: 4
    autoMitigate: true
    scopes: [
      logAnalyticsWorkspace.id
    ]
  }
}

module errorRatioAlertRule 'br/public:avm/res/insights/scheduled-query-rule:0.1.2' = {
  name: 'mygreeterv3-CreateStorageAccount-error-ratio'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'mygreeterv3-CreateStorageAccount-error-ratio'
    alertDescription: 'This is the alert for mygreeterv3 for the calculated metric error ratio.'
    criterias: {
      allOf: [
        {
          metricMeasureColumn:'ratio'
          operator: 'GreaterThan'
          query: 'let method = "CreateStorageAccount";\nContainerLogV2\n| where ContainerName == "servicehub-mygreeterv3-server"\n| where LogMessage["component"] == "server"\n| where LogMessage["method"] == method\n| where LogMessage["msg"] == "finished call"\n| extend timePoint = todatetime(LogMessage["time"])\n| extend code = tostring(LogMessage["code"])\n| summarize \n total = count(),\nerror = countif(code != "OK") \nby bin(timePoint,5m)\n| extend ratio = round(error * 100.0/total, 3)\n| project timePoint, ratio\n\n'
          threshold: '0.05'
          timeAggregation: 'Maximum'
        }
      ]
    }
    enabled: true
    windowSize: 'PT5M'
    evaluationFrequency: 'PT5M'
    severity: 4
    autoMitigate: true
    scopes: [
      logAnalyticsWorkspace.id
    ]
  }
}

module latencyAlertRule 'br/public:avm/res/insights/scheduled-query-rule:0.1.2' = {
  name: 'mygreeterv3-CreateStorageAccount-latency-by-error-code'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'mygreeterv3-CreateStorageAccount-latency-by-error-code'
    alertDescription: 'This is the alert for mygreeterv3 for the calculated metric error ratio.'
    criterias: {
      allOf: [
        {
          metricMeasureColumn:'latency'
          operator: 'GreaterThan'
          query: 'let method = "CreateStorageAccount";\nlet binSizeMinute = timespan(5m);\nlet binSizeSecond = binSizeMinute / 1s;\nContainerLogV2\n| where ContainerName == "servicehub-mygreeterv3-server"\n| where LogMessage["component"] == "server"\n| where LogMessage["method"] == method\n| where LogMessage["msg"] == "finished call"\n| summarize latency = avg(todouble(LogMessage["time_ms"]))\n    by\n    tostring(LogMessage["code"]),\n    bin(todatetime(LogMessage["time"]), binSizeMinute)\n\n'
          threshold: '1600'
          timeAggregation: 'Maximum'
        }
      ]
    }
    enabled: true
    windowSize: 'PT5M'
    evaluationFrequency: 'PT5M'
    severity: 4
    autoMitigate: true
    scopes: [
      logAnalyticsWorkspace.id
    ]
  }
}