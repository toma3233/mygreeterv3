# Using "productShortName" vs "servicehub" for internal aks-deployment

- Overall our templates folders no longer hard-code the term "servicehub" when naming resources or deployments in most places.
- However, the one discrepancy currently exists under each server and demoserver's templates folder, in the 'deployment.yaml' file.
  - ***Both the deployment's metadata name and container name must have the "servicehub" prefix as it is required by the AKS standalone logging (where we have hard coded the prefix in the fluent bit log collection configuration).*** Due to such constraint, the two names in deployment.yaml will not be changed by poductShortName (<<.sharedInput.productShortName>>).
  - Internal AKS users can change this to fit their requirements instead, if they no longer want to use "servicehub".
