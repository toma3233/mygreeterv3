// Auto generated. Can be modified.
package main

import (
	"dev.azure.com/service-hub-flg/service_hub/_git/service_hub.git/testing/canonical-output/mygreeterv3/server/internal/server"
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

	startCmd.Flags().IntVar(&options.Port, "port", 50051, "the addr to serve the api on")
	startCmd.Flags().BoolVar(&options.JsonLog, "json-log", false, "The format of the log is json or user friendly key-value pairs")
	startCmd.Flags().StringVar(&options.SubscriptionID, "subscription-id", "", "The subscription ID to connect to")
	startCmd.Flags().BoolVar(&options.EnableAzureSDKCalls, "enable-azureSDK-calls", false, "Toggle to run azureSDK CRUDL calls if cluster is enabled with workload-id")
	startCmd.Flags().IntVar(&options.HTTPPort, "http port", 50061, "the addr to serve the gRPC-Gateway on")
	startCmd.Flags().StringVar(&options.RemoteAddr, "remote-addr", "", "the demo server's addr for this server to connect to")
	startCmd.Flags().Int64Var(&options.IntervalMilliSec, "interval-milli-sec", options.IntervalMilliSec,
		"The interval between two requests. Negative numbers mean sending one request.")
	startCmd.Flags().StringVar(&options.IdentityResourceID, "identity-resource-id", "", "the MSI used to authenticate to Azure from E2E env")
}

func start(cmd *cobra.Command, args []string) {
	server.Serve(options)
}
