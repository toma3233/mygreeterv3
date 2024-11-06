#!/bin/bash

resourceGroupName="servicehubtest_bicep_rg"
location="eastus"

echo "Creating resource group: $resourceGroupName"
az group create --name $resourceGroupName --location $location

resourceGroupInfo=$(az group show --name $resourceGroupName --output json)

if [ -n "$resourceGroupInfo" ]; then
    echo "Resource group '$resourceGroupName' created successfully."
else
    echo "Failed to retrieve information for resource group '$resourceGroupName'."
    exit 1
fi

totalTests=0
passedTests=0
failedTests=()

directories=$(find . -maxdepth 1 -type d)

for dir in $directories; do
    # Skip the current directory and parent directory
    if [ "$dir" == "." ] || [ "$dir" == ".." ]; then
        continue
    fi

    # Define a variable to pass parameters into the test block
    directoryName=$(basename $dir)
    bicepFile="$dir/test.bicep"
    
    ((totalTests++))

    # Validate deployment
    echo "Validating the deployment of $directoryName"
    if grep -q "targetScope = 'subscription'" "$bicepFile"; then
        validationResult=$(az deployment sub validate --location EastUs --template-file $bicepFile)
    else
        validationResult=$(az deployment group validate --resource-group $resourceGroupName --template-file $bicepFile)
    fi

    # Check for errors
    if [[ "$validationResult" == *"\"error\": null"* && "$validationResult" == *"\"provisioningState\": \"Succeeded\""* ]]; then
        echo "Validation passed."
        ((passedTests++))
    else
        echo "Validation failed."
        failedTests+=("$directoryName: $validationResult")
    fi
done

echo
echo "Total tests: $totalTests"
echo "Passed tests: $passedTests"
echo "Percentage of tests passed: $((passedTests * 100 / totalTests))%"

# Display the list of failed tests
if [ ${#failedTests[@]} -gt 0 ]; then
    echo "Failed tests:"
    for failedTest in "${failedTests[@]}"; do
        echo "- $failedTest"
    done
else
    echo "All tests passed!"
fi
