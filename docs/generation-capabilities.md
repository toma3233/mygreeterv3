# Service Hub Generation Capabilities

## Generation flags

| Flag   |  Required |  Default | Definition | Data accessible as |
|---|---|---|---|---|
| generatorConfig | For all component generation |../config-files/generator-config.yaml|  Stores information about what the template paths for the generation are. | GeneratorInputs map[string]string
| commonConfig | For all component generation |../config-files/local-development/external-generation/common-config.yaml| Used for shared variables between all modules. | SharedInput sharedInput `yaml:"sharedInput"`ResourceInput resourceInput `yaml:"resourceInput"` PipelineInput pipelineInput `yaml:"pipelineInput"`
| envInformation | For all component generation |../config-files/local-development/external-generation/env-information.yaml| Stores environment information that is outputted from terraform generation. | EnvInformation envInformation `yaml:"envInformation"`
| serviceConfig | Only for "pipeline" or "service" component generation |../config-files/local-development/external-generation/service-configs| The path folder that stores the config files for all the services you want to generate or the single file path if you are only generating one. | ServiceInput serviceInput  `yaml:"serviceInput"`
| ignoreGarbageDeletion |For all component generation| False |If state file does not exist, this ignores garbage deletion of previously generated files only if flag is set to true |---
| user |For all component generation| external |Which user is generating (must match a user in templateSpec.csv) | User
| generationType |For all component generation| resource |Which type of module we are generating. Must be one of "service", "resource", or "pipeline" |---
| extraConfig | Optional |""| Used to store extra inputs not defined in the schema of other config files. | ExtraInputs map[string]string

Refer to [generationInterface](../go-templating/generationInterface.go) for the schema and examples of the variables.

### External Users

#### Shared Resources

- To generate:

    ```bash
    #To use default config files
    make genResources
    #-------------
    #If you want to specify what flags, command is as following (with default info)
    make genResources commonConfig=../config-files/local-development/external-generation/common-config.yaml generatorConfig=../config-files/generator-config.yaml envInformation=../config-files/local-development/external-generation/env-information.yaml ignoreGarbageDeletion=false
    ```

#### Service (any general module)

- If the path is a folder, **ALL** the files under ***serviceConfig*** will be taken into account as separate services, and generation will occur for all of them.
- If you only want to generate one service for example, you would pass in path to that one service config.
- To generate:

    ```bash
    #To use default config files
    make genService
    #-------------
    #If you want to specify what flags, command is as following (with default info)
    make genService serviceConfig=../config-files/local-development/external-generation/service-configs commonConfig=../config-files/local-development/external-generation/common-config.yaml generatorConfig=../config-files/generator-config.yaml envInformation=../config-files/local-development/external-generation/env-information.yaml ignoreGarbageDeletion=false
    ```

#### Pipeline

- In the default case, generating the pipeline folder will also take into account **ALL** the files under ***serviceConfig*** (if the provided path is a folder) as separate services and add each service that has **"runPipeline"** set to true to the pipeline file that runs the overarching service hub deployment flow.
- You can pass an empty string "" into the serviceConfig flag if you don't want the test/dev pipeline to deploy the seperate services.
- To generate:

    ```bash
    #To use default config files
    make genPipeline
    #-------------
    #If you want to specify what flags, command is as following (with default info)
    make genPipeline serviceConfig=../config-files/local-development/external-generation/service-configs commonConfig=../config-files/local-development/external-generation/common-config.yaml generatorConfig=../config-files/generator-config.yaml envInformation=../config-files/local-development/external-generation/env-information.yaml ignoreGarbageDeletion=false
    ```

#### Generate all modules

- To generate:

    ```bash
    #To use default config files
    make generateAll
    #-------------
    #If you want to specify what flags, command is as following (with default info)
    make generateAll serviceConfig=../config-files/local-development/external-generation/service-configs commonConfig=../config-files/local-development/external-generation/common-config.yaml generatorConfig=../config-files/generator-config.yaml envInformation=../config-files/local-development/external-generation/env-information.yaml ignoreGarbageDeletion=false
    ```

### AKS Users

#### Service

- If the path is a folder, **ALL** the files under ***serviceConfig*** will be taken into account as separate services, and generation will occur for all of them.
- If you only want to generate one service for example, you would pass in path to that one service config.
- To generate:
  - Modify [aks mygreeterv3-config.yaml](../config-files/local-development/aks-generation/service-configs/mygreeterv3-config.yaml) vars to customize the service

    ```bash
    #To use default config files
    make genAKSService
    #-------------
    #If you want to specify what flags, command is as following (with default info)
    make genService serviceConfig=../config-files/local-development/aks-generation/service-configs commonConfig=../config-files/local-development/aks-generation/common-config.yaml generatorConfig=../config-files/generator-config.yaml envInformation=../config-files/local-development/aks-generation/env-information.yaml ignoreGarbageDeletion=false user=aks
    ```

## How input is used for code generation

### Variables

- Inputs passed into template execution

    ```bash
    allInput allInput 
    ```

- Variables defined as structs vs maps
  - The below structs are used to store the information that is "required" or well defined for the purpose of generation.

    ```bash
    SharedInput sharedInput `yaml:"sharedInput"`
    ResourceInput resourceInput `yaml:"resourceInput"`
    PipelineInput pipelineInput `yaml:"pipelineInput"`
    EnvInformation envInformation `yaml:"envInformation"`
    ServiceInput serviceInput  `yaml:"serviceInput"`
    ```

  - The below maps are used to store the variables that can be added to or expanded in order to provide the user additional flexibility.

    ```bash
    # Multiple services
    ServiceNameToRunPipeline map[string]bool `name:"serviceNameToRunPipeline"`
    # Additional inputs not previously defined
    ExtraInputs map[string]string `name:"extraInputs"`
    # Additional template options
    GeneratorInputs map[string]string `name:"generatorInputs"`
    ```

- Variables shared amongst components:
  - Refer to [generationInterface](../go-templating/generationInterface.go) for detailed information about the breakdown to each struct and input variable.
- Variables used for templating code:
  - allInput.GeneratorInputs in combination with generationType and the component's templateName is used to get the **directory of the template** for template execution during generation.
    - allInput.ResourceInput.templateName utilised to get shared-resources template.
    - allInput.ServiceInput.templateName utilised to get the module template.
    - allInput.PipelineInput.templateName utilised to get the pipelines template.
  - generationType in combination with the component's directoryName is used to used to extract the  **destination directory** for template execution during generation.
    - allInput.ResourceInput.DirectoryName utilised to generate shared-resources directory.
    - allInput.ServiceInput.DirectoryName utilised to generate the module directory.
    - allInput.PipelineInput.DirectoryName utilised to generate pipelines directory.
  - allInput struct gets converted into a map to maintain consistency with the variables defined in the yaml file. generatorInputs will **NOT** be accessible to template vars as we do not want to give access to this information.

### Generation strategies and further details

- General modules/components are considered to fall under the "service" category for generation.
- If a path to a folder is passed in for serviceConfig folder, all the config files in the folder will be used for "service" generation one by one. This is to allow for multiple service/component generation at once.
- The "pipeline" generation utilizes the service information to decide whether or not to run the service's pipeline as a part of the all-encompassing pipeline that also creates the shared resources. This is to allow ease of adding service information to pipelines if you have created a new service.
