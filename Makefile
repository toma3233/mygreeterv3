# The file path is relative to repo_root/go-templating
generatorConfig?=../config-files/generator-config.yaml
commonConfig?=../config-files/local-development/external-generation/common-config.yaml
serviceConfig?=../config-files/local-development/external-generation/service-configs
aksServiceConfig?=../config-files/local-development/aks-generation/service-configs
aksCommonConfig?=../config-files/local-development/aks-generation/common-config.yaml
aksEnvInformation?=../config-files/local-development/aks-generation/env-information.yaml
envInformation?=../config-files/local-development/external-generation/env-information.yaml
initServiceYaml?=mygreeterv3-config.yaml
ignoreOldState?=false
threshold?=5
defaultUser?=false
REGISTRY_NAME = servicehubregistry
IMG_NAME = installations
DATE = $(shell date +%Y%m%d)
# -------------------------------------------------------------------------------
genServiceTemplateSpec:
	cd go-templating; go run . genTemplateSpec --generatorConfig ${generatorConfig} --templateChoice mygreeterGoTemplate --defaultUser ${defaultUser}
genResourcesTemplateSpec:
	cd go-templating; go run . genTemplateSpec --generatorConfig ${generatorConfig} --templateChoice resourcesTemplate --defaultUser ${defaultUser}
genPipelineTemplateSpec:
	cd go-templating; go run . genTemplateSpec --generatorConfig ${generatorConfig} --templateChoice pipelineTemplate --defaultUser ${defaultUser}
genBasicServiceTemplateSpec:
	cd go-templating; go run . genTemplateSpec --generatorConfig ${generatorConfig} --templateChoice basicGoTemplate --defaultUser ${defaultUser}
genAllTemplateSpec: genServiceTemplateSpec genResourcesTemplateSpec genPipelineTemplateSpec genBasicServiceTemplateSpec
# -------------------------------------------------------------------------------
generateExternal:
	cd go-templating; go run . generate --generationType ${generationType} --generatorConfig ${generatorConfig} --commonConfig ${commonConfig} --envInformation $(envInformation) --serviceConfig ${serviceConfig} --ignoreOldState=${ignoreOldState}
genPipeline:
	$(MAKE) generateExternal generationType="pipeline"
genResources:
	$(MAKE) generateExternal generationType="resource"
genService:
	$(MAKE) generateExternal generationType="service"
generateAll:
	$(MAKE) genService
	$(MAKE) genResources
	$(MAKE) genPipeline
genAKSService:
	cd go-templating; go run . generate --user "aks" --generationType "service" --generatorConfig ${generatorConfig} --serviceConfig ${aksServiceConfig} --commonConfig ${aksCommonConfig} --ignoreOldState=${ignoreOldState} --envInformation $(aksEnvInformation)
# -------------------------------------------------------------------------------
dumpExternal:
	cd go-templating; go run . dump --generationType ${generationType} --generatorConfig ${generatorConfig} --commonConfig ${commonConfig} --envInformation $(envInformation) --serviceConfig ${serviceConfig}/mygreeterv3-config.yaml
dumpPipeline:
	$(MAKE) dumpExternal generationType="pipeline"
dumpService:
	$(MAKE) dumpExternal generationType="service"
dumpResources:
	$(MAKE) dumpExternal generationType="resource"
dumpOperators:
	$(MAKE) dumpExternal generationType="operators"
dumpAll:
	$(MAKE) dumpService
	$(MAKE) dumpResources
	$(MAKE) dumpPipeline
	$(MAKE) dumpOperators
dumpAKSService:
	cd go-templating; go run . dumpAKSService --generationType "service" --generatorConfig ${generatorConfig} --serviceConfig ${aksServiceConfig}/mygreeterv3-config.yaml --commonConfig ${aksCommonConfig} --envInformation $(aksEnvInformation)
# -------------------------------------------------------------------------------
testOutput:
	bash ./testing/testOutput.sh
testExternalServiceCoverage:
	threshold=${threshold} bash ./testing/canonical-output/mygreeterv3/server/test/testCoverageOutput.sh
testBasicServiceCoverage:
	threshold=${threshold} bash ./testing/canonical-output/basicservice/server/test/testCoverageOutput.sh
testGoTemplatingCoverage:
	threshold=${threshold} StagingDir=../testing/CoverageReports bash ./testing/testCoverageRepo.sh
testCoverage: testExternalServiceCoverage testBasicServiceCoverage testGoTemplatingCoverage
testSuites:
	bash ./testing/testSuites.sh
testAll: testSuites testOutput testExternalServiceCoverage testGoTemplatingCoverage
#This assumes we have an environment variable called subscriptionId
userEnv:
	bash ./user-environment/createInitialUserEnv.sh ${initServiceYaml}
dockerInstallations:
	if ! docker buildx ls | grep -q svchubbuilder; then \
		docker buildx create --name svchubbuilder --driver docker-container --bootstrap --use; \
	else \
		docker buildx use svchubbuilder; \
	fi
	az acr login --name $(REGISTRY_NAME)
	docker buildx build --build-arg AKS_GOPROXY_TOKEN --build-arg GOPROXY --build-arg GOPRIVATE --build-arg GONOPROXY --platform linux/amd64,linux/arm64 --tag $(REGISTRY_NAME).azurecr.io/$(IMG_NAME):$(DATE) --push -f testing/Dockerfile ..
	docker buildx use --builder default
	docker buildx rm --builder svchubbuilder
