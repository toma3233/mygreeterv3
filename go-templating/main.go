package main

import (
	"github.com/spf13/cobra"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Generate",
	Short: "Generate service and/or middleware modules",
	Long: `This service allows you to generate the service or middleware modules or both modules
by taking in the provided config files. The generator requires a config file in order to know where
the service/middleware templates are stored. The service and middleware both require config files
to pass in specified information regarding the modules.`,
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
