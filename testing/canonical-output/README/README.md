# Setting up Service Hub Generated Code

## Table of Contents

- [Prerequisites](#prerequisites)
  - [Operating System](#operating-system)
  - [Installations](#installations)
  - [Personal Auth](#personal-authentication)
  - [Generated Directory Check](#generated-directories)
- [Terminology](#terminology)
- [Setting up](#creating-your-setup)
  - [Quick setup - Manual](#quick-setup-manually)
  - [Quick setup - Development pipeline](#quick-setup-using-development-pipeline)
- [Deleting resources](#deleting-your-resources)

## Warning

Do not run any commands with sudo.

## Prerequisites

### Operating system

Service hub functions on Linux.

- If you are on Windows, follow the steps to install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) if you do not already have it

### Installations

For all of these installations, we recommend checking if they already exist on your WSL set up before installing.

- Follow the steps to install [Go](https://go.dev/doc/install). Make sure to add Go to your path.

- Follow the steps to install [Docker](https://docs.docker.com/engine/install/).  Also turn on [Docker Desktop WSL 2](https://docs.docker.com/desktop/wsl/) such that docker works with WSL.

- Install [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli).
  - Azure Subscription Requirements

        In order to create all the resources (shared and service specific), you must have role permissions in your subscription to:

        1. Create, update and delete resources
        2. Assign roles

        Note that `Contributor` access does not have the permissions to assign roles. `Owner` access has all the permissions necessary.
  - After installing, make sure to log in and set your subscription

        ```bash
        az login
        az account set --subscription $subscriptionId
        ```

- Check Bicep version to see if it was installed after Azure CLI is installed

    ```bash
    az bicep version
    # to upgrade
    az bicep upgrade
    ```

- You can choose to install [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) if you would like, or you can use our pre-built docker container for a clean environment that has everything installed once deploying your service.

- We assume that Bash is used as the default shell. If you are using Zsh, please make the necessary adjustments throughout the instructions accordingly.

### Personal Authentication

If you have not already done the following steps.

Set up access for go build to be able to download modules from your organization for generation.

Creating a Personal Access Token (PAT) is essential to obtain the necessary access. There are two methods to achieve this. The first method involves setting up a manual PAT, which expires in a week and requires manual renewal each time it expires. The preferred method is to use a rotator that automatically fetches and injects a new token for you. Detailed instructions can be found [here](https://dev.azure.com/service-hub-flg/_git/service_hub?path=/developer-guide.md&version=GBmaster&_a=preview). If you choose to create the PAT manually, please follow the steps below.

- Create a personal authentication token under your organization to gain access to modules stored in your repository.
  - At your org's ADO(AzureDevOps e.g. [here](https://service-hub-flg.visualstudio.com/)) page, click the person with settings button next to your initials at the top right.
  - Click on Personal Access Tokens
- Create a new token for org account under "Code" scope with "Read" access. You can specify an expiration date. Make sure Organization is selected to be service-hub-flg, cross organization PAT won't work.
- Save token as environment variable READPAT (either by exporting in terminal or adding to bashrc profile).
- Also do the same for the GOPRIVATE variable below.
- Preferably add the lines to bashrc profile so that you do not have to set the variables again and again everytime you open a new terminal.

    ```bash
    export READPAT=[Your Personal Access Token of Azure Devops of your org]
    export GOPRIVATE="dev.azure.com"
    ```

- If you have the GOPROXY variable set for AKS or other uses, you must do one of the below.

    ```bash
    #you must either comment it out of your bashrc and refresh by running
    source ~/.bashrc
    #or you can unset the GOPROXY variable in your terminal for now. If you only unset it, you must do it every time your
    # bashrc refreshes, you use a new terminal.
    unset GOPROXY
    ```

- Configure git to access repo with PAT instead. This will automatically authenticate downloading from this repository.

    ```bash
    git config --global url."https://$READPAT@dev.azure.com/service-hub-flg/service_hub/_git/service_hub".insteadOf "https://dev.azure.com/service-hub-flg/service_hub/_git/service_hub"
    ```

- By setting up GOPRIVATE and the git configuration as mentioned above, golang will be able to access modules stored in private azure devops repo. For more information regarding this refer to the [Azure Devops Git docs](https://learn.microsoft.com/en-us/azure/devops/repos/git/go-get?view=azure-devops)

### Generated Directories

We assume you have generated the below directories:

- **shared-resources**
- **serviceDirectoryName**
- **pipeline-files**

If you have not generated these directories, please return to the `service_hub` repo.

## Terminology

| Term   |  Examples | Definition |
|---|---|---|
| shared resources |---|  Azure resources shared and used by multiple services |
| service resources |---| Azure resources that are just for a single service |
| generated directory |---|  The directory that contains all the required directories (and the one that currently holds this README.md file) |
| service directory |mygreeterv3| The directory within the generated directory that contains the generated service  |
| shared resource directory |shared-resources| The directory within the generated directory that contains code for generating shared resources |

# Creating your setup

You can choose to either perform your set up manually or through the development pipeline as seen in the next two sections.

## Quick Setup (manually)

The following shows you the steps on how to get a working setup quickly.

If you want to understand why we are doing each of the steps, or require more detail, please look at the corresponding README's in the [shared-resources folder](shared-resources/README.md) and the serviceDirectoryName folder -> serviceDirectoryName/README.md

```bash
# --------------------------------
# Generate environment config (env-config.yaml)
# Both shared-resources and serviceDirectoryName will use the file, so it will sit
# in the generated directory.
# The following makefile target assumes you have an environment variable called $subscriptionId.
# Either set the environment variable before running the command, or copy and paste your subscription
# id into env-config.yaml
# Make sure the subscriptionId matches the id you used to log into Azure CLI.
# Make sure resourcesName value is present. Otherwise, set a value using only alpha numberic characters.
cd README
make genEnvConfig
# --------------------------------
# Enter shared resources directory
cd shared-resources
# Templates env-config.yaml values into all the required files. We assume env-config.yaml exists in your generated folder. (i.e. the folder that stores the generated directories)
make template-files
# Create shared resources
make deploy-resources # takes about 10 minutes
cd ..
# --------------------------------
# Enter service directory to perform required steps
cd serviceDirectoryName
# In svc directory initialize.
./init.sh
cd server
# Templates env-config.yaml values into all the required files. We assume env-config.yaml exists in your generated folder. (i.e. the folder that stores the generated directories)
make template-files
# Build image (can be done in parallel with deploy-resources). If you are using macOS, building a Docker image locally won't work on an AKS cluster because we currently don't support cross-platform build. Instead, you need to use DevBox to build the Docker image. The following commands should be executed on DevBox. If you do not have a DevBox set up, please follow the instructions in https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/CloudNativeCompute.wiki/358303/Dev-Box
make build-image
# Create service specific resources
make deploy-resources # takes 5-10 minutes
# Push image to acr
make push-image
# (If svc running on aks cluster) Upgrade service on AKS cluster
make upgrade
# (If svc not running on aks cluster) Deploy service to AKS cluster
make install
# Refresh your bashrc such that your terminal has updated kubeconfig
source ~/.bashrc
```
<!-- TODO: Support cross-platform image build -->

### Refer to the following section in the service directory readme to check if your service is deployed correctly

Check if Service Deployment Successful -> found in serviceDirectoryName/README.md

## Quick Setup (using development pipeline)

### Start the dev pipeline

(TODO: Link dev pipeline)

Go into pipelines, select the development pipeline and manually trigger the pipeline by pressing "Run Pipeline". Choose the branch you would like to create your set up with.

[Optional] If you would like to change the subscription id that the dev pipeline uses.

- Under "Advanced options" select variables.
  - Change SUBSCRIPTION_ID to your preferred id.

Once the dev pipeline has finished running, quick set up will be different as a lot of the commands have already been done by the dev pipeline.

### Finding your resources

The development pipeline will display warnings that will point you to the resource group that your resources were created in.

Pressing on the warning will direct you to a link to the azure portal for that resource group.

### Copy over required files

From the pipeline run, download the following files and move them to their appropriate locations.

- env-config.yaml : move to root of generated directory.
- .Azuresdk_properties_outputs.yaml : move to serviceDirectoryName/server/artifacts

### Make sure you are logged into azure

```bash
# Log into azure 
az login
# Set you subscription to the same subscription id that was used for the dev pipeline
az account set --subscription $subscriptionId
```

### Template all files

```bash
# Template shared resources files
cd ../shared-resources
#Templates env-config.yaml values into all the required files. We assume env-config.yaml exists in your generated folder. (i.e. the folder that stores the generated directories)
make template-files
cd ..
# Template service files
cd serviceDirectoryName/server
#Templates env-config.yaml values into all the required files. We assume env-config.yaml exists in your generated folder. (i.e. the folder that stores the generated directories)
make template-files
```

### Connect to your aks cluster

```bash
cd serviceDirectoryName/server/generated
# Gets credentials for a managed kubernetes cluster
make connect cluster
#Refresh your bashrc such that your terminal has updated kubeconfig
source ~/.bashrc
```

### Refer to the following section in the service directory readme to check if your service is deployed correctly

Check if Service Deployment Successful -> found in serviceDirectoryName/README.md

## Deleting your resources

### Through the development pipeline

1. Go into pipelines, select the development pipeline and manually trigger the pipeline by pressing "Run Pipeline".

2. Under "Advanced options"
    - Select variables.
        - Set RESOURCES_NAME to be the unique id that belongs to your resources. (If you previously created resources through dev pipeline, it would have been outputed as a warning, or if you did it manually, it is resourcesName variable in your env-config.yaml file).
        - Set DELETE to be true.
    - Select stages to run. ONLY select the "Delete all resources" stage. Ignore the warning about the skipped stage.
3. Press run pipeline. This pipeline will delete your resources if the service principal pre-set for the pipeline has access to the subscription the resources were created in. (If you used a different subscription, refer to [manual deletion](#manual-deletion))

### Manual deletion

Run the following command, where $resourcesName is your unique id (found in env-config.yaml)

```bash
az group delete -n servicehub-$resourcesName-rg --yes
```

## Making changes to Bicep Resources

### Modifying Pre-defined Bicep Resources

If you need to adjust the parameter values for your resources, follow these steps:

1. Check the [Bicep Module Registry](https://github.com/Azure/bicep-registry-modules/tree/main) for available parameter values that you can modify. Look for parameters in the corresponding Bicep files, indicated by the syntax param. These directions do not apply for `aks-managed-cluster` and `subscription-role-assignment` bicep resource. See *note* below.
2. Modify the parameter values in the relevant Bicep files according to your requirements.
3. After making the necessary changes, execute the command `make deploy-resources` to apply the updated configurations.

*Note: For AKS managed cluster or subscription role assignment, refer to the bicep-modules directory in the Service Hub registry. Locate the corresponding directory and review the service-hub-generated-module.bicep for available parameters.*

### Adding Bicep Resources

To explore available Bicep modules and their definitions, check out the following links:

- [AVM Bicep Resource Modules](https://azure.github.io/Azure-Verified-Modules/indexes/bicep/bicep-resource-modules/): This resource provides a list of Bicep modules along with their locations in the Bicep Module Registry.
- [Bicep Module Registry](https://github.com/Azure/bicep-registry-modules/tree/main): The Bicep Module Registry contains modules used by the AVM Bicep Resource Modules.
Keep in mind that we’ve already set up resources for you using the AVM Bicep Resource Modules. If you wish to modify or add resources, refer to the Bicep resource modules and registry for guidance. Detailed examples are provided in their README.md files.
- [Azure Container Registry Example](https://github.com/Azure/bicep-registry-modules/tree/main/avm/res/container-registry/registry). This links to the azure container registry directory in the bicep module registry. In the README.md, there is documentation and examples of how to add the azure container registry. Make sure to specify the version.

Keep in mind that we’ve already set up resources for you using the AVM bicep resource modules. If you want to modify or add resources, refer to the bicep resource modules and registry to learn how to.

### "Define Once, Reference Everywhere Else" Rule

We introduce the concept of defining a resource once and only referencing it in other locations. An example in our repo is the log analytics workspace. In `resources/Main.SharedResources.Template.bicep` from the `shared-resources` directory, we define the log analytics workspace:

`resources/Main.SharedResources.Template.bicep` (source file of workspace):

```bicep
// Defined here
module workspace 'br/public:avm/res/operational-insights/workspace:0.3.4' = {
  name: 'servicehub-${resourcesName}-workspaceDeploy'
  scope: resourceGroup(subscriptionId, resourceGroupName)
  params: {
    name: 'servicehub-${resourcesName}-workspace'
    location: rg.outputs.location
  }
}
```

`Azuresdk.ServiceResources.Template.bicep` in service directory

```bicep
// referenced here
resource logAnalyticsWorkspace 'Microsoft.OperationalInsights/workspaces@2022-10-01' existing = {
  name: 'servicehub-${resourcesName}-workspace'
  scope: resourceGroup(subscriptionId, resourceGroupName)
}
```

Note that we use the `resource` and `existing` syntax to reference the resource. Both help us prevent the modification of the resource from it's non-source file.
