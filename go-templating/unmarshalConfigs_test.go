package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnmarshalConfigs", func() {
	// Test setup
	commonConfigFile := "go-test-config-files/common-config.yaml"
	envInformationFile := "go-test-config-files/env-information.yaml"
	generatorConfigFile := "go-test-config-files/generator-config.yaml"
	// The expected value of allInput struct
	expectedAllInput := allInput{
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
	Describe("unmarshalRequiredVars", func() {

		It("should unmarshal the required variables", func() {
			// Test execution
			actualAllInput, err := unmarshalRequiredVars(commonConfigFile, envInformationFile, generatorConfigFile)
			// Test verification
			Expect(err).NotTo(HaveOccurred())
			Expect(actualAllInput).To(Equal(expectedAllInput))
		})
		It("should ignore extra information", func() {
			extraCommonConfigFile := "go-test-config-files/common-config-extra.yaml"
			// Test execution
			actualAllInput, err := unmarshalRequiredVars(extraCommonConfigFile, envInformationFile, generatorConfigFile)
			// Test verification
			Expect(err).NotTo(HaveOccurred())
			Expect(actualAllInput).To(Equal(expectedAllInput))
		})
		It("should work with less information", func() {
			lessCommonConfigFile := "go-test-config-files/common-config-less.yaml"
			// Test execution
			actualAllInput, err := unmarshalRequiredVars(lessCommonConfigFile, envInformationFile, generatorConfigFile)
			// Test verification
			Expect(err).NotTo(HaveOccurred())
			lessInput := actualAllInput
			lessInput.PipelineInput = pipelineInput{}
			lessInput.ResourceInput = resourceInput{}
			Expect(actualAllInput).To(Equal(lessInput))
		})
		It("should return an error if commonConfigFile is empty", func() {
			emptyCommonConfigFile := ""
			// Test execution
			_, err := unmarshalRequiredVars(emptyCommonConfigFile, envInformationFile, generatorConfigFile)
			// Test verification
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if envInformationFile is empty", func() {
			emptyEnvInformationFile := ""
			// Test execution
			_, err := unmarshalRequiredVars(commonConfigFile, emptyEnvInformationFile, generatorConfigFile)
			// Test verification
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if generatorConfigFile is empty", func() {
			emptyGeneratorConfigFile := ""
			// Test execution
			_, err := unmarshalRequiredVars(commonConfigFile, envInformationFile, emptyGeneratorConfigFile)
			// Test verification
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if any of the config files cannot be read", func() {
			nonExistentConfigFile := "non-existent-config.yaml"
			// Test execution
			_, err := unmarshalRequiredVars(nonExistentConfigFile, envInformationFile, generatorConfigFile)
			// Test verification
			Expect(err).To(HaveOccurred())
		})

	})
	Describe("unmarshalIntoStruct", func() {
		It("should return an error if config file is invalid", func() {
			// Config file is invalid because types of field is not correct.
			// Service input expects a boolean for "runPipeline" but the config file has a string.
			invalidConfigFile := "go-test-config-files/invalid-config.yaml"
			// Test execution
			allInput := allInput{}
			err := unmarshalIntoStruct(invalidConfigFile, &allInput)
			// Test verification
			Expect(err).To(HaveOccurred())
		})
	})
})
