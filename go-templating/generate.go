package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

type generateOptions struct {
	sConfig         string
	cConfig         string
	extraConfig     string
	generatorConfig string
	envInformation  string
	generationType  string
	templateChoice  string
	ignoreOldState  bool
	user            string
	defaultUser     bool
}

var options = generateOptions{}

// Define commands and initialize flags for given commands.
func init() {
	// Cobra command for generating a module
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a module",
		RunE:  generate,
	}
	// Cobra command for generating state file
	var genTemplateSpecCmd = &cobra.Command{
		Use:   "genTemplateSpec",
		Short: "Generate template spec file for provided template",
		RunE:  genTemplateSpec,
	}

	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&options.generationType, "generationType", "resource", "Type of folder you want to generate. Must be one of \"pipeline\", \"service\", or \"resource\"")
	generateCmd.Flags().StringVar(&options.sConfig, "serviceConfig", "../config-files/local-development/external-generation/service-configs", "Folder path with service config files to generate.")
	generateCmd.Flags().StringVar(&options.cConfig, "commonConfig", "../config-files/local-development/external-generation/common-config.yaml", "Config file to generate common inputs.")
	generateCmd.Flags().StringVar(&options.extraConfig, "extraConfig", "", "Config file for extra inputs.")
	generateCmd.Flags().StringVar(&options.envInformation, "envInformation", "../config-files/local-development/external-generation/env-information.yaml", "Terraform output file.")
	generateCmd.Flags().StringVar(&options.generatorConfig, "generatorConfig", "../config-files/generator-config.yaml", "Generator config file.")
	generateCmd.Flags().BoolVarP(&options.ignoreOldState, "ignoreOldState", "i", false, "Whether or not you ignore garbage deletion.")
	generateCmd.Flags().StringVar(&options.user, "user", "external", "User that is running generation.")

	rootCmd.AddCommand(genTemplateSpecCmd)
	genTemplateSpecCmd.Flags().StringVar(&options.templateChoice, "templateChoice", "mygreeterGoTemplate", "Chosen template to generate state file for.")
	genTemplateSpecCmd.Flags().StringVar(&options.generatorConfig, "generatorConfig", "generator-config.yaml", "Generator config file.")
	genTemplateSpecCmd.Flags().BoolVar(&options.defaultUser, "defaultUser", false, "If you want it to add a default user to all the files in the template spec.")
	genTemplateSpecCmd.Flags().StringVar(&options.user, "user", "external", "User that you want to add to all files.")
}

// This function generates the module based on the inputs provided.
// It handles the generation of service, pipeline, and resource modules.
func generate(cmd *cobra.Command, args []string) error {
	// Unmarshal all inputs from the common config files
	allInput, err := unmarshalRequiredVars(options.cConfig, options.envInformation, options.generatorConfig)
	if err != nil {
		return err
	}
	// If we have an extra inputs config file, unmarshal it
	if options.extraConfig != "" {
		err = unmarshalIntoMap(options.extraConfig, &allInput.ExtraInputs)
		if err != nil {
			return err
		}
	}
	allInput.ServiceNameToRunPipeline = map[string]bool{}
	serviceConfigFiles, err := handleServiceConfigFiles(options.sConfig)
	if err != nil {
		return err
	}
	// Unmarshal each service config file into all inputs struct
	// and generate each service module if generationType is service.
	// We also build a map of service names to run in the pipeline for pipeline generation.
	for _, serviceConfigFile := range serviceConfigFiles {
		err := unmarshalIntoStruct(serviceConfigFile, &allInput)
		if err != nil {
			return err
		}
		_, exists := allInput.ServiceNameToRunPipeline[allInput.ServiceInput.DirectoryName]
		if exists {
			fmt.Println("Service name already exists in map. Duplicate service names are not allowed.")
			return errors.New("Duplicate service names")
		}
		allInput.ServiceNameToRunPipeline[allInput.ServiceInput.DirectoryName] = allInput.ServiceInput.RunPipeline
		if options.generationType == serviceStr {
			err = genFromTemplate(allInput, options)
			if err != nil {
				return err
			}
			fmt.Println("Successfully generated", allInput.ServiceInput.DirectoryName)
		}
	}
	// Generate resource or pipeline module.
	if options.generationType != serviceStr {
		err = genFromTemplate(allInput, options)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated", options.generationType)
	}
	// Copy over config files used in generation to the destination folder.
	err = copyConfigFiles(allInput, options.generationType)
	if err != nil {
		return err
	}
	return nil
}

// This function handles one vs multiple service config file case.
// It returns a list of service config files
func handleServiceConfigFiles(sConfig string) ([]string, error) {
	info, err := os.Stat(sConfig)
	serviceConfigFiles := []string{}
	// If the service config path does not exist, we do not generate service modules.
	if err != nil && options.generationType == serviceStr {
		fmt.Printf("Could not find service config path: %v \n", err)
		return serviceConfigFiles, err
	}
	// If the service config path is a directory, we glob all the files in the directory.
	if err == nil {
		if info.IsDir() {
			serviceConfigFiles, err = filepath.Glob(sConfig + "/*.yaml")
			if err != nil {
				fmt.Printf("Could not combine file names: %v \n", err)
				return serviceConfigFiles, err
			}
		} else {
			serviceConfigFiles = append(serviceConfigFiles, sConfig)
		}
	}
	return serviceConfigFiles, nil
}

// This function copies the config files used to generate the current module
// to the destination folder. It will overwrite older config files with the same name.
func copyConfigFiles(allInput allInput, generationType string) error {
	destFolder := filepath.Join(allInput.SharedInput.DestinationDirPrefix, configFilesPath)
	err := os.Mkdir(destFolder, 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Could not create config files destination folder: %v \n", err)
		return err
	}
	// All generations require the environment information file
	err = copyFile(options.envInformation, destFolder)
	if err != nil {
		return err
	}
	// All generations require common config file
	err = copyFile(options.cConfig, destFolder)
	if err != nil {
		return err
	}
	switch generationType {
	//Both service generation and pipeline generation require service config files
	case serviceStr, pipelineStr:
		destFolder = filepath.Join(destFolder, serviceConfigFilesPath)
		err = os.Mkdir(destFolder, 0777)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("Could not create service config files destination folder: %v \n", err)
			return err
		}
		serviceConfigFiles, err := handleServiceConfigFiles(options.sConfig)
		if err != nil {
			return err
		}
		for _, serviceConfigFile := range serviceConfigFiles {
			err = copyFile(serviceConfigFile, destFolder)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

func copyFile(srcPath string, destFolder string) error {
	fileName := filepath.Base(srcPath)
	destPath := filepath.Join(destFolder, fileName)
	input, err := ioutil.ReadFile(srcPath)
	if err != nil {
		fmt.Printf("Error reading source file when copying: %v \n", err)
		return err
	}
	err = ioutil.WriteFile(destPath, input, 0644)
	if err != nil {
		fmt.Printf("Error creating file when copying: %v \n", destPath)
		return err
	}
	return nil
}

// This function generates the template spec file for the template folder of choice.
func genTemplateSpec(cmd *cobra.Command, args []string) error {
	// Create inputs for the templating
	generatorInputs := map[string]string{}
	err := unmarshalIntoMap(options.generatorConfig, &generatorInputs)
	if err != nil {
		return err
	}
	err = generateTemplateSpec(generatorInputs[options.templateChoice], options)
	if err != nil {
		return err
	}
	fmt.Println("Successfully generate template spec for", options.templateChoice)
	return nil
}
