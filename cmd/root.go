package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var etcd string
var cfg string

func init() {
	rootCmd.AddCommand(consumerCmd)
	rootCmd.AddCommand(providerCmd)
	rootCmd.PersistentFlags().StringVarP(&cfg, "cfg", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&etcd, "etcd", "e", "", "etcd registry")
	viper.BindPFlag("cfg", rootCmd.PersistentFlags().Lookup("cfg"))
	viper.BindPFlag("etcd", rootCmd.PersistentFlags().Lookup("etcd"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.SetConfigFile("/root/dists/config.toml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("can't read config", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "sidecar for consumer and provider services",
	Long:  "sidecar for consumer and provider services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("What do you want? consumer agent or provider agent")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
