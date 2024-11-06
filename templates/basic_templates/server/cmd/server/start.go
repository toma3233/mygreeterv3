// Auto generated. Can be modified.
package main

import (
	"<<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/internal/server"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the service",
	Run:   start,
}

var options = server.Options{}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVar(&options.Port, "port", 50151, "the addr to serve the api on")
	startCmd.Flags().BoolVar(&options.JsonLog, "json-log", false, "The format of the log is json or user friendly key-value pairs")
	startCmd.Flags().IntVar(&options.HTTPPort, "http port", 50161, "the addr to serve the gRPC-Gateway on")
	startCmd.Flags().StringVar(&options.RemoteAddr, "remote-addr", "", "the demo server's addr for this server to connect to")
	startCmd.Flags().Int64Var(&options.IntervalMilliSec, "interval-milli-sec", options.IntervalMilliSec,
		"The interval between two requests. Negative numbers mean sending one request.")
}

func start(cmd *cobra.Command, args []string) {
	server.Serve(options)
}
