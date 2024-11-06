uuid := $(shell LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 10; echo)
envConfig=../env-config.yaml
location=westus
#This assumes we have an environment variable called subscriptionId
genEnvConfig:
	echo "# A string that will be used for naming all resources that are created for this project." > ${envConfig}
	echo "resourcesName: $(uuid)" >> ${envConfig}
	echo "# The subscription Id to which we deploy our resources." >> ${envConfig}
	echo "subscriptionId: $$subscriptionId" >> ${envConfig}
	echo "# The location of the resource group to which we deploy our resources." >> ${envConfig}
	echo "location: $(location)" >> ${envConfig}
	echo "# The resource group to which we deploy our resources." >> ${envConfig}
	echo "resourceGroupName: servicehub-$(uuid)-rg" >> ${envConfig}
