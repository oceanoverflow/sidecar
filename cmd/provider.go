package cmd

import (
	"fmt"

	"github.com/oceanoverflow/sidecar/provider"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port, dubbo string

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "start provider agent",
	Long:  "start the provider side agent",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify which type of agent you want to run")
			return
		}
		switch args[0] {
		case "small":
			port = viper.GetString("provider.small.port")
			dubbo = viper.GetString("provider.small.dubbo")
		case "medium":
			port = viper.GetString("provider.medium.port")
			dubbo = viper.GetString("provider.medium.dubbo")
		case "large":
			port = viper.GetString("provider.large.port")
			dubbo = viper.GetString("provider.large.dubbo")
		default:
		}
		provider.ServeCommunicate(args[0], port, dubbo)
	},
}
