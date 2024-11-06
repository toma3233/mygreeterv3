package main

import (
	"go-templating/maintain-service/internal/populatemethods"
	"go-templating/maintain-service/internal/resourcelinks"
	"go-templating/templateutil"
	"log"

	"github.com/spf13/cobra"
)

// TODO: Consider moving this function to a separate file
func init() {
	var svcDirName string
	var protoFilePath string
	var svcOutputFileName string
	var deleteGarbageFiles bool
	var templatePrefix string
	var templatePath string
	var envConfig string
	var leftDelim string
	var rightDelim string

	var populateMethodFilesCmd = &cobra.Command{
		Use:   "populateMethodFiles",
		Short: "A command that updates all the method files",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := populatemethods.PopulateAllMethodFiles(svcDirName, protoFilePath, deleteGarbageFiles)
			return err
		},
	}
	var generateResourceFilesCmd = &cobra.Command{
		Use:   "generateResourceFiles",
		Short: "A command that updates all the method files",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := resourcelinks.GenerateAllResourceFiles(svcDirName, svcOutputFileName)
			return err
		},
	}
	var executeTemplatesInFolderCmd = &cobra.Command{
		Use:   "executeTemplatesInFolder",
		Short: "Generates files based on template files under the specified directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			delimiter := []string{}
			delimiter = append(delimiter, leftDelim)
			delimiter = append(delimiter, rightDelim)
			err := templateutil.ExecuteTemplatesInFolder(templatePrefix, templatePath, envConfig, delimiter)
			return err
		},
	}

	rootCmd.AddCommand(populateMethodFilesCmd)
	populateMethodFilesCmd.Flags().StringVar(&svcDirName, "svcDirName", "/container-data", "Name of the directory in the container")
	populateMethodFilesCmd.Flags().StringVar(&protoFilePath, "protoFilePath", "proto/api.proto", "The proto file path relative to service directory")
	populateMethodFilesCmd.Flags().BoolVar(&deleteGarbageFiles, "deleteGarbageFiles", false, "Whether or not to delete files automatically")

	rootCmd.AddCommand(generateResourceFilesCmd)
	generateResourceFilesCmd.Flags().StringVar(&svcDirName, "svcDirName", ".", "Name of the relative path where we start the traversal")
	generateResourceFilesCmd.Flags().StringVar(&svcOutputFileName, "svcMdFileName", "resources.md", "The name of the  file we generate")

	rootCmd.AddCommand(executeTemplatesInFolderCmd)
	executeTemplatesInFolderCmd.Flags().StringVar(&templatePrefix, "templatePrefix", "template-", "The prefix that defines which files are templates.")
	executeTemplatesInFolderCmd.Flags().StringVar(&templatePath, "templatePath", ".", "The path for the folder that is traversed to search for templates.")
	executeTemplatesInFolderCmd.Flags().StringVar(&envConfig, "envConfig", ".", "The env config file used.")
	executeTemplatesInFolderCmd.Flags().StringVar(&leftDelim, "leftDelim", "{{", "The left delimiter used for the template.")
	executeTemplatesInFolderCmd.Flags().StringVar(&rightDelim, "rightDelim", "}}", "The right delimiter used for the template.")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "MaintainService",
	Short: "Maintain service",
	Long:  `This service allows you to make automated changes to the generated directory, which includes but is not limited to the creation, deletion and modification of files in the generated directory.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error executing root command %s", err)
	}
}

func main() {
	Execute()
}
