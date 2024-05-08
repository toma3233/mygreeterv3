function azLogin() {
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    NC='\033[0m'
    #--------------------------------------------------------
    #Login with provided service principal and provision resources
    az login --service-principal --username $AZ_APP_ID --tenant $AZ_TENANT_ID --password $AZ_PASSWORD
    if [ $? -ne 0 ]
    then
        echo -e "${RED}Azure login was not successfull${NC}"
        exit 1
    fi
    echo -e "${GREEN}Azure login was successfull${NC}"
}