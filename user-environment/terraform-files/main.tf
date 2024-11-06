terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">= 0.1.0"
    }
  }
}

provider "azuredevops" {
  org_service_url       = var.org_url
  personal_access_token = data.azurerm_key_vault_secret.project_pat.value
}
provider "azurerm" {
  features {}
}

# Provider for the key-vault subscription.
# If the key_vault_subscription_id is not provided, by setting the id to null it will use default subscription.
provider "azurerm" {
  features {}
  subscription_id = var.key_vault_subscription_id != "" ? var.key_vault_subscription_id : null
  alias           = "kv_sub"
}

provider "azuread" {}

data "azurerm_subscription" "current" {
}

data "azurerm_subscription" "kv_sub" {
  provider = azurerm.kv_sub
}

# Create a project
# --------------------------------
# Dependencies:
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat"
resource "azuredevops_project" "project" {
  name       = var.project_name
  description        = "This project was created by Terraform"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  lifecycle {
    // This prevents the project from being deleted by Terraform by ignoring changes to the visibility block
    // in the scenario that we are importing an already created project.
    ignore_changes = [
      visibility,
    ]
  }
}

# Create a repository under the above project.
# --------------------------------
# Dependencies:
# "azuredevops_project.project"
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat"
resource "azuredevops_git_repository" "repo" {
  project_id = azuredevops_project.project.id
  name       = var.repo_name
  initialization {
    init_type = "Uninitialized"
  }
  lifecycle {
    // This prevents the repository from being deleted by Terraform by ignoring changes to the initialization block
    // in the scenario that we are importing an already created repository.
    ignore_changes = [
      initialization,
    ]
  }
}

# Create a variable group that stores the readpat.
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "azuredevops_project.project",
# "azuredevops_serviceendpoint_azurerm.service_connection",
# "azurerm_role_assignment.assign_sp_kv_permissions",
# "data.azuread_client_config.current",
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat",
# "data.azurerm_subscription.current",
# "data.azurerm_subscription.kv_sub"
resource "azuredevops_variable_group" "pat_variable_group" {
  project_id   = azuredevops_project.project.id
  name         = "ADO_PAT"
  description  = "Azure Devops PATs"
  allow_access = true

  key_vault {
    name                = var.pat_key_vault_name
    service_endpoint_id = azuredevops_serviceendpoint_azurerm.service_connection.id
  }

  variable {
    name = "READPAT"
  }
  depends_on = [
    azurerm_role_assignment.assign_sp_kv_permissions,
  ]
}

# This creates a pipeline within the created repository and uses
# the given pipeline path as the yaml file to use for the pipeline.
# This pipeline is the test pipeline for the user and
# it is triggered by changes to the default branch (in this case master)
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "azuredevops_git_repository.repo",
# "azuredevops_project.project",
# "azuredevops_serviceendpoint_azurerm.service_connection",
# "azuredevops_variable_group.pat_variable_group",
# "azurerm_role_assignment.assign_sp_kv_permissions",
# "data.azuread_client_config.current",
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat",
# "data.azurerm_subscription.current",
# "data.azurerm_subscription.kv_sub"
resource "azuredevops_build_definition" "test_build" {
  project_id = azuredevops_project.project.id
  name       = "Service Resource And Code Test Pipeline"

  ci_trigger {
    use_yaml = true
  }

  repository {
    repo_type = "TfsGit"
    repo_id     = azuredevops_git_repository.repo.id
    branch_name = azuredevops_git_repository.repo.default_branch
    yml_path    = var.pipeline_path
  }
  variable_groups = [
    azuredevops_variable_group.pat_variable_group.id
  ]
  variable {
    name         = "SUBSCRIPTION_ID"
    value        = data.azurerm_subscription.current.subscription_id
  }
  variable {
    name         = "skipComponentGovernanceDetection"
    value        = true
  }
  variable {
    name         = "DELETE"
    value        = true
  }
  variable {
    name         = "RESOURCES_NAME"
    value        = ""
  }
}

# This creates a pipeline within the created repository and uses
# the given pipeline path as the yaml file to use for the pipeline.
# This is the development pipeline for the user, it doesn't get
# automatically triggered, it must be manually triggered.
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "azuredevops_git_repository.repo",
# "azuredevops_project.project",
# "azuredevops_serviceendpoint_azurerm.service_connection",
# "azuredevops_variable_group.pat_variable_group",
# "azurerm_role_assignment.assign_sp_kv_permissions",
# "data.azuread_client_config.current",
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat",
# "data.azurerm_subscription.current",
# "data.azurerm_subscription.kv_sub"
resource "azuredevops_build_definition" "dev_build" {
  project_id = azuredevops_project.project.id
  name       = "Service Resource And Code Development Pipeline"

  repository {
    repo_type = "TfsGit"
    repo_id     = azuredevops_git_repository.repo.id
    branch_name = azuredevops_git_repository.repo.default_branch
    yml_path    = var.pipeline_path
  }
  variable_groups = [
    azuredevops_variable_group.pat_variable_group.id
  ]
  variable {
    name         = "SUBSCRIPTION_ID"
    value        = data.azurerm_subscription.current.subscription_id
  }
    variable {
    name         = "skipComponentGovernanceDetection"
    value        = true
  }
    variable {
    name         = "DELETE"
    value        = false
  }
    variable {
    name         = "RESOURCES_NAME"
    value        = ""
  }
}


#TODO: provide specific vars to client config and subscription instead of 
#having terraform read from az set up on machine.
# This data source is used to access the configuration of the azure active directory (aad) provider
# It reads the azure information of the person running the terraform, (authenticated principal)
# and can be used to link them as owners to the aad application and service principal.
data "azuread_client_config" "current" {
}

