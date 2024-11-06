#!/bin/bash
# TODO: for long term maintenance, we either need to write in more formal language and test the code, or just use simple test (e.g., get the max length of region name and other var name. Then compute the resource group name length directly.)
set -e
deploymentType=$1
if [[ $deploymentType == "ev2" ]]; then
    # Get the parameters from the scope bindings
    parametersFile="Ev2/Ev2Specs/ScopeBinding.json"
    # Get resource group name from service model
    azureResourceGroupName=$(jq -r '.serviceResourceGroupDefinitions[0].azureResourceGroupName' Ev2/Ev2Specs/ServiceModel.json)
else
    # Get the parameters from the parameters json file
    parametersFile="resources/Main.SharedResources.Parameters.json"
    # Get resource group name from parameters json file
    azureResourceGroupName=$(jq -r '.parameters.resourceGroupName.value' $parametersFile)
fi
# Check if the parameters file exists
if [ ! -f "$parametersFile" ]; then
    echo "Parameters file $parametersFile does not exist."
    exit 1
fi
# Get longest possible azure region that ev2 resources can be deployed in
longestLocation=$(az account list-locations --query "[].name" | jq -r 'max_by(length)')
# If the resource group name contains $location(), replace it with the longest possible azure region
if [[ ${azureResourceGroupName} == *"$location()"* ]]; then
    azureResourceGroupName="${azureResourceGroupName//\$location()/${longestLocation}}"
fi
# If the resource group name contains $config(regionShortName), replace it with the longest possible azure region
if [[ ${azureResourceGroupName} == *"$config(regionShortName)"* ]]; then
    # Currently longest location in region short name has length 5, for example ussw1 which is the short name for southwestus
    longestShortName="ussw1"
    azureResourceGroupName="${azureResourceGroupName//\$config(regionShortName)/${longestShortName}}"
fi
# Check if the resource group name is longer than 90 characters
if [ ${#azureResourceGroupName} -gt 90 ]; then
    echo "Resource group name will exceed the maximum length and lead to ev2 rollout failure."
    exit 1
fi
az bicep build --file "resources/Main.SharedResources.Template.bicep" 
# Get the resource definition for the managed cluster
managedClusterResource=$(jq '.resources[] | select(.type == "Microsoft.Resources/deployments" and .properties.template.resources[]?.type == "Microsoft.ContainerService/managedClusters")' resources/Main.SharedResources.Template.json)
# Get the name of the managed cluster
clusterName=$(jq -r '.resources[] | select(.type == "Microsoft.Resources/deployments") | .properties.template.resources[] | select(.type == "Microsoft.ContainerService/managedClusters") | .name' resources/Main.SharedResources.Template.json)
# If the name of the managed cluster is a parameter, get the value of the parameter
if [[ $clusterName =~ ^\[parameters\(.*\)\]$ ]]; then
    parameterName=$(echo "$clusterName" | sed -e 's/\[parameters(\(.*\))\]/\1/' -e 's/[^a-zA-Z0-9]//g')
    # Get the value of the parameter and remove square brackets
    parameterValue=$(echo "$managedClusterResource" | jq -r --arg pn "$parameterName" '.properties.parameters[$pn].value' | sed 's/\[//g; s/\]//g')
    parameters=$(cat "$parametersFile")
    # Remove everything before the first curly brace
    parameters=$(echo "$parameters" | sed -n '/^{/,$p')
    # If the parameter value is a format string, replace the parameters in the string with their values
    if [[ $parameterValue =~ ^format\(.*\)$ ]]; then
        # Remove format() from around the string
        formattedString=${parameterValue#format(}
        formattedString=${formattedString%)}
        # Split the string by commas
        IFS=',' read -ra stringArray <<"<<<">> "$formattedString"
        firstElement=${stringArray[0]//\'/}
        firstElement=${firstElement//\"/}
        # Separate the first element
        IFS=' ' read -ra firstElementArray <<"<<<">> "$firstElement"
        firstElement=${firstElementArray[0]}
        # Replace each of the parameters in the format string with their corresponding values
        # Example: format('<<.sharedInput.productShortName>>-{0}-cluster', parameters('resourcesName')) -> <<.sharedInput.productShortName>>-<value of resourcesName>-cluster
        # Value is extracted from scope bindings for example -> <<.sharedInput.productShortName>>-official-cluster
        for ((i = 1; i < ${#stringArray[@]}; i++)); do
            # Remove extra spaces, single quotes, double quotes, and parameters() from around the element
            currentElement=${stringArray[i]// /}
            currentElement=${currentElement#parameters(}
            currentElement=${currentElement%)}
            currentElement=${currentElement//\'/}
            currentElement=${currentElement//\"/}
            # Get the index of the parameter
            index=$(($i-1))
            # Get the value of the parameter from the scope bindings or values.json depending on deployment type
            if [[ $deploymentType == "ev2" ]]; then
                evaluatedValue=$(echo "$parameters" | jq -r --arg pname "$currentElement" '.scopeBindings[0].bindings[] | select(.find == "{{." + $pname + "}}") | .replaceWith')
            else
                evaluatedValue=$(echo "$parameters" | jq -r --arg pname "$currentElement" '.parameters[$pname].value')
            fi
            # Replace the parameter in the format string with its value
            firstElement=${firstElement//\{$index\}/${evaluatedValue}}
        done
        # Set the resource name to the formatted string
        resourceName=$firstElement
    else
        # If the parameter value is not a format string, set the resource name to the parameter value
        resourceName=$parameterValue
        echo "parameterValue does not match the format string"
    fi
else
    # If the name of the managed cluster is not a parameter, set the resource name to the cluster name
    resourceName=$clusterName
fi
rm -rf resources/Main.SharedResources.Template.json
# Build AKS node resource group name in accordance to https://learn.microsoft.com/en-us/troubleshoot/azure/azure-kubernetes/create-upgrade-delete/aks-common-issues-faq#what-naming-restrictions-are-enforced-for-aks-resources-and-parameters-
# AKS node resource group name structure is MC_resourceGroupName_resourceName_location
aksNodeRGName="MC_${azureResourceGroupName}_${resourceName}_${longestLocation}"
echo "Worse case scenario will take place if cluster is deployed in the ${longestLocation} region"
echo
if [ ${#aksNodeRGName} -gt 80 ]; then
    echo "Aks node resource group name $aksNodeRGName will exceed the maximum length, as it is at ${#aksNodeRGName} and needs to be less than 80"
    exit 1
else 
    echo "Aks node resource group name $aksNodeRGName is less than maximum length, as it is at ${#aksNodeRGName} and needs to be less than 80"
fi
