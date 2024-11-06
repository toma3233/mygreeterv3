#!/bin/bash
#Define color codes for printing to help analysis.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#---------
set -e


cd testing/canonical-output/basicservice
cd server

export READPAT=$READPAT
git config --global url."https://$READPAT@dev.azure.com/service-hub-flg/service_hub/_git/service_hub".insteadOf "https://dev.azure.com/service-hub-flg/service_hub/_git/service_hub"

go build ./...
if [ $? -ne 0 ]
then
    echo "${RED}Building the server module failed.${NC}"
    exit 1
fi
go test ./...
if [ $? -ne 0 ]
then
    echo "${RED}Testing the server module failed.${NC}"
    exit 1
fi
echo "${GREEN}Server module build and test was successful.${NC}"
