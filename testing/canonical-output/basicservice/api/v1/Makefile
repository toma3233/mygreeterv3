all: service
MAIN_DIR := $(shell dirname -- $(shell dirname -- $(shell dirname -- $(shell pwd))))
USER_ID := $(shell id -u):$(shell id -g)
REGISTRY_NAME=servicehubregistry
IMG_NAME=service_hub_environment
IMG_DATE_TAG=20240912
GOPATH := $(shell go env GOPATH)

# api.swagger.json file is automatically generated using npx buf generate command via the openapi plugin
# Set up goproxy variables in accordance to if goModuleNamePrefix is the user's repository or internal go.goms.io
service:
	# If we re-add make service into deploy-resources, these commands will be needed for pipeline to successfully
	# access the private repository
	# if [ -z "$(READPAT)" ]; then \
	#	echo "READPAT does not exist. Please rerun the command with this variable"; \
	#	exit 1; \
	# fi
	# export GOPRIVATE='dev.azure.com'; \
	# git config --global url."https://$(READPAT)@dev.azure.com/service-hub-flg/service_hub/_git/service_hub".insteadOf \
	#	"https://dev.azure.com/service-hub-flg/service_hub/_git/service_hub"; \

	docker run --mount src=$(MAIN_DIR),target=/app/binded-data,type=bind \
		$(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(IMG_DATE_TAG) \
		/bin/bash -c \
	"cd binded-data/basicservice/api/v1; \
	cd proto; buf mod update; cd ..; buf generate; \
	swagger-codegen-cli generate -i api.swagger.json -l go -o restsdk -c swagger-config.json; \
	chown -R $(USER_ID) ../..; \
	cd /app; ./maintain_service populateMethodFiles \
		--protoFilePath binded-data/basicservice/api/v1/proto/api.proto \
		--svcDirName binded-data/basicservice --deleteGarbageFiles=false;"

	# go mod related work is done outside the docker container 
	# to avoid repeated download in the docker container
	go install go.uber.org/mock/mockgen@latest
	go mod tidy # Update go.mod after code generation.
	go get go.uber.org/mock/mockgen/model # Required by mockgen.
	$(GOPATH)/bin/mockgen -package mock \
		dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/basicservice/api/v1 \
		BasicServiceClient > mock/mock_client.go
