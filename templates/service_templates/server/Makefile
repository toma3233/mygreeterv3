<<if eq .user "external">>
REGISTRY_NAME=servicehubregistry
IMG_NAME=service_hub_environment
IMG_DATE_TAG=20240912
MAIN_DIR := $(shell dirname -- $(shell dirname -- $(shell pwd)))
DIR_NAME=<<.serviceInput.directoryName>>
#The config file used to populate resourceName and subscriptionId, default assumes there is a config file
#in the main directory of project.
ENVCONFIG?="env-config.yaml"

all: tidy test build

template-files:
	@if [ -f ../../$(ENVCONFIG) ]; then \
		docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) \
		./maintain_service executeTemplatesInFolder --templatePrefix "template-" --templatePath binded-data/$(DIR_NAME) --envConfig binded-data/$(ENVCONFIG) --leftDelim "{{" --rightDelim "}}"; \
		echo "Templating server files was successful"; \
	else \
		echo "../../$(ENVCONFIG) does not exist. Please make sure either the path exists or rerun the command with the flag --ENVCONFIG 'pathtoconfigfile'"; \
		exit 1; \
	fi

# (Prerequisite) To deploy the shared resources needed, go to the directory of the shared resources and run the Makefile.
# Creates all the service specific resources and generates markdown files
# TODO: Investigate how to express dependency across different directories rather than write imperative code.
deploy-resources: template-files
	cd resources; make all
	docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) ./maintain_service generateResourceFiles --svcDirName=binded-data/$(DIR_NAME); \

tidy: template-files
	go mod tidy

test:
	go test ./...

build: template-files
	go build -o bin/ <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/cmd/client
	go build -o bin/ <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/cmd/demoserver
	go build -o bin/ <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/cmd/server
	
build-image: template-files
	cd generated; make build-image

push-image: template-files
	cd generated; make push-image


AZURESDK_PATH := artifacts/.Azuresdk_properties_outputs.yaml

# The install and uninstall information (k8s secret) is stored in the current namespace.
install: template-files
	cd generated; make connect-cluster
	if [ -f $(AZURESDK_PATH) ]; then \
		docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) /bin/bash -c \
		"cd binded-data/$(DIR_NAME)/server; export KUBECONFIG=\$$(pwd)/generated/.kube/config; \
		cd deployments; helm install <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-client . -f values-common.yaml -f values-client.yaml; \
		helm install <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server . -f values-common.yaml -f values-server.yaml -f ../$(AZURESDK_PATH); \
		helm install <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-demoserver . -f values-common.yaml -f values-demoserver.yaml;"; \
	else \
		echo ".Azuresdk_properties_outputs.yaml file does not exist. Run the make deploy-resources cmd"; \
		docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) /bin/bash -c \
		"cd binded-data/$(DIR_NAME)/server; export KUBECONFIG=\$$(pwd)/generated/.kube/config; \
		cd deployments; helm install <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-client . -f values-common.yaml -f values-client.yaml; \
		helm install <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server . -f values-common.yaml -f values-server.yaml; \
		helm install <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-demoserver . -f values-common.yaml -f values-demoserver.yaml;"; \
	fi

upgrade: template-files
	cd generated; make connect-cluster
	if [ -f $(AZURESDK_PATH) ]; then \
		docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) /bin/bash -c \
		"cd binded-data/$(DIR_NAME)/server; export KUBECONFIG=\$$(pwd)/generated/.kube/config; \
		cd deployments; helm upgrade <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-client . -f values-common.yaml -f values-client.yaml; \
		helm upgrade <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server . -f values-common.yaml -f values-server.yaml -f ../$(AZURESDK_PATH); \
		helm upgrade <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-demoserver . -f values-common.yaml -f values-demoserver.yaml;"; \
	else \
		echo ".Azuresdk_properties_outputs.yaml file does not exist. Run the make deploy-resources cmd"; \
		docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) /bin/bash -c \
		"cd binded-data/$(DIR_NAME)/server; export KUBECONFIG=\$$(pwd)/generated/.kube/config; \
		cd deployments; helm upgrade <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-client . -f values-common.yaml -f values-client.yaml; \
		helm upgrade <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server . -f values-common.yaml -f values-server.yaml; \
		helm upgrade <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-demoserver . -f values-common.yaml -f values-demoserver.yaml;"; \
	fi

uninstall: template-files
	docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) /bin/bash -c \
	"cd binded-data/$(DIR_NAME)/server; export KUBECONFIG=\$$(pwd)/generated/.kube/config; \
	helm uninstall <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-client; \
	helm uninstall <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-demoserver; \
	helm uninstall <<.sharedInput.productShortName>>-<<.serviceInput.directoryName>>-server;";
<<else>>
# Assumptions: 
# The directory of this file is under aks-rp. This is to find aksbuilder.sh.
# The aks-rp directory is under $HOME. The azureconfig.yaml is at aks-rp/bin/azureconfig.yaml.
# This is to find azureconfig.yaml used by deploy.
# Modify the path in this file if the assumption of this file is broken.
bootstrap:
	../../hack/aksbuilder.sh bootstrap -w ./

tidy:
	../../hack/aksbuilder.sh tidy -w ./

build:
	../../hack/aksbuilder.sh build -w ./

deploy:
	../../hack/aksbuilder.sh localdeploy --workspace ./ --usespoperator --e2eazureconfig ~/aks-rp/bin/azureconfig.yaml
<<end>>
