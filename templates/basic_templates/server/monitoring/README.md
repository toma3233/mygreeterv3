# Monitoring

## Azure Data Explorer (ADX) Dashboard

### Creating the Dashboard

_The file you will need to create the dashboard will be named something like `dashboard/adx-dashboard.json`_.

1. Navigate to <https://dataexplorer.azure.com/dashboards>. You may need to login if you haven't already.
2. In the main dashboard page, select `New dashboard` > Import from file.
<IMG  src="https://learn.microsoft.com/en-us/azure/data-explorer/media/adx-dashboards/import-dashboard-file.png"  alt="Screenshot of dashboard, showing the import from file option."/>
1. Select the monitoring/adx-dashboard.json file.

- You will need to click through a path. (i.e. Linux/ubuntu/home/<username>/go/src/goms.io/aks/rp/<servicename>/monitoring/adx-dasbhoard.json)

4. A prompt will then appear with the title `New Dashboard`. Your dashboard name should already be populated with the name `<your service name> Dashboard`. Click `Create`.
5. Your dashboard will then appear. Make sure to save the dashboard so you have access to it in the future.

Refer to the [documentation](https://learn.microsoft.com/en-us/azure/data-explorer/azure-data-explorer-dashboards#to-create-new-dashboard-from-a-file) for information on Azure Data Explorer Dashboards.

### Sharing the Dashboard

To share the dashboard, you must first manage permissions, grant permissions, and then share the dashboard link.

#### Manage Permissions

1. Select the Share menu item in the top bar of the dashboard.
2. Select Manage permissions from the drop-down.

<IMG  src="https://learn.microsoft.com/en-us/azure/data-explorer/media/adx-dashboards/share-dashboard.png"  alt="Share dashboard drop-down."/>

#### Grant Permissions

To grant permissions to a user in the Dashboard permissions pane:

1. Enter the Azure AD user or Azure AD group in Add new members.
2. In the Permission level, select one of the following values: Can view or Can edit.
3. Select Add.

<IMG  src="https://learn.microsoft.com/en-us/azure/data-explorer/media/adx-dashboards/dashboard-permissions.png"  alt="Manage dashboard permissions."/>

#### Share the dashboard link

To share the dashboard link, do one of the following:

1. Select Share and then select Copy link
2. In the Dashboard permissions window, select Copy link.

#### Change the Parameters

Should you want to change the parameters (i.e. Time range) you are looking at, see the parameters in the top left area of the dashboard and select the different values to change the visuals.

<IMG  src="https://learn.microsoft.com/en-us/azure/data-explorer/media/dashboard-parameters/top-five-states.png
"  alt="Change the parameters."/>

<<if eq .user "aks">>
### Change Hardcoded Values

Currently, `namespace`, used in base query `_baseQuery`, is hardcoded and uses our provided namespace.

#### Change namespace in _baseQuery

1. Make sure you are in editing mode. If not, please click the `viewing` dropdown and switch to editing mode.
2. Click `Base queries` in the to view the parameters.
3. Click the pencil icon on the `_baseQuery` tile to edit the `_baseQuery` query.
4. Add your desired namespaces in the query to replace our example namespace.
5. Click `Run` to run the kql query.
6. You can now click `Done` and save your changes.

## Alert Rules

We provide example adx mon alert rules for error ratio, latency, and QPS in the directory `alert_rules/`.
1. To create these alert rules, please add the example adx mon alert rule files in the [aks-rp-alerts repo](https://msazure.visualstudio.com/CloudNativeCompute/_git/aks-rp-alerts) in the directory `charts/templates`.
2. Please rename the files to replace `svc` with the name of your service.
3. To build and release the alert rules, follow the instructions in the repo.
<<end>>
