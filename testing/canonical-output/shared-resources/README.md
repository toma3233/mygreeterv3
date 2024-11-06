# Resource Provisioning

## Create or Update Shared Resources

```bash
# (optional) change parameter values in resources/Main.SharedResources.Template.bicep
make deploy-resources
```
[Optional] Should you want to modify the parameter values, change parameter values in `resources/Main.SharedResources.Template.bicep` and run `make deploy-resources`. Follow the instructions in the README section [Making changes to Bicep Resources](../README.md) at the root of the repo.

## View All Resources and Dependencies

See [shared_resources.md](sharedresources_resources.md). This provides a high-level overview of all your deployments.

This file will only exist after you have run `make deploy-resources`. To see the resources you have created and their dependencies, click the different links in this file. Each link is a different markdown file that is associated with a bicep deployment. Each bicep deployment associated file has:

- list of resources you have created via bicep file
- links to the resources in Azure portal
- the dependencies of each resource
