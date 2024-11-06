# Overview

The hierachy of user visible data model.

- Project/app (mapped to one ADO project, a large service team) (e.g., AKS)
  - project config
  - Shared resources (mapped to a shared or standalone ADO repo)
    - AKS, registry, role assignment
    - Monitoring: workspace, collection rule
  - Service (multiple) (mapped to a shared or standalone ADO repo) (e.g., mygreeter)
    - service config
    - api
    - Microservice (multiple) (e.g., client, demoserver, server under the server directory)
      - Methods
    - test
    - build/push image
    - release/deployments
      - resources
      - binary


# Project

# Service

# Microservice

# Methods
