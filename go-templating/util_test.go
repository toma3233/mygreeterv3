package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("ConvertStructToMap", func() {
		Context("when given a general struct object", func() {
			It("should convert the struct to a map", func() {
				// Sample struct object
				type Person struct {
					Name  string
					Age   int
					Email string
				}
				person := Person{
					Name:  "Example User",
					Age:   30,
					Email: "exampluser@microsoft.com",
				}

				// Convert the struct to a map
				result := ConvertStructToMap(person)
				// Without a tag, the first letter of the variable name is converted to lowercase
				// Verify the map contains the expected values in yaml formatting
				Expect(result).To(HaveKeyWithValue("name", "Example User"))
				Expect(result).To(HaveKeyWithValue("age", 30))
				Expect(result).To(HaveKeyWithValue("email", "exampluser@microsoft.com"))
			})
		})

		Context("when given a struct object with nested structs", func() {
			It("should recursively convert the nested structs to maps", func() {
				// Sample struct object with nested structs
				type Address struct {
					Street  string
					City    string
					Country string
				}
				type Person struct {
					Name    string
					Age     int
					Address Address
				}
				person := Person{
					Name: "Example User",
					Age:  30,
					Address: Address{
						Street:  "Belmont Ave",
						City:    "Seattle",
						Country: "USA",
					},
				}

				// Convert the struct to a map
				result := ConvertStructToMap(person)
				// Without a tag, the first letter of the variable name is converted to lowercase
				// Verify the map contains the expected values in yaml formatting
				Expect(result).To(HaveKeyWithValue("name", "Example User"))
				Expect(result).To(HaveKeyWithValue("age", 30))
				Expect(result).To(HaveKey("address"))
				Expect(result["address"]).To(HaveKeyWithValue("street", "Belmont Ave"))
				Expect(result["address"]).To(HaveKeyWithValue("city", "Seattle"))
				Expect(result["address"]).To(HaveKeyWithValue("country", "USA"))
			})
		})
		Context("when given pre-defined object with tags", func() {
			It("should convert the struct to a map with tag keys", func() {
				// Create a predefined struct object that already has tags.
				allInput := allInput{
					SharedInput: sharedInput{
						DestinationDirPrefix: "../testing/generated-output",
						DirectoryPath:        "testing/canonical-output/",
					},
					ResourceInput: resourceInput{
						DirectoryName:    "shared-resources",
						Ev2ResourcesName: "official",
						ContactEmail:     "",
						ServiceTreeId:    "3c3a9111-8d68-418f-8868-96641e1510d0",
						TemplateName:     "resourcesTemplate",
					},
					PipelineInput: pipelineInput{
						DirectoryName: "pipeline-files",
						TemplateName:  "pipelineTemplate",
					},
					EnvInformation: envInformation{
						GoModuleNamePrefix:    "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output",
						ServiceConnectionName: "ServiceHubServiceConnection",
					},
					ServiceInput:             serviceInput{DirectoryName: "", ServiceName: "", RunPipeline: false, TemplateName: ""},
					ServiceNameToRunPipeline: nil,
					ExtraInputs:              nil,
					User:                     "",
					GeneratorInputs: map[string]string{
						"mygreeterGoTemplate": "../templates/service_templates",
						"resourcesTemplate":   "../templates/resource-provisioning_templates",
						"pipelineTemplate":    "../templates/pipeline_templates",
						"readmeTemplate":      "../templates/readme_templates",
					},
				}
				// Convert the struct to a map
				result := ConvertStructToMap(allInput)
				// Verify the map contains the expected values with the tag keys
				Expect(result).To(HaveKeyWithValue("sharedInput", HaveKeyWithValue("destinationDirPrefix", "../testing/generated-output")))
				Expect(result).To(HaveKeyWithValue("sharedInput", HaveKeyWithValue("directoryPath", "testing/canonical-output/")))
				Expect(result).To(HaveKeyWithValue("resourceInput", HaveKeyWithValue("directoryName", "shared-resources")))
				Expect(result).To(HaveKeyWithValue("resourceInput", HaveKeyWithValue("ev2ResourcesName", "official")))
				Expect(result).To(HaveKeyWithValue("resourceInput", HaveKeyWithValue("contactEmail", "")))
				Expect(result).To(HaveKeyWithValue("resourceInput", HaveKeyWithValue("serviceTreeId", "3c3a9111-8d68-418f-8868-96641e1510d0")))
				Expect(result).To(HaveKeyWithValue("resourceInput", HaveKeyWithValue("templateName", "resourcesTemplate")))
				Expect(result).To(HaveKeyWithValue("pipelineInput", HaveKeyWithValue("directoryName", "pipeline-files")))
				Expect(result).To(HaveKeyWithValue("pipelineInput", HaveKeyWithValue("templateName", "pipelineTemplate")))
				Expect(result).To(HaveKeyWithValue("envInformation", HaveKeyWithValue("goModuleNamePrefix", "dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output")))
				Expect(result).To(HaveKeyWithValue("envInformation", HaveKeyWithValue("serviceConnectionName", "ServiceHubServiceConnection")))
				Expect(result).To(HaveKeyWithValue("serviceInput", HaveKeyWithValue("directoryName", "")))
				Expect(result).To(HaveKeyWithValue("serviceInput", HaveKeyWithValue("serviceName", "")))
				Expect(result).To(HaveKeyWithValue("serviceInput", HaveKeyWithValue("runPipeline", false)))
				Expect(result).To(HaveKeyWithValue("serviceInput", HaveKeyWithValue("templateName", "")))
				Expect(result).To(HaveKey("serviceNameToRunPipeline"))
				Expect(result).To(HaveKey("extraInputs"))
				Expect(result).To(HaveKeyWithValue("user", ""))
				Expect(result).To(HaveKeyWithValue("generatorInputs", HaveKeyWithValue("mygreeterGoTemplate", "../templates/service_templates")))
				Expect(result).To(HaveKeyWithValue("generatorInputs", HaveKeyWithValue("resourcesTemplate", "../templates/resource-provisioning_templates")))
				Expect(result).To(HaveKeyWithValue("generatorInputs", HaveKeyWithValue("pipelineTemplate", "../templates/pipeline_templates")))
				Expect(result).To(HaveKeyWithValue("generatorInputs", HaveKeyWithValue("readmeTemplate", "../templates/readme_templates")))
			})
		})
	})
})
