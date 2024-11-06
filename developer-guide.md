# Developer Guide

## Using our PAT Rotator to avoid manually changing PATs

Our PAT rotator will be rotating pats into the azure key vault "servicehubkv" automatically on a weekly basis.
By following the below instructions, rather than having to manually change READPAT, everytime you refresh your terminal/bashrc,
it will grab the most recent PAT and add to the git config file.

1. Make sure you have at least role of "Key Vault Secrets User" for the subscription "AKS Long Running Things". Check here: <https://ms.portal.azure.com/#@microsoft.onmicrosoft.com/resource/subscriptions/359833f5-8592-40b6-8175-edc664e2196a/users>
2. Add the following lines to your bashrc file

    ```bash
    export READPAT=$(az keyvault secret show --name READPAT --vault-name servicehubkv --query value -o tsv)
    if [ -n "$READPAT" ]; then
        # Look for any url lines in the gitconfig file that match "service_hub" as that means there is an older
        # configuration already set up that we must remove.
        git config --global --get-regexp 'url.*service_hub' | while read line; do
          key=$(echo "$line" | cut -d ' ' -f 1) # Extract the key from the line
          git config --global --unset "$key" # Unset the configuration
        done
        git config --global url."https://$READPAT@dev.azure.com/service-hub-flg/service_hub/_git/service_hub".insteadOf "https://dev.azure.com/service-hub-flg/service_hub/_git/service_hub"
    fi
    ```

3. If you are developing on temp_userenv also add the following lines.

    ```bash
    if [ -n "$READPAT" ]; then
        # Look for any url lines in the gitconfig file that match "temp_userenv" as that means there is an older
        # configuration already set up that we must remove.
        git config --global --get-regexp 'url.*temp_userenv' | while read line; do
          key=$(echo "$line" | cut -d ' ' -f 1) # Extract the key from the line
          git config --global --unset "$key" # Unset the configuration
        done
        git config --global url."https://$READPAT@dev.azure.com/service-hub-flg/temp_userenv/_git/temp_userenv_repo".insteadOf "https://dev.azure.com/service-hub-flg/temp_userenv/_git/temp_userenv_repo"
    fi
    ```

4. Refresh your source file or restart terminal for changes to take effect.

    ```bash
    source ~/.bashrc
    ```

## Simple End to End flow example

1. If you want to develop with generated code instead of templates
    - Follow steps in [Pre-requisites and Quick Start to generating a service section of README.md](README.md)
2. If you are developing in ~/projects (skip if changes were made directly in templates)
    - Make changes to files in ~/projects, follow steps mentioned in the corresponding folder's readmes and test as required within the folders.
    - Convert the changes into the _templates folders and make sure to manually convert instance value to its template variable any variables that were used.
      - i.e. mygreeterv3 => <<.serviceInput.directoryName>>  
2. If any new files were added to template folders or file names were changed

    ```bash
    #To use generate all template specs
    make genAllTemplateSpec
    ```

    Make sure to check the template_spec.csv files to confirm they look as expected.