# Key vault that stores the PATs for the pipelines and project
# It uses the provider defined by the key vault subscription id.
# If the key_vault_subscription_id is not provided, key_vault is assumed to exist in default sub.
data "azurerm_key_vault" "key_vault" {
  provider = azurerm.kv_sub
  name                = var.pat_key_vault_name
  resource_group_name = var.key_vault_rg_name
}

# PAT that has permissions to create projects
# Must include at least the following permissions
# - Build (Read & Execute)
# - Code (Full)
# - Project and Team (Read, Write & Manage)
# - Pipeline Resources (Use and manage)
# - Variable Group (Read, create, and manage)
# - Service Connections (Read, query and manage)
data "azurerm_key_vault_secret" "project_pat" {
  provider = azurerm.kv_sub
  name         = "PROJECTPAT"
  key_vault_id = data.azurerm_key_vault.key_vault.id
}

# PAT that has permissions to read code
# Must include at least the following permissions
# - Code (Read)
data "azurerm_key_vault_secret" "read_pat" {
  provider = azurerm.kv_sub
  name         = "READPAT"
  key_vault_id = data.azurerm_key_vault.key_vault.id
}

# Service principals are directly associated with an application within Azure active directory.
# In order to create a service principal, an application must also be created, and the sp will
# be connected to the application.
# Both the application and the service principal is linked to the person generating the terraform
# (authenticated principal). By linking the object id as owner, terraform can have the necessary permissions
# to manage the application, and only the owner will be able to make changes to the application + sp. 
# --------------------------------
# Dependencies:
# "data.azuread_client_config.current"
resource "azuread_application" "app" {
  display_name = var.application_name
  owners                       = [data.azuread_client_config.current.object_id]
  service_management_reference = var.service_tree_id
}

# Dependencies:
# "data.azuread_client_config.current",
# "azuread_application.app"
resource "azuread_service_principal" "sp" {
  client_id = azuread_application.app.client_id
  owners                       = [data.azuread_client_config.current.object_id]
}

# We assign the user identity to have owner permissions such that it is capable
# of CRUD operations within this subscription
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "data.azuread_client_config.current",
# "data.azurerm_subscription.current"
resource "azurerm_role_assignment" "assign_sp_permissions" {
  principal_id           = azuread_service_principal.sp.object_id
  role_definition_name  = "Owner"  
  scope                 = data.azurerm_subscription.current.id
}

# We assign the user identity to have key vault secrets user permissions such
# that the pipeline/service connection can read secrets from the key vault
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "data.azuread_client_config.current",
# "data.azurerm_subscription.kv_sub"
resource "azurerm_role_assignment" "assign_sp_kv_permissions" {
  principal_id           = azuread_service_principal.sp.object_id
  role_definition_name  = "Key Vault Secrets User"  
  scope                 = data.azurerm_subscription.kv_sub.id
}

# Create a service connection endpoint that will be used by pipelines to authenticate interactions with azure
# This service connection uses workload identity federation for authentication 
# and is linked to the previously created application and service principal
# Will be created at the subscription scope and require owner permissions
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "azuredevops_project.project",
# "data.azuread_client_config.current",
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat",
# "data.azurerm_subscription.current"
resource "azuredevops_serviceendpoint_azurerm" "service_connection" {
  project_id                             = azuredevops_project.project.id
  service_endpoint_name                  = var.serviceConnectionName
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "WorkloadIdentityFederation"
  credentials {
    serviceprincipalid = azuread_service_principal.sp.client_id
  }
  azurerm_spn_tenantid      = azuread_service_principal.sp.application_tenant_id
  azurerm_subscription_id   = data.azurerm_subscription.current.subscription_id
  azurerm_subscription_name = data.azurerm_subscription.current.display_name
}

# Generate a federated identity credential for the service connection and 
# that is linked to previously created application
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "azuredevops_project.project",
# "azuredevops_serviceendpoint_azurerm.service_connection",
# "data.azuread_client_config.current",
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat",
# "data.azurerm_subscription.current"
resource "azuread_application_federated_identity_credential" "federated_creds" {
  application_id = azuread_application.app.id
  display_name   = "federated-credential"
  audiences            = ["api://AzureADTokenExchange"]
  issuer              = azuredevops_serviceendpoint_azurerm.service_connection.workload_identity_federation_issuer
  subject             = azuredevops_serviceendpoint_azurerm.service_connection.workload_identity_federation_subject
}

# Grant our generated pipeline permissions to use the service connection
# --------------------------------
# Dependencies:
# "azuread_application.app",
# "azuread_service_principal.sp",
# "azuredevops_build_definition.dev_build",
# "azuredevops_git_repository.repo",
# "azuredevops_project.project",
# "azuredevops_serviceendpoint_azurerm.service_connection",
# "azuredevops_variable_group.pat_variable_group",
# "azurerm_role_assignment.assign_sp_kv_permissions",
# "data.azuread_client_config.current",
# "data.azurerm_key_vault.key_vault",
# "data.azurerm_key_vault_secret.project_pat",
# "data.azurerm_subscription.current",
# "data.azurerm_subscription.kv_sub"
resource "azuredevops_pipeline_authorization" "dev_pipeline_auth" {
  project_id  = azuredevops_project.project.id
  resource_id = azuredevops_serviceendpoint_azurerm.service_connection.id
  type        = "endpoint"
  pipeline_id = azuredevops_build_definition.dev_build.id
}
resource "azuredevops_pipeline_authorization" "test_pipeline_auth" {
  project_id  = azuredevops_project.project.id
  resource_id = azuredevops_serviceendpoint_azurerm.service_connection.id
  type        = "endpoint"
  pipeline_id = azuredevops_build_definition.test_build.id
}