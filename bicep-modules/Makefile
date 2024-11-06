
FILE_NAME := servicehub-generated-module.bicep
REGISTRY_NAME := servicehubregistry
BICEP_REGISTRY_DIRECTORY := bicep/modules
DIRECTORY_NAME ?= aks-managed-cluster
# Do NOT use existing/current tag.
# Should you publish a module using the existing/current tag, you will override the current published modules.
# This will cause other users' future deployments to potentially break.
# Increase the version number only
TAG ?= v5

MODULEDIRS := $(patsubst %/,%,$(wildcard */))

testAllModules:
	./testAllModules.sh

# PUBLISHING A MODULE
# i.e. make publishModule DIRECTORY_NAME=aks-managed-cluster
publishModule:
	az bicep publish --file $(DIRECTORY_NAME)/$(FILE_NAME) --target br:$(REGISTRY_NAME).azurecr.io/$(BICEP_REGISTRY_DIRECTORY)/$(DIRECTORY_NAME):$(TAG)
publishAll: $(MODULEDIRS)
	@for dir in $^; do \
		az bicep publish --file $$dir/$(FILE_NAME) --target br:$(REGISTRY_NAME).azurecr.io/$(BICEP_REGISTRY_DIRECTORY)/$$dir:$(TAG); \
	done
