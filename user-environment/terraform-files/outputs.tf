output "repo_url" {
    value = azuredevops_git_repository.repo.remote_url
}
output "project_name" {
    value = azuredevops_project.project.name
}
output "repo_name" {
    value = azuredevops_git_repository.repo.name
}
output "service_connection_name" {
    value = azuredevops_serviceendpoint_azurerm.service_connection.service_endpoint_name
}
// PATS are utilised by environment creation script to set up git configuration.
output "project_pat" {
    value = data.azurerm_key_vault_secret.project_pat.value
    sensitive = true
}
output "read_pat" {
    value = data.azurerm_key_vault_secret.read_pat.value
    sensitive = true
}
