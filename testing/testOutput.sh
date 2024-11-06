#!/bin/bash -x
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
#Set config files
servConfigFolder="../config-files/testing/service-configs"
genConfig="../config-files/generator-config.yaml"
commonConfig="../config-files/testing/common-config.yaml"
terraformOut="../config-files/testing/env-information.yaml"
#Generate pipeline, service and shared resources
rm -rf testing/generated-output
mkdir testing/generated-output
make generateAll serviceConfig=$servConfigFolder commonConfig=$commonConfig generatorConfig=$genConfig envInformation=$terraformOut
if [ $? -ne 0 ]
then
    echo "Test Failed: Generation failed."
    exit 1
fi
echo "Test Passed: Generation was successful"

# Generate mygreeterv3
cd testing/generated-output/mygreeterv3
./init.sh
if [ $? -ne 0 ]
then
    echo "Test Failed: Init.sh failed."
    exit 1
fi
cd ..
cd ..
cd ..

# Generate basicservice
cd testing/generated-output/basicservice
./init.sh
if [ $? -ne 0 ]
then
    echo "Test Failed: Init.sh failed."
    exit 1
fi
cd ..
cd ..
cd ..

#Compare difference between canonical outputs and currently generated outputs
#However ignore all .mod and .sum files as they change upon each run.
diff -r testing/generated-output testing/canonical-output  -x '*.mod' -x '*.sum' -x '*.work' -x 'buf.lock'
if [ $? -eq 0 ]
then
        diff -r testing/generated-output testing/canonical-output
        if [ $? -eq 0 ]
        then
                echo -e "${GREEN}Test Passed: Generated output is same as canonical output"
        else
                echo -e "${GREEN}Test Passed: Generated output is same as canonical output${NC}"
                echo -e "${RED}However, it as shown above, the mod/sum/work/lock files need to be copied over in order for the canonical output to be the latest version.${NC}"
        fi
else
        echo "Test Failed: Generated output is different from canonical output"
        exit 1
fi
