# Navigating service hub through its use cases and e2e flows

## What are the current configuration files available/required for generation

- For more information about generation capabilities refer to [generation-capabilities.md](generation-capabilities.md)

### For all users (aks and external)

- generator-config.yaml:
  - The generator config is the config file that defines where the template folders exist. Generating will take the files from these template locations.
  - Currently all generation uses the same [generator-config.yaml](../config-files/generator-config.yaml)
- env-information.yaml:
  - For the user environment this file would be automatically generated, by the script that creates the user environment (extracting the value from terraform output) however for our regular development and testing it is manually generated as there is no terraform resources to extract the output from.
- service-config.yaml:
  - Used to define the variables that are used when generating the service.
- common-config.yaml
  - Used to define the shared variables used for generation.

### For external users only

- pipelineInput:
  - Used to define the variables that are used when generating the pipeline files.
- resourceInput:
  - Used to define the variables that are used when generating the resources module.

# 1. AKS generation development

## a. Code stored within the aks-rp repository

### i. Generating code locally (make genAKSService)

- Config files exist under [config-files/local-development/aks-generation](../config-files/local-development/aks-generation).
- The following config files exist:
  - same [generator-config.yaml](../config-files/generator-config.yaml)
  - [service-configs folder](../config-files/local-development/aks-generation/service-configs/)
    - [mygreeterv3-config.yaml](../config-files/local-development/aks-generation/service-configs/mygreeterv3-config.yaml)
  - [env-information.yaml](../config-files/local-development/aks-generation/env-information.yaml)
  - [common-config.yaml](../config-files/local-development/aks-generation/common-config.yaml)

Variables defined in config files relevant to this use case.

- Where code is generated (defined in common-config.yaml)
  - destinationDirPrefix: ~/projects/aks-rp  
- Which goModuleNamePrefix (defined in env-information.yaml) is used for both api and server modules (hardcoded by us since no terraform resources are created)
  - goModuleNamePrefix: go.goms.io/aks/rp
  - When go.work file exists (which is default set up), local api module will be used. When go.work is removed, the api module in the repo is used. However, docker image build and aksbuilder operations does not have access to go.work file so the api module in the repo is used.

- Path from root of repo to directory (blank as root of repo contains the directories that are generated). There is no directories between root of repo and service directory
  - directoryPath: ""

### ii. Using local api module vs api module stored in repository (aks-rp)

