#!/bin/bash
source ./testing/canonical-output/mygreeterv3/server/test/loginWithSP.sh

#Define color codes for printing to help analysis.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#--------------------------------------------------------
#Login with provided service principal and provision resources
azLogin || exit 1
# We are assuming resource provisioning has been complete to deploy service on resources

# TODO: Some sort of check that resources are provisioned
#---------
cd testing/canonical-output/mygreeterv3
cd server
make deploy-resources
if [ $? -ne 0 ]
then
    echo -e "${RED}Provisioning service specific resources failed${NC}"
    exit 1
fi
echo -e "${GREEN}Provisioning service specific resources failed was successful!${NC}"