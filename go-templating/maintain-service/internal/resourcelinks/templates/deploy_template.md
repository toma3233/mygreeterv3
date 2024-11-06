# Resource Information

## Resource and Resource Type

| Name         | ResourceType |
|:-------------|:-------------|
<<range .Resources>>| [<<.ResourceName>>](<<.ResourcePortalLink>>) |  <<.ResourceType>> |
<<end>>

## Deployments and Dependencies

| Name                  | Depends On  |
|:----------------------|:------------|
<<range .Dependencies>>| <<.DeploymentName>> | <ul><<range .DependencyList>><li><<.>></li><<end>></ul>  |
<<end>>