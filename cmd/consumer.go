package cmd

import (
	"github.com/oceanoverflow/sidecar/consumer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "start consumer agent",
	Long:  "start the consumer side agent",
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetString("consumer.port")
		consumer.ListenAndServe(port)
	},
}