Refer to [Using local api module vs api module stored in repository](#ii-using-local-api-module-vs-api-module-stored-in-repository) for more information.

### iii. Testing code manually (aks-rp)

- To test the go modules manually:
  - init.sh in serviceDirectoryName will rerun go build and test.
- Testing on aks-standalone
  - Push changes to aks-rp
  - Deploy [standalone pipeline](https://msazure.visualstudio.com/CloudNativeCompute/_build?definitionId=68881)
  - Complete e2e flow with steps detailed in Run service on AKS internal standalone section of [README.md](../templates/service_templates/README.md).
  - For more information or help use this wiki detailing: [AKS Standalone Environment Usage (Internal Customer)](https://service-hub-flg.visualstudio.com/service_hub/_wiki/wikis/service_hub.wiki/26/AKS-Standalone-Environment-Usage-(Internal-Customer))

### iv. How to merge your changes into the repository (aks-rp)

- Create a pull request for your branch in aks-rp, and follow standard procedures for aks-rp.

## b. Possible future use for internal AKS users

- AKS generated service has capability to use service hub created resources. I.E. aks repo switches to deploying resources via bicep.

# 2. External generation development

## a. Code stored within the service hub repository

This use case is when service hub is used to generate non-aks version of code for development, but that development uses modules stored at the service hub repository when local api is not being used.

### i. Generating code locally (make generateAll)

- Config files exist under [config-files/local-development/external-generation](../config-files/local-development/external-generation/external-generation).
- The following config files exist:
  - same [generator-config.yaml](../config-files/generator-config.yaml)
  - [service-configs folder](../config-files/local-development/external-generation/service-configs/)
    - [mygreeterv3-config.yaml](../config-files/local-development/external-generation/service-configs/mygreeterv3-config.yaml)
    - [readme-config.yaml](../config-files/local-development/external-generation/service-configs/readme.yaml)
  - [common-config.yaml](../config-files/local-development/external-generation/common-config.yaml)
  - [env-information.yaml](../config-files/local-development/external-generation/env-information.yaml)

Variables defined in config files relevant to this use case.

- Where code is generated  (defined in common-config.yaml)
  - destinationDirPrefix: ~/projects/external
- Which goModuleNamePrefix  (defined in env-information.yaml) is used for both api and server modules (hardcoded by us since no terraform resources are created)
  - goModuleNamePrefix: dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output
  - Please refer to the next section for more information about the goModuleNamePrefix use for api module.

- Path from root of repo to directory (blank as root of folder will be projects/external)
  - directoryPath: ""

### ii. Using local api module vs api module stored in repository

In the previous section, it was mentioned that the goModuleNamePrefix defines the api module's prefix and is referencing a module that is tagged within the service hub repository. However, there are technically 2 "versions" of the api module.

- Using local api
  - Currently, the code in the service is already set up to use the local api module in most locations. The go.work file at the root of the service directory forces the server module to use the local api module.
  - However, the [Dockerfile](../templates/service_templates/server/Dockerfile) that builds the service image is the one location that does not use the local api module, and it instead pulls the external api module defined by the goModuleNamePrefix. If you want it to use local api, you have to change the dockerfile to include all the files from serviceDirectoryName not just serviceDirectoryName/server.
- Using the external module
  - If you want the service to only use the external api module, in the [go.work](../testing/canonical-output/mygreeterv3/go.work) file that gets created under the service, you must remove ./api. This way it will pull the external api module.

### iii. Testing code manually

- To test the go modules manually:
  - init.sh will rerun go build and test.
- To test the whole e2e process manually:
  - Follow the "quick setup steps" mentioned in the README.md at the root of your generated directory (~/projects/external)
- Refer to section [v. Testing the code with the service hub pipeline (make testOutput)](#v-testing-the-code-with-the-service-hub-pipeline-make-testoutput) below to see how you can test more robustly

### iv. How to merge your changes into the repository

- If you are a service hub developer please refer to the [developer guide](../developer-guide.md) as it details an example e2e flow of how to make changes and how to merge into repository.
- If you are an external user:
  - Within the generated folder (~/projects/external), initialize git and set up git for the repository you would like to push this code to.

### v. Testing the code with the service hub pipeline (make testOutput)

This explains how to test non-aks generation within service hub with pipelines (test and dev pipeline). As a result, it performs non-aks generation and uses the module name that accesses modules tagged in service hub repository.

- Config files exist under [config-files/testing](../config-files/testing/).
  - same [generator-config.yaml](../config-files/generator-config.yaml)
  - [service-configs folder](../config-files/testing/service-configs/)
    - [mygreeterv3-config.yaml](../config-files/testing/service-configs/mygreeterv3-config.yaml)
    - [readme-config.yaml](../config-files/testing/service-configs/readme.yaml)
  - [common-config.yaml](../config-files/testing/common-config.yaml)
  - [env-information.yaml](../config-files/testing/env-information.yaml)

Variables defined in config files relevant to this use case.

- Where code is generated  (defined in common-config.yaml)
  - destinationDirPrefix: ../testing/generated-output
- Which goModuleNamePrefix (defined in env-information.yaml) is used for both api and server modules (hardcoded by us since no terraform resources are created)
  - goModuleNamePrefix: dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output
  - Although generation occurs in generated-output, the goModuleNamePrefix references canonical-output, this is because service hub uses the idea of "golden tests". Where we compare what was created due to changes (generated-output) to the last fully functioning version (canonical-output). This is explained further in the how to test section directly below.
- Path from root of repo to directory (i.e. from service hub to canonical output to give pipelines the right path)
  - directoryPath: testing/canonical-output/

How to test if you are a service hub developer. (Refer to the [developer guide](../developer-guide.md) for more in depth info)

- In root of service hub repo run

    ```bash
    make testOutput
    ```

- The above command compares generated-output content with canonical output as developers need to make sure they didn't accidentally change something or break something in generation. Take a look at the diff check to determine if everything changed was as expected.
- Update the canonical-output directory. Service hub engineers need to make sure that generated-output matches canonical output as the pipeline build uses canonical-output.

    ```bash
    rm -rf testing/canonical-output/*; cp -Rp testing/generated-output/* testing/canonical-output/
    ```

- Push your code to service hub, and run [our test pipeline](https://dev.azure.com/service-hub-flg/service_hub/_build?definitionId=8) either through your pull request or manually in pipelines.

Structure for code generated within testing directory <https://dev.azure.com/service-hub-flg/_git/service_hub?path=/testing/canonical-output>:

- testing
  - canonical-output (copy of most updated generated repo)
    - serviceDirectoryName (i.e. mygreeterv3)
    - resourcesDirectoryName (i.e shared-resources)
    - pipelineDirectoryName (i.e. pipeline-files)
  - generated-output (in git ignore so not checked in)
    - serviceDirectoryName (i.e. mygreeterv3)
    - resourcesDirectoryName (i.e shared-resources)
    - pipelineDirectoryName (i.e. pipeline-files)

What does pipeline look for?

- As mentioned above, our test and dev pipelines use canonical output to mimic a fully functioning external user repository, as a result, the pipeline looks for [testing/canonical-output/pipeline-files/testServiceResourceAndCode.yaml](../testing/canonical-output/pipeline-files/testServiceResourceAndCode.yaml)
  - Both the [Resource Provisioning Dev Pipeline](https://dev.azure.com/service-hub-flg/service_hub/_build?definitionId=60) and the [Resource Provisioning Test Pipeline](https://dev.azure.com/service-hub-flg/service_hub/_build?definitionId=8) use this file as the starting yaml template.

When are pipelines triggered?

- Test pipeline:
  - Can be manually triggered
  - Must be manually triggered and passed as a part of PR check before merging to master. At the top of each PR, it is provided as a required check to pass before merging is possible. If you do not click queue to trigger the pipeline build run and it doesn't pass, the repository policy will not allow the PR to merge into master.
- Dev pipeline:
  - Only manual trigger

## b. Code stored within user repo generated by service hub

### i. Generating the code for the external repository (make userEnv)

This use case is when service hub creates an ado project and repository, and generates code to be pushed to that repository.

- Config files exist under [config-files/user-environment](../config-files/user-environment/).
- The following config files exist:
  - same [generator-config.yaml](../config-files/generator-config.yaml)
  - [service-configs folder](../config-files/user-environment/service-configs/)
    - [mygreeterv3-config.yaml](../config-files/user-environment/service-configs/mygreeterv3-config.yaml)
    - [readme-config.yaml](../config-files/user-environment/service-configs/readme.yaml)
  - [common-config.yaml](../config-files/user-environment/common-config.yaml)
  - [env-information.yaml](../config-files/user-environment/env-information.yaml)

Variables defined in config files relevant to this use case.

- Where code is generated in service hub (defined in common-config.yaml)
  - destinationDirPrefix: ../user-environment/generated-output
- Which goModuleNamePrefix (defined in terraform.yaml) is used for both api and server modules (Created by the output of terraform, the following is an example with temp_userenv)
  - goModuleNamePrefix: dev.azure.com/service-hub-flg/temp_userenv/_git/temp_userenv_repo.git
- Path from root of repo to directory (blank as root of folder will be where generated repo goes)
  - directoryPath: ""

Repository structure once code is pushed to created repository:
example : <https://dev.azure.com/service-hub-flg/temp_userenv/_git/temp_userenv_repo>

- serviceDirectoryName (i.e. mygreeterv3)
- resourcesDirectoryName (i.e shared-resources)
- pipelineDirectoryName (i.e. pipeline-files)
- terraform-files

Refer to this section: [external environment pipeline testing](#2-external-user-repository-test) for information about next steps such as how testing occurs for this user environment repository.

### ii. Using local api module vs api module stored in repository (external repo)

Refer to [Using local api module vs api module stored in repository](#ii-using-local-api-module-vs-api-module-stored-in-repository) for the previous use case for more information

### iii. Testing code manually (external repo)

- To test the go modules manually:
  - init.sh in serviceDirectoryName will rerun go build and test.
- To test the whole e2e process manually:
  - Follow the "quick setup steps" mentioned in the README.md at the root of your generated directory
- Refer to section [v. External user repository test.](#v-external-user-repository-test) below to see how you can test more robustly

### iv. How to merge your changes into the repository (external repo)

- Currently the externally generated repositories don't have pull request policies set up, so you can merge directly into master.

### v. External user repository test

This use case does not involve any code generation as code already exists in the user repository. Rather, it is to test the code that exists in the user repository with pipelines (test and dev pipeline).

This testing case is directly related to the above generation case as it tests what was generated and sits in the user repository.

Access example pipelines: <https://dev.azure.com/service-hub-flg/temp_userenv/_build>

What does pipeline look for?

- [pipeline-files/testServiceResourceAndCode.yaml](https://dev.azure.com/service-hub-flg/temp_userenv/_git/temp_userenv_repo?path=/pipeline-files/testServiceResourceAndCode.yaml)

When are pipelines triggered?

- Test pipeline:
  - Triggered upon any merge into master branch.
- Dev pipeline:
  - Only manual trigger
