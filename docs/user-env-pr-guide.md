# Navigating the PR
## Table of contents
 - [Part 1: Changes to testing pipeline](#part-1:-changes-to-testing-pipeline) 
 - [Part 2: User environment](#part-2:-user-environment) 
 - [Part 3: Docker image for installations](#part-3:-docker-image-for-installations) 
 - [Part 4: Miscelaneous changes to repo.](#part-4:-miscelaneous-changes-to-repo) 

## Part 1: Changes to testing pipeline
### 1. Splitting pipeline into seperate jobs.
As a part of this PR, one of the biggest changes was to split the testing pipeline into seperate jobs. The new pipeline yaml file that kickstarts the testing pipeline is: [testServiceResourceAndCode.yaml](../templates/pipeline_templates/testServiceResourceAndCode.yaml) 

Here are the steps the pipeline now takes:
 1. Test Coverage for a service. (<<.serviceInput.directoryName>>TestCoverage) (1st pipeline job in parallel)
    - Job first mentioned [here](../templates/service_templates/deployServicePipeline.yaml) and uses [test suite script](../templates/service_templates/server/test/testSuites.sh) and [test coverage script.](../templates/service_templates/server/test/testCoverageOutput.sh).
    - Since this job does not depend on anything it will run in parallel with GenerateAndPublishEnvConfig as the first job. 
    - This job has 3 tasks, running the test suite check (confirming all folders with go files has a test suite), running the tests to generate the coverage reports, and publishing the coverage reports as an artifact. 
 1. Generate and publishing the environment config. (GenerateAndPublishEnvConfig) (1st pipeline job in parallel)
    - Job first mentioned [here](../templates/pipeline_templates/testServiceResourceAndCode.yaml) 
    - This task creates the environment config using an inline script and publishes it as an artifact such that other jobs can access it.
 2. Provisioning shared resources. (ProvisionSharedResources) (2nd pipeline job)
    - Job first mentioned [here](../templates/pipeline_templates/testServiceResourceAndCode.yaml). Job uses this [yaml file](../resource-provisioning_templates/provisionSharedResourcesPipeline.yaml) as a template. And the yaml file uses this [script.](../resource-provisioning_templates/provisionSharedResources.sh)
    - This job waits to run until GenerateAndPublishEnvConfig is complete.
    - This task first downloads the environment config, then uses a [script](../resource-provisioning_templates/provisionSharedResources.sh) to take care of logging into azure and provisioning the shared resources.

 3. Service deployment related tasks: 
    - Template mentioned [here](../templates/pipeline_templates/testServiceResourceAndCode.yaml) leads to service deployment [yaml](../templates/service_templates/deployServicePipeline.yaml) file
    1. Build and Push Image. (BuildAndPushImage) (3rd/4th pipeline job run in parallel with below)
        - Job first mentioned [here](../templates/service_templates/deployServicePipeline.yaml) and uses this [script.](../templates/service_templates/server/test/buildAndPushImage.sh)
        - This job waits to run until ProvisionSharedResources (#3 mentioned above) is complete, thus can run in parallel with ProvisionServiceResources
        - This task first downloads the environment config, then uses a [script](../templates/service_templates/server/test/buildAndPushImage.sh) to take care of logging into azure before building and pushing image.
        - This job passes in either READPAT or GOPROXY vars depending on if this pipeline is being generated for user's repository or our repository. READPAT allows the docker image to pull the api module from the users repo (dev.azure.com..) while the GOPROXY vars allow the docker image to pull from go.goms.io
    2. Provision service resources. (ProvisionServiceResources) (3rd/4th pipeline job run in parallel with above)
        - Job first mentioned [here](../templates/service_templates/deployServicePipeline.yaml) and uses this [script.](../templates/service_templates/server/test/provisionServiceResources.sh)
        - This job waits to run until ProvisionSharedResources (#3 mentioned above) is complete, thus can run in parallel with BuildAndPushImage
        - This task first downloads the environment config, then uses a [script](../templates/service_templates/server/test/provisionServiceResources.sh) to take care of logging into azure before provisioning service specific resources, and then publishing the created artifacts. (In this case .azuresdk_properties_outputs.yaml that is used in the deploy service task).
    3. Deploying the Service. (DeployService) (5th pipeline job)
        - Job first mentioned [here](../templates/service_templates/deployServicePipeline.yaml) and uses this [script.](../templates/service_templates/server/test/deployService.sh)
        - This job waits to run until BuildAndPushImage and ProvisionServiceResources are complete.
        - This task first downloads the environment config, the downloads the artifact that was published in then uses a [script](../templates/service_templates/server/test/deployService.sh) to take care of logging into azure before deploying the service then checking if it was deployed as expected by checking if [pods are running](../templates/service_templates/server/test/checkServicePodStatus.sh), and if [logs are as expected](../templates/service_templates/server/test/checkServicePodLogs.sh).
 4. Deleting Resources (DeleteResourceGroup). (Final pipeline job)
     - Job first mentioned [here](../templates/pipeline_templates/testServiceResourceAndCode.yaml) and the yaml file uses this [script.](../resource-provisioning_templates/deleteResourceGroup.sh)
    - This job waits to run until DeployService is complete. However, if DeployService never runs, then it will still occur only if ProvisionSharedResources was successful.
    - This task first downloads the environment config, then uses a [script](../resource-provisioning_templates/deleteResourceGroup.sh) to take care of logging into azure and then deleting the resource group that was created in shared resources, (also indirectly eliminating service specific resources since they were generated in the same resource group)

### 2. Using canonical output as testing code.
- Since we are now using seperate jobs as a part of our testing pipeline, the jobs need to be able to access files that already exist in the repository otherwise pipelines cannot be run at all. Thus, we need a version of our most updated generated code to be actually checked into the repository. 
- As canonical output was being checked in already, rather than having two versions of our code checked in, it made most sense to instead run the testing pipeline on canonical output. 
- Changes to canonical output as a result were:
    - It no longer uses go.goms.io as the goModuleNamePrefix, and now uses dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output.
    - PAT was set up in the [yaml file](../testingPipeline.yaml) and when making [testOutput](../testing/testOutput.sh)
    - [testOutput](../testing/testOutput.sh) still acts as it used to before, however, now it prints a warning if there were changes to any of the files that we were previously ignoring (go.sum/go.mod/buf.lock etc) as those files must be updated to the latest version for canonical output testing to be accurate. Development team members MUST pay attention to the warning when running this test and copy over those files if they were changed.

### 3. The artifacts directory.
- We now have an [artifacts directory](../templates/service_templates/server/artifacts/) for the service. Currently there is a .gitignore file under the folder as git doesn't allow us to generate an empty folder, and we want that folder to be accessible in the pipeline not generated in each job.
- Currently .azuresdk_properties_outputs.yaml gets produced in the artifact directory when provisioning service specific resources gets called.

## Part 2: User environment
### 1. Terraform files
- The [terraform files](../user-environment/terraform-files/) are where the files required for terraform exist..
    - [main.tf](../user-environment/terraform-files/main.tf) define the resources that are created as a requirement for the user environment including ADO project, repo, pipeline, and when service principal comes back online the aad application, service principal, service principal password, and role assignment.
    - [variables.tf](../user-environment/terraform-files/variables.tf) define the variables that are used for the terraform resources. (More information about the variables can be found in the [readme](README.md))
    - [output.tf](../user-environment/terraform-files/output.tf) defines the outputs that we need to use for settings such as git configuration, goModuleNamePrefix, etc.
- When you apply terraform files, terraform.tfstate and terraform.tfstate.backup are created but they are under gitignore as they hold credentials/secrets. NEVER check these files into the repo.

### 2. Reworking goModuleNamePrefix
- For the user environment to function properly we want to use the service's api module within the repository we are creating. That means everywhere we would importing from go.goms.io for our testing/development, the user's version should be using their own api module. Thus, goModuleNamePrefix is what changes.
- For the user environment the goModuleNamePrefix is extracted from the terraform output, and then passed into generation with the [env-information.yaml](../config-files/user-environment/env-information.yaml) config file. (Side Note: Since there is no terraform extraction for service hub repo testing, we just pre-set go.goms.io in our [env-information.yaml](../config-files/testing/env-information.yaml))
- The extraction for the user-environment takes place in the script mentioned in [4. Creation of user environment](#4-creation-of-user-environment) .

### 3. Changes to repo to adjust to downloading from personal and private repository (goModuleNamePrefix)
- In order for go mod tidy and other go commands to be able to download the api module correctly when building and pushing image the following changes were necessary
    1. Conditionally template in [dockerfile](../templates/service_templates/server/Dockerfile) to include PAT as env var and run git configuration to be able to access private repository or include GOPROXY vars for accessing go.goms.io.
    2. Conditionally template in [template-Makefile](../templates/service_templates/server/generated/template-Makefile) to build the docker image with READPAT or GOPROXY vars depending on if its user environment or our testing environment
    3. Conditionally template in [deployServicePipeline.yaml](../templates/service_templates/deployServicePipeline.yaml) to include READPAT or GOPROXY vars as environment variables into the scripts depending on if its user environment or our testing environment
    4. Conditionally template in [api/v1 Makefile](../templates/service_templates/api/v1/Makefile) specifically when running the bash command for container. If we are in user environment we must add GOPRIVATE and git configuration for docker container to work properly.

### 4. Creation of user environment
- To create the user environment the following steps are required. They are done by the [createInitialUserEnv script](../user-environment/createInitialUserEnv.sh) :

    1. Exports environment variables (project/read PATs, service principal values) such that they can be used by terraform.
    2. Extracts pipelinePath and serviceDirectoryName from config files
    3. Applies terraform files state
    4. From terraform output, extract values use for the goModuleNamePrefix and for gitConfiguration. Without git configuration, we will not be able to download api module when building server in init.sh
    5. Generate code. (It is generated to user-environment/generated-output, but we do not check it in).
    6. Copy terraform files over to generated code folder, and [remove secret vars](../user-environment/removeSecretVars.sh) from terraform files. 
    7. Add git remote origin to be the url we created from terraform outputs.
    8. Run [init.sh](../templates/service_templates/init.sh). When init.sh is run, the module is tagged and pushed if we are doing it for user environment as it was conditionally templated in.
    9. Wait 30 seconds to make sure the terraform created pipeline is settled before pushing to the repository. If we push too soon the pipeline isn't automatically triggered as it takes a few mins for pipeline to settle after it was created. 
    10. Push to repo.


### 5. The user's testing pipeline (exactly the same as our testing pipeline) refer to:  - [1. Splitting pipeline into seperate jobs.](#1.-splitting-pipeline-into-seperate-jobs) 

## Part 3: Docker image for installations
- An [image](../go-templating/installations/Dockerfile) was created to help users and developers avoid having to download external packages. It installs the following
    - The right package of golang depending if arm or amd.
    - wget, default-jre, pip, npm
    - go packages for buf, mock, grpc, protobuf, grpc-ecosystem.
    - swagger-nodegen-cli
- It also sets up aks goproxy for when required.
- It is used in [api/v1 Makefile](../templates/service_templates/api/v1/Makefile) as a wrapper around the commands that used to be make service.
- It also works for successfully running init.sh upon generating such that testing no longer requires installations (both locally and on pipeline).
- As a result generateCode can be entirely eliminated since we no longer require the installations for pipeline.
## Part 4: Miscelaneous changes to repo.
### 1. Addition of <<.sharedInput.directoryPath>> to all config files
- directoryPath's purpose is a little bit difficult to understand. It exists to help create the correct paths for pipelines such that we do not need two sets of pipeline files depending on if we are testing or if its user's pipeline.
- For example:
    - directoryPath in our [config files](../config-files/testing/service-config.yaml) for our testing is testing/canonical-output/ while for user environment [config files](../config-files/user-environment/service-config.yaml) it is "" (blank).
- This is because pipeline files need paths to be already set and relative to where the pipeline yaml file sits. In our testing pipeline [yaml file](../testing/canonical-output/pipeline-files/testServiceResourceAndCode.yaml) you can see that targetPaths and filePaths start with testing/canonical-output (i.e. testing/canonical-output/shared-resources/deleteResourceGroup.sh) as relative to our repository, that is where the files sit (the generated code repository we are testing sits in canonical output). While relative to the user's repository, the files sit at the root of the repository thus the directoryPath is blank (i.e. shared-resources/deleteResourceGroup.sh)
### 2. Changes to goModuleNamePrefix for generation and files to adjust to this change
- refer to [2. Reworking goModuleNamePrefix](#2-reworking-gomodulenameprefix) and [3. Changes to repo to adjust to downloading from personal and private repository (goModuleNamePrefix)](#3-changes-to-repo-to-adjust-to-downloading-from-personal-and-private-repository-gomodulenameprefix)
### 3. Making init.sh more robust to failures. 
- Because [init.sh](../templates/service_templates/init.sh) is being placed in a couple scripts, there have been instances where if commands in init.sh fail its hard to notice as things continue on as normal. I have added some exits to help catch when commands are failing such that we catch those errors instead of spending extra time debugging why things weren't working at a latter stage. 
