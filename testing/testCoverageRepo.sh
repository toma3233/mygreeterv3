#!/bin/bash -x
source ./testing/runCoverage.sh
#Define color codes for printing to help analysis.
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#Add go path to bash
PATH=$PATH:$(go env GOPATH)/bin
#Create directory within staging directory on agent machine
cd go-templating
mkdir -p $StagingDir
go install github.com/onsi/ginkgo/v2/ginkgo@latest
fileName=go-templating-coverage-report
#Perform test coverage for given folder and save output result
coverage $StagingDir $fileName
#Exit if coverage results were unsatisfactory.
if [ $? -eq 0 ]
then
    exit 1
fi