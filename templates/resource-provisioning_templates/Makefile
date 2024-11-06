DIR_NAME=<<.resourceInput.directoryName>>
MAIN_DIR := $(shell dirname -- $(shell pwd))
#The config file used to populate resourceName and subscriptionId, default assumes there is a config file
#in the main directory of project.
ENVCONFIG?="env-config.yaml"
REGISTRY_NAME=servicehubregistry
IMG_NAME=service_hub_environment
IMG_DATE_TAG=20240912

.PHONY: template-files deploy-resources

template-files:
	@docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) \
	./maintain_service executeTemplatesInFolder --templatePrefix "template-" --templatePath binded-data/$(DIR_NAME) --envConfig binded-data/$(ENVCONFIG) --leftDelim "{{" --rightDelim "}}"

deploy-resources: template-files
	echo "Starting resource provisioning."
	RESOURCES_NAME=$(shell jq -r '.parameters.resourcesName.value' resources/Main.SharedResources.Parameters.json); \
	az deployment sub create --name "<<.sharedInput.productShortName>>-$$RESOURCES_NAME-deploy" --location EastUS --template-file resources/Main.SharedResources.Template.bicep --parameters resources/Main.SharedResources.Parameters.json -o json > .$(DIR_NAME)_tmp.json
	@if cat .$(DIR_NAME)_tmp.json | grep '"provisioningState": "Succeeded"' > /dev/null 2>&1; then \
		echo "Resource provisioning succeeded."; \
		cp .$(DIR_NAME)_tmp.json .$(DIR_NAME)_output.json; \
	else \
		echo "Resource provisioning did not succeed."; \
		cat .$(DIR_NAME)_tmp.json; \
		exit 1; \
	fi
	@rm .$(DIR_NAME)_tmp.json
	@docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) ./maintain_service generateResourceFiles --svcDirName=binded-data/$(DIR_NAME); \

