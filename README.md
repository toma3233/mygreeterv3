# Service Hub

This repo holds the code that helps engineers to bootstrap a service. It takes care of the following tasks for engineers:

- Development environment: ADO project, repo, pipeline, service connection, etc.
- Service skeleton: api, service implementation, middleware, test, build, release, etc.
- Resource provisioning: shared resource, service specific resource declaration, build, and release.
- Logs generation, collection, dashboard, and alerting.

For more details, see [Service Hub wiki](https://dev.azure.com/service-hub-flg/service_hub/_wiki/wikis/service_hub.wiki/4/About-Service-Hub).

The aks-middleware module that we use is in a public [github repo](https://github.com/Azure/aks-middleware).

# User Guide

This file provides the concrete steps on how to use Service Hub. See [high level user guide](https://dev.azure.com/service-hub-flg/service_hub/_wiki/wikis/service_hub.wiki/96/User-Guide) to understand the big picture.

**Developers**: If you are developing features for this repo, please also refer to the [Developer Guide](./developer-guide.md).

# Pre-requisites

## Operating system

Service hub runs on Linux and MacOS.

- If you are on windows, follow the steps to install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install).

## Installations

- Follow the steps to install [Go](https://go.dev/doc/install).

# Creating a new MSFT service

Feel free to skip the steps if you have already performed them or have an existing service.

The steps for service tree entry, security groups, Jit policy, Ev2 registration, and AIRS registration are for deploying to production in tenants such as AME. If you run the generated the service on your local machine or deploy to corp subscriptions only, you can skip the steps.

1. Create a new service in service tree, by simply providing the required information in the [Create New Service Tab](https://microsoftservicetree.com/create/AddNewService).
2. Create the required 4 security groups.
    - The four groups are
        1. Security Group for Admin Access to Development Resources
        2. Security Group for Admin Access to Production Resources
        3. Security Group for Standard Access to Development Resources
        4. Security Group for Standard Access to Production Resources
    - How to create security groups for development resources.
        1. Using your corp account access the [CoreIdentity platform](https://coreidentity.microsoft.com/)
        2. Under common tasks select the "Create Entitlement" option.
        3. Where it says Name for entitlement details enter TM-**'your service's short name'**. Fill out your description and make sure the resource locations match your domain. Press next.
        4. Add your primary and secondary owners. Press next and submit changes.
        5. Once your entitlement creation request has been approved you will have to add a new role, by following the next steps. By default creating an entitlement only creates a "Reader" security group label, and that should not be the role we assign access to resources to.
        6. Press on your entitlement under "Entitlements" on your dashboard. Then select "Roles"
        7. Create a new role and label it the role you would like and add a description. For example "ReadWrite".
        8. After some time, this will automatically create a security group (that is also visible on [Microsoft Entra](https://entra.microsoft.com/#view/Microsoft_AAD_IAM/GroupsManagementMenuBlade/~/AllGroups/menuId/AllGroups)) that you can use for **STANDARD Access to Development Resources**.
        9. Once the role has been created, add your team members to that role too by going to the "Members" tab and pressing the Add button. List your members, press Add, press Next, and select their corresponding roles.
        10. Repeat steps 2-9 , but this time with group name AP-**'your service's short name'**. This will create the service group you will use for **ADMIN Access to Development Resources**.
    - How to create security groups for production resources.
        1. On your SAW, using your AME account access the [OneIdentity platform](https://oneidentity.core.windows.net/Group)
        2. Where it says your group name enter TM-**'your service's short name'** and press next.
        3. Add your primary and secondary owners and any members you would like to this security group.
        4. Press submit changes. This will create the service group you will use for **STANDARD Access to Production Resources**.
        5. Repeat steps 2-4, but this time with group name AP-**'your service's short name'**. This will create the service group you will use for **ADMIN Access to Production Resources**.
3. Link your security groups to your service's access control metadata in Service Tree
    1. Make sure all your security groups have been created, including the development groups being visible in [Microsoft Entra](https://entra.microsoft.com/#view/Microsoft_AAD_IAM/GroupsManagementMenuBlade/~/AllGroups/menuId/AllGroups).
    2. Navigate back to your service's entry in the service tree, and press on the "Metadata" tab.
    3. Under the left hand list, select "Access Control Configuration".
    4. Press on Edit Metadata.
        1. For the tabs for Security Groups for Development Resources, insert the security group that was created tied to the roles you have created. Standard Access will be the SG role in TM-**'your service's short name'** entitlement, and Admin Access will be the SG role in AP-**'your service's short name'** entitlement. They will likely start with "redmond\crg-"
        2. For the tabs for Security Groups for Production Resources, directly insert the security group names. Standard Access will be AME\TM-**'your service's short name'**, and Admin Access will be AME\AP-**'your service's short name'**.
4. Set your release approvers and submitters in your service's Service Tree metadata.
    1. Navigate back to your service's entry in the service tree, and press on the "Metadata" tab.
    2. Under the left hand list, select "OneBranch Services".
    3. Press on Edit Metadata.
        1. For "Release Submitters", insert the security group tied to your Admin Access, that will be the SG role in AP-**'your service's short name'** entitlement. They will likely start with "redmond\crg-"
        2. For "Release Approvers", it must be the AME security group with Admin Access, which will be AME\AP-**'your service's short name'**.
5. Setting up dSCM for JIT Policies.
    1. Onboard your service to the [dSCM platform](https://ui.dscm.core.windows.net/onboard) if it does not already exist. Provide your service tree id and the required information on the page and press Onboard.
    2. Once onboarded, on the menu on the left hand side select "JIT configuration", then Policies.
    3. Select "Subscription". We will be setting a basic policy to meet the security requirements such that acccess to production resources and service releases are secure.
    4. Select "Add policy". For a basic subscription policy, you can set Condition to be Requester is a member of any of the following Security Group(s): AME\AP-**'your service's short name'**. Set how many hours you want to grant access to the subscription, and whatever Approval/rejection message you want.
        - This policy will grant only a member of your security group **"ADMIN Access to Production Resources"**, access to any "Prod" Subscriptions you create that are tied to your service for the specified number of hours when their JIT request is approved.
        - You can vary the policies to best suit your situation, however for security it is best to have a well defined policy for your subscriptions as it defines who can access your production resources and what they can do with it.
    5. For more information about JIT visit [the internal documentation on it](https://review.learn.microsoft.com/en-us/identity/access-management/just-in-time/?branch=main)
6. Register your service with Ev2. Follow steps 2 and 3 in the [Ev2 Service Registration Docs](https://ev2docs.azure.net/features/policies/service-registration.html#2-author-your-service-specification-file) to set up your service with Ev2.
7. Provide your **"STANDARD Access to Production Resources"** security group (AME\TM-**'your service's short name'**) with Submitter (Reader) access to production rollouts, and your **"ADMIN Access to Production Resources"** security group (AME\AP-**'your service's short name'**) with Approver (Contributer) access to rollouts, by following these [instructions](https://ev2docs.azure.net/features/rollout-orchestration/release-pipelines/production.html?q=approver)
8. <a name="ev2-subscriptions"></a>Setting up subscriptions for Ev2

    Currently our Ev2 rollout provisions a subscription to each region that we release in. This was a decision made in order to avoid any restrictions and limitations to resources when they are all located within the same subscription. We follow the subscription provisioning model defined by [Ev2](https://ev2docs.azure.net/features/service-artifacts/actions/subscriptionProvisioningParameters.html#subscription-provisioning-rollout-parameters).
    1. If you do not already have an AIRS registration previously used to create Production Subscriptions, follow the [instructions](https://ev2docs.azure.net/features/service-artifacts/actions/subscriptionProvisioningParameters.html#airs-configuration) to create one.
    2. Give Ev2 permissions to the created registration using [these steps](https://ev2docs.azure.net/features/service-artifacts/actions/subscriptionProvisioningParameters.html#give-ev2-permissions-to-the-created-registration)
    3. The value used for AIRS_REGISTERED_OWNER_OBJECT_ID in the above step is defined in the subscription parameters block in the rollout parameters file for 'airsRegisteredUserPrincipalId' property. It will be later used in the config files when generating with Service Hub
9. Follow the instructions in section [Using Service Hub for generation](#using-service-hub-for-generation) to generate the code and configs.
10. In your generated location, navigate to the README folder for further instructions.

# Using Service Hub for generation

The following steps assume you either have performed the steps in the previous section, or already have an existing service. Choose your form of generation from one of the below options.

We classify users to two types: MSFT user and AKS user. See [user type definitions](https://dev.azure.com/service-hub-flg/service_hub/_wiki/wikis/service_hub.wiki/165/User-types) to decide your type.

As an MSFT user, you can either 1) create a full user environment if you are a new service team or 2) generate some modules and you merge the generated modules to your existing repo.

- [MSFT user generates a full user environment](#msft-user-user-environment): Generate an Azure Devops project, repository, and associated resources such as service connection, application registration, service principal, test/dev pipelines, and the service and its resources.
- [MSFT user generates selected components to a local folder](#msft-user-component): Generate the service and its resources.

As an AKS user, you generate some modules and merge the generated code to existing AKS repos. AKS user is different from MSFT user because AKS has its own build and release system. Service Hub won't generate build and release pipeline for AKS users.

- [AKS user generates selected components to a local folder](#aks-user-component)

If you want to see more concrete use cases, see [how Service Hub generate code to test different use cases](docs/servicehub-use-cases-guide.md).

## <a name="msft-user-user-environment"></a> MSFT user generates a full user environment

- Within the service hub repository, create your own branch.
- Edit the configuration files under [config-files/user-environment](/config-files/user-environment/) to fit your service's needs. Make sure to set the subscription information about your service in [common-config](/config-files/user-environment/common-config.yaml) gained from [Ev2 subscription step](#ev2-subscriptions) of the previous section.
- For more information about our config files and generation capabilities please view [docs/generation-capabilities.md](/docs/generation-capabilities.md).

### Azure subscription and log in

The user environment creation (terraform) uses implicit azure subscription and credentials. This means whatever credentials you are logged in with on your machine will be used to generate the resources.

- Azure Subscription Requirements

    In order to create all the required azure resources (not ADO resources), you must have role permissions in your subscription to:

    1. Create, update and delete resources
    2. Assign roles

    Note that `Contributor` access does not have the permissions to assign roles. `Owner` access has all the permissions necessary.

- To set your subscription for the user environment

    ```bash
    az login
    az account set --subscription $subscriptionId
    ```

### Required variables and secrets

To create the user environment you must edit [variables.tf](user-environment/terraform-files/variables.tf) to store the specific variables that you require.

For the process to function as expected, you **must** have an Azure Key Vault (AKV) that stores "PROJECTPAT" and "READPAT" as secrets. Your AKV can sit in any subscription. If it does not sit in the subscription you are logged into when initializing the user environment, you must define the "key_vault_subscription_id" variable in variables.tf

- PROJECTPAT
  - Reason: The Service Hub Terraform implementation uses the personal access token to access the Azure Devops resources (project, repository, pipelines). The automation script also uses it to handle git operations  (pushing code to the repository after user environment creation).
  - PAT Scope: Must include at least Build (Read & Execute), Code (Full), Project and Team (Read, Write & Manage), Pipeline Resources (Use and manage), Variable Group (Read, create, and manage), Service Connections (Read, query and manage)
  - It acts similarly to how we use the current user's identity to access resources.

- READPAT
  - Reason: It is added as a secret variable to the pipeline definition. The pipeline uses it to access go modules while building the go code. It is also used by developers to build code on their local machine.
  - PAT Scope: Must include at least Code (Read)

If you want to set up a PAT rotator such that your PATS get automatically rotated into your Azure Key Vault, please refer to [1ES secret rotator](https://eng.ms/docs/experiences-devices/opg/office-es365/idc-fundamentals-security/oe-secret-rotator/secret-rotator-tool/guide). Without rotator, you need to update the secrets by yourself. Otherwise the pipeline will fail due to expired PAT.

### Avoiding GOPROXY clash

**Only applicable to AKS team members whose dev machine has the GOPROXY environment variable.**

Service Hub does not use go.goms.io as the goModuleNamePrefix. It uses the created repository directly. Because no proxy is used, you MUST remove GOPROXY from the environment variables.

### Final creation

- Once you have edited variables.tf, run

```bash
make userEnv
```

- Everything will have generated in your chosen Azure Devops Organization. Navigate to the corresponding org, and newly created project and repo. The README folder in your repo will have further instructions.

## <a name="msft-user-component"></a>MSFT user generates selected components

This generates a subset of resources from [MSFT user generates a full user environment](#msft-user-user-environment). It won't create the user environment (ADO project, repo, and pipeline).

- Within the service hub repository, create your own branch.
- Edit the configuration files under [config-files/local-development/external-generation](/config-files/local-development/external-generation) to fit your service's needs. Make sure to set the subscription information about your service in [common-config](/config-files/local-development/external-generation/common-config.yaml) gained from [Ev2 subscription step](#ev2-subscriptions) of the previous section..
- For more information about our config files and generation capabilities please view [docs/generation-capabilities.md](/docs/generation-capabilities.md).

- Create the required directory structure
  - If you do not have, create a projects directory under root

      ```bash
      mkdir -p ~/projects/external
      ```

- Generate with the steps below

    ```bash
    # The default config files will generate a service, shared resources, and pipeline folder under ~/projects/external
    # To generate a specific component only, look into the Makefile to find the targets used by generateAll.
    make generateAll
    # Once generated make your way to the folder
    cd ~/projects/external
    ```

- In your generated location, navigate to the README folder for further instructions

## <a name="aks-user-component"></a> AKS user generates selected components

This generates a service for AKS engineers. The service can be built and deployed with existing AKS Dev Infra tools (AKSbuilder, Overlay deployer, AKS release pipeline, etc)

- Within the service hub repository, create your own branch.
- Edit the configuration files under [config-files/local-development/aks-generation](/config-files/local-development/aks-generation/) to fit your generation needs.
- For more information about our config files and generation capabilities please view [docs/generation-capabilities.md](/docs/generation-capabilities.md).

- If you do not have one create a projects directory under root

  ```bash
  mkdir ~/projects
  ```

- Clone the [aks-rp](https://msazure.visualstudio.com/CloudNativeCompute/_git/aks-rp) repository into the `~/projects` directory, as the service will be generated within `aks-rp`
- Generate the service

  ```bash
  # The default config files will generate a service under ~/projects/aks-rp
  make genAKSService
  # Once generated, make your way to the `mygreeterv3` folder within the `aks-rp` repository and access the README for next steps
  cd ~/projects/aks-rp/mygreeterv3
  ```

## More information on generating modules

- Refer to all of our [generation-capabilities](docs/generation-capabilities.md)

## Dump Information to get a better understanding of existing variables

### MSFT users

- To dump service variables information

  ```bash
  make dumpService
  ```

- To dump resources variables information

  ```bash
  make dumpResources
  ```

- To dump pipeline variables information

  ```bash
  make dumpPipeline
  ```

- To dump available operators information

  ```bash
  make dumpOperators
  ```

- To dump all information stored.

  ```bash
  make dumpAll
  ```

### AKS Users

- To dump service variables information

    ```bash
    make dumpAKSService
    ```
