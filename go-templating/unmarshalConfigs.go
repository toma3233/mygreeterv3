package main

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	//"path/filepath"
)

// This function unmarshalls the config files that are always required for generation into the umbrella inputs struct.
// All modules require the common config inputs, terraform output inputs, and generator config inputs.
func unmarshalRequiredVars(commonConfigFile string, envInformationFile string, generatorConfigFile string) (allInput, error) {
	allInput := allInput{}
	err := unmarshalIntoStruct(commonConfigFile, &allInput)
	if err != nil {
		return allInput, err
	}
	err = unmarshalIntoStruct(envInformationFile, &allInput)
	if err != nil {
		return allInput, err
	}
	// Unmarshal generator template paths
	err = unmarshalIntoMap(generatorConfigFile, &allInput.GeneratorInputs)
	if err != nil {
		return allInput, err
	}
	// Expand the destination directory if it is a home directory
	allInput.SharedInput.DestinationDirPrefix, err = homedir.Expand(allInput.SharedInput.DestinationDirPrefix)
	if err != nil {
		log.Printf("Error in expanding destination: %v \n", err)
		return allInput, err
	}
	return allInput, nil
}

// This is a general unmarshalling function that unmarshalls a config file into a map.
func unmarshalIntoMap(configFile string, inputs *map[string]string) error {
	conf, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Could not read config file: %v \n", err)
		return err
	}
	err = yaml.Unmarshal(conf, &inputs)
	if err != nil {
		log.Printf("Could not unmarshal config file: %v \n", err)
		return err
	}
	return nil
}

// This function unmarshalls a config file into our defined inputs umbrella inputs struct.
func unmarshalIntoStruct(configFile string, allInput *allInput) error {
	conf, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Could not read %s config file: %v \n", configFile, err)
		return err
	}
	err = yaml.Unmarshal(conf, allInput)
	if err != nil {
		log.Printf("Could not unmarshal %s config file: %v \n", configFile, err)
		return err
	}
	return nil
}
