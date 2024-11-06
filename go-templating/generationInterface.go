package main

import (
	"path/filepath"
)

// Inputs used across all generation types.
type sharedInput struct {
	// The destination folder where generation takes place.
	DestinationDirPrefix string `yaml:"destinationDirPrefix" example:"~/projects/external"`
	// The path from the root of the repo to where the destination directory is (exclusive).
	// Only relevant for service hub testing.
	DirectoryPath string `yaml:"directoryPath" example: "testing/canonical-output/"`
	ProductDisplayName string `yaml:"productDisplayName" example:"AKS Service Hub"`
	ProductShortName string `yaml:"productShortName" example:"servicehub"`
}

// Inputs for resource generation. Variables are shared with service and pipeline module.
type resourceInput struct {
	DirectoryName string `yaml:"directoryName" example:"shared-resources"`
	// Both resource and service module require these inputs for Ev2
	Ev2ResourcesName string `yaml:"ev2ResourcesName" example:"official"`
	ContactEmail     string `yaml:"contactEmail" example:"test@microsoft.com"`
	ServiceTreeId    string `yaml:"serviceTreeId" example:"3c3a9111-8d68-418f-8868-96641e1510d0"`
	SecondLevelServiceTreeNodeId string `yaml:"secondLevelServiceTreeNodeId" example:"ef733b4f-da1d-4909-8495-73785ce205aa"`
	AdminSecurityGroupId string `yaml:"adminSecurityGroupId" example:"6779196b-41a9-45e9-9cf4-e1192539fbce"`
	AirsRegisteredUserPrincipalId string `yaml:"airsRegisteredUserPrincipalId" example:"23551938-26fb-4713-bb60-456716564972"`
	AirsRegisteredUserTenantId string `yaml:"airsRegisteredUserTenantId" example:"33e01921-4d64-4f8c-a055-5bdaffd5e33d"`
	PcCode string `yaml:"pcCode" example:"P84536"`
	CostCategory string `yaml:"costCategory" example:"FR"`
	// Only used by internal templating code to determine which template to use.
	TemplateName string `yaml:"templateName" example:"resourcesTemplate"`
}

// Inputs for pipeline generation. Variables are shared with service and resource module.
type pipelineInput struct {
	// Directory name shared as both shared-resources and service module reference it
	// to access pipeline files.
	DirectoryName string `yaml:"directoryName" example:"pipeline-files"`
	// Only used by internal templating code to determine which template to use.
	TemplateName string `yaml:"templateName" example:"pipelineTemplate"`
}

// Inputs generated/outputted by terraform resource creation
type envInformation struct {
	GoModuleNamePrefix    string `yaml:"goModuleNamePrefix" example:"shared-resources"`
	ServiceConnectionName string `yaml:"serviceConnectionName" example:"shared-resources"`
}

// Inputs for service generation. Variables are shared with pipeline module.
type serviceInput struct {
	// Directory name is also used by pipeline module to access service pipeline files.
	DirectoryName string `yaml:"directoryName" example:"shared-resources"`
	ServiceName   string `yaml:"serviceName" example:"shared-resources"`
	ContactEmail  string `yaml:"contactEmail" example:"test@microsoft.com"`
	// Used by pipeline module to determine whether or not the service needs to be
	// run in general overarching pipeline.
	RunPipeline bool `yaml:"runPipeline" example:"shared-resources"`
	// Only used by internal templating code to determine which template to use.
	TemplateName string `yaml:"templateName" example:"shared-resources"`
}

type allInput struct {
	SharedInput    sharedInput    `yaml:"sharedInput"`
	ResourceInput  resourceInput  `yaml:"resourceInput"`
	PipelineInput  pipelineInput  `yaml:"pipelineInput"`
	EnvInformation envInformation `yaml:"envInformation"`
	ServiceInput   serviceInput   `yaml:"serviceInput"`

	// The following variables have tags to maintain consistency
	// in variable naming when structs are converted to maps to pass into
	// template executions.

	// Used to store service information generation.
	ServiceNameToRunPipeline map[string]bool   `name:"serviceNameToRunPipeline"`
	ExtraInputs              map[string]string `name:"extraInputs"`
	User                     string            `name:"user"`
	// Only used before templating execution, inaccessible to template files.
	GeneratorInputs map[string]string
}

// This function generates the variables used for templating. It returns
// the destination string and the template name based on generationType
func getTemplateVars(allInput allInput, generationType string) (string, string) {
	switch generationType {
	case serviceStr:
		return filepath.Join(allInput.SharedInput.DestinationDirPrefix, allInput.ServiceInput.DirectoryName), allInput.ServiceInput.TemplateName
	case resourceStr:
		return filepath.Join(allInput.SharedInput.DestinationDirPrefix, allInput.ResourceInput.DirectoryName), allInput.ResourceInput.TemplateName
	case pipelineStr:
		return filepath.Join(allInput.SharedInput.DestinationDirPrefix, allInput.PipelineInput.DirectoryName), allInput.PipelineInput.TemplateName
	default:
		return "", ""
	}
}