3. If you made changes to service directory's api module
    - Refer to [updating api module](#updating-api-module-to-service-hub)
    - Make sure to wait sometime for tag to settle.
4. If api module tag is completely up to date and settled. Generate canonical output

    ```bash
    make testOutput
    ```

    Take a look at the diff check to determine if everything changed was as expected.
5. Update the canonical-output directory.

    ```bash
    rm -rf testing/canonical-output/*; cp -Rp testing/generated-output/* testing/canonical-output/
    ```

6. Push to service hub branch
7. Generate a Pull Request, add your chosen service hub memeber as a reviewer.
8. In order for pull request to be approved for merging, your reviewer must approve and you must queue and successfully pass the two check pipelines at the top of the pull request.

Optional steps:

- Additional test commands

  ```bash
  #To test coverage tests
  make testCoverage
  # To run all tests
  make testAll
  ```

- Check if any file differences between test output and canonical

  ```bash
  diff -r testing/generated-output testing/canonical-output
  ```

## Further details on development steps

### Generate template spec file

- The template spec file (.templateSpec.csv) should exist in the templates folder that you are using to generate the required files. It stores File name, Overwrite capabilities, and a user list.
  - You can either manually create/edit the file however a command for automatic generation exists as defined below.
  - Template spec file MUST be in lexicographical order for generation to work as expected. (Directories first, then files)
  - If user is provided when generation, it will only generate files that have the user in its user list. Current default is "aks" in every files user list.
  - Changes to overwrite capabilities and user list must be done manually. Regeneration can be done if you want files to be automatically added in lexicographic order to template spec. The command will maintain the settings of the file in the previously existing spec.
- To automatically generate template specs

```bash
#To use generate all template specs
make genAllTemplateSpec
#-------------
#To use generate service template spec
make genServiceTemplateSpec
#-------------
#To use generate resources template spec
make genResourcesTemplateSpec
#-------------
```

## Switching between AKS development and external user development

Currently aks development exists under ~/projects/aks-rp, while external user development is under ~/projects/external.

When you are switching back and forth, the main error that will arise is pulling the golang modules correctly.

This set up should always exist:

- Make sure your GOPRIVATE includes "dev.azure.com" (in addition to the ones used for AKS), for example.

  ```bash
  export GOPRIVATE="goms.io/aks/*,go.goms.io/aks/*,go.goms.io/fleet*,dev.azure.com"
  ```

- Git config is set up as mentioned in README.md

    ```bash
    # Check if you already have it set up
    vi ~/.gitconfig
    #If not already there, run the following
    git config --global url."https://$READPAT@dev.azure.com/service-hub-flg/service_hub/_git/service_hub".insteadOf "https://dev.azure.com/service-hub-flg/service_hub/_git/service_hub"
    ```

From AKS to external:

- You must unset GOPROXY variable as what it does is rerout any sites under GOPRIVATE to instead use the authentication method defined in the GOPROXY variable. If you do not unset this variable for external generation, when it tries to download the api module, it will try to reroute authentication through the GOPROXY variable.

```bash
unset GOPROXY
```

From external to AKS :

```bash
source ~/.bashrc
```

## Updating API Module to service hub

Development for service hub will use the following goModuleNamePrefix

```bash
goModuleNamePrefix: dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output
```

This is because, the api module that we have checked in at the most updated form will exist under testing/canonical-output/mygreeterv3/api.

As a result, when you tag the api module, you must do the following

```bash
#Make sure canonical output is most updated version
make testOutput
rm -rf testing/canonical-output/*; cp -Rp testing/generated-output/* testing/canonical-output/
# Tag and push api module.
cd testing/canonical-output
cd mygreeterv3
git add ./api
git commit -m "api module v$(versionNumber)"
git tag testing/canonical-output/mygreeterv3/api/v$(versionNumber)
git push origin testing/canonical-output/mygreeterv3/api/v$(versionNumber)
```

If you want service hub to use your updated tag version of api module, you must change the version of api that is being required in [init.sh](templates/service_templates/init.sh) in the line that says the following, make sure to update the version number to match yours. If you want to test using the pipeline, you must make sure you change this step otherwise the pipeline will simply use the last stable version of the api module.

```bash
go mod edit -require <<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>@v0.0.2
```

## Updating middleware module to reflect in service hub

Similar to updating api module above, we are using the last stable version for the middleware module.

- If you want service hub to use your updated version of middleware module, you must change the version of middleware that is being required in [init.sh](templates/service_templates/init.sh) in the line that says the following, make sure to update the version number to match yours. If you want to test using the pipeline, you must make sure you change this step otherwise the pipeline will simply use the last stable version of the middleware module.

  ```bash
  go mod edit -require github.com/Azure/aks-middleware@v0.0.19
  ```

## Resource Provisioning

### Publishing Bicep Modules

In our directory `bicep-modules`, there is a list of directories, each representing a resource. Each directory contains a bicep module that helps create the resource. If you make changes to a bicep module, make sure you publish the module.

Example:
`make publishModule DIRECTORY_NAME=aks-managed-cluster`

## Adding a new set of service templates

### Testing normal deployment

Before testing Ev2 deployment it is important to make sure normal deployment is working for your service.

Run the **Resource and service deployment** pipeline in service_hub, and make sure your service is deploying correctly in the logs.

If your service is not showing up in the pipeline, make sure to set runPipeline: true in your service-config.yaml file.

### Testing Ev2 Deployment

#### Generate into temp_userenv_repo

1. Clone [temp_userenv_repo](https://dev.azure.com/service-hub-flg/temp_userenv/_git/temp_userenv_repo)
2. Create a new branch in the repository
3. Change **destinationDirPrefix** in [user-environment/common-config.yaml](config-files/user-environment/common-config.yaml) to match the path where you cloned the repository.
4. Generate by running the following commands in service_hub

    ```bash
    #Set config files
    genConfig="../config-files/generator-config.yaml"
    servConfigFolder="../config-files/user-environment/service-configs"
    commonConfig="../config-files/user-environment/common-config.yaml"
    terraformOut="../config-files/user-environment/env-information.yaml"
    #Generate pipeline files, service and shared resources
    make generateAll serviceConfig=$servConfigFolder commonConfig=$commonConfig generatorConfig=$genConfig envInformation=$terraformOut
    ```

#### Initialize your services

1. For every service, run ./init.sh.
2. If it is a new service or there have been changes to the api module to an older service, you will need to tag the new api module for it.

    ```bash
    # Inside temp_userenv_repo.
    cd **serviceDirectoryName**
    git add ./api
    git commit -m "api module v$(versionNumber)"
    git tag **serviceDirectoryName**/api/v$(versionNumber)
    git push origin **serviceDirectoryName**/api/v$(versionNumber)
    ```

#### Create build and release pipeline

**NOTE**: Builds will work on your own branch, however releases will not. Once you are sure the build pipeline was successful and the artifacts look as expected, you will need to create a PR to merge into master. Then re-run the build on master, and run the release when ready.

Follow the instructions in the **Build and Release with EV2 through OneBranch Pipelines** section of README/Ev2_README.md

## New Service

If you need to create a new service from scratch that fulfills one of your product requirements, we have created the basic_templates to help you with this process. The basic_templates are a very minimal service without any added features (no Async or azure sdk, only restsdk and the basic client/server) that you can use as a baseline for your service. Simply follow these steps to begin creating a new type of service:
1. Duplicate basic_template directory and name it with your new service name (e.g. authentication_templates).
2. In `config-files/generator-config.yaml` add the path to your new directory following the format in that same file, providing a name and path.
3. In `config-files/local-development` under both the `external-generation` and `aks-generation` duplicate the `basic-config.yaml` file and modify it's name and contents to match the name of your new service, making reference to the name of the path to your new templates done in step 2. Repeat this process for the `testing` and `user-environment` folders.
4. Run `make genService`, to the location defined as `destinationDirPrefix` in your `common-config` file.
5. Make changes to the templates to modify the generated code.
