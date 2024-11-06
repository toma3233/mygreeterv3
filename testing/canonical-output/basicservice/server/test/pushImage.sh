#!/bin/bash

#Define color codes for printing to help analysis.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#--------------------------------------------------------
# We are assuming resource provisioning has been complete to deploy service on resources

# TODO: Some sort of check that resources are provisioned
#---------
cd testing/canonical-output/basicservice
cd server
make push-image
if [ $? -ne 0 ]
then
    echo -e "${RED}Docker image push failed with exit code $?${NC}"
    exit 1
fi
echo -e "${GREEN}Docker image push was successfull!${NC}"