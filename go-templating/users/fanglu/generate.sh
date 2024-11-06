#!/bin/sh -x

# The file path is relative to the go-templating directory.
# The Makefile is at the root of the repo.

# Get to the root of the repo.
cd ../../.. 

# make genServiceTemplateSpec

# Generate the code.
make generate serviceConfig=users/fanglu/service-config.yaml middlewareConfig=users/fanglu/middleware-config.yaml

# Hardcoded to be consistent with what defined in middleware-config.yaml and service-config.yaml.

# Init the aks-middleware module.
cd ~/aks-rp/aks-middleware
make

# Init the service.
cd ~/aks-rp/mygreeterv2

# Only needed when use one's own aks-middleware copy.
# No need to do this when using a shared module.
go work init
go work use ../aks-middleware

./init.sh

mv go.work go.work.temp

make bootstrap
make tidy
make build
make deploy
