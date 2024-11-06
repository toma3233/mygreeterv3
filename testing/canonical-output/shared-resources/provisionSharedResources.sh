#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#TODO: currently folder name is hard-coded, however in the future we should template this 
#by seperating destination dir to destination dir prefix and shared resources folder name
#in resources-config.yaml
cd testing/canonical-output/shared-resources
make deploy-resources
if [ $? -ne 0 ]
then
    echo -e "${RED}Resource deployment failed $?${NC}"
    exit 1
fi
echo -e "${GREEN}Resource Provisioning was successful!${NC}"