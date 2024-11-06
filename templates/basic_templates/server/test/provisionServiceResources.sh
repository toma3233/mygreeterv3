#!/bin/bash

#Define color codes for printing to help analysis.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#--------------------------------------------------------
# We are assuming resource provisioning has been complete to deploy service on resources

# TODO: Some sort of check that resources are provisioned
#---------
cd <<.sharedInput.directoryPath>><<.serviceInput.directoryName>>
cd server
# If we re-add make service into deploy-resources, these arguments will 
# be needed for pipeline to successfully access the private repository.
<<- if contains .envInformation.goModuleNamePrefix "dev.azure.com">>
# export READPAT=$1
<<- else>>
# export AKS_GOPROXY_TOKEN=$1
# export GOPROXY=$2
# export GOPRIVATE=$3
# export GONOPROXY=$4<<end>>
make deploy-resources
if [ $? -ne 0 ]
then
    echo -e "${RED}Provisioning service specific resources failed${NC}"
    exit 1
fi
echo -e "${GREEN}Provisioning service specific resources failed was successful!${NC}"