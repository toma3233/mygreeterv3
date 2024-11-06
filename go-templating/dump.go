package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var dumpOptions = generateOptions{}

// Define commands and initialize flags for given commands.
func init() {
	// Cobra commands for dumping
	var dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump all variable and operator examples and explanations",
		Run:   dump,
	}
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVar(&dumpOptions.generationType, "generationType", "resource", "Type of folder you want to dump variables for. Must be one of \"pipeline\", \"service\", \"resource\", \"operators\"")
	dumpCmd.Flags().StringVar(&dumpOptions.sConfig, "serviceConfig", "", "Folder path with service config files to generate.")
	dumpCmd.Flags().StringVar(&dumpOptions.cConfig, "commonConfig", "common-config.yaml", "Config file to generate common inputs.")
	dumpCmd.Flags().StringVar(&dumpOptions.envInformation, "envInformation", "env-information.yaml", "Terraform output file.")
	dumpCmd.Flags().StringVar(&dumpOptions.generatorConfig, "generatorConfig", "generator-config.yaml", "Generator config file.")
}

// Dump variable examples and explanations for generating service
func dump(cmd *cobra.Command, args []string) {
	if dumpOptions.generationType == "operators" {
		dumpContent("operators.yaml")
	} else {
		fmt.Println("All generation uses the common config file")
		dumpContent(dumpOptions.cConfig)
		fmt.Println("All generation uses the terraform output file")
		dumpContent(dumpOptions.envInformation)
		switch dumpOptions.generationType {
		case serviceStr, pipelineStr:
			fmt.Println("Service and pipeline generation uses the service config yaml files, the following is an example")
			dumpContent(dumpOptions.sConfig)
		}
	}
}

// Using provided file, dump content where header comments and comments have different colors
// for ease of reading.
func dumpContent(file string) {
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorBlue := "\033[34m"
	colorReset := "\033[0m"

	readFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), "# #") {
			fmt.Println(string(colorGreen), "\n"+fileScanner.Text(), string(colorReset))
		} else {
			if strings.Contains(fileScanner.Text(), "#") {
				fmt.Println(string(colorRed), fileScanner.Text(), string(colorReset))
			} else {
				fmt.Println(fileScanner.Text())
			}
		}
	}
	fmt.Println(string(colorBlue), "----------------------", string(colorReset))
	readFile.Close()
}
