package cmd

import (
	"fmt"
	"os"

	"github.com/ory/viper"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)
var RootCmd = &cobra.Command{
	Use:   "policy",
	Short: "Policy is a very high performance RBAC Services",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}
func Execute() {
	viper.AutomaticEnv()

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := homedir.Dir()
		// if err != nil {
		// 	er(err)
		// }

		// Search config in home directory with name ".cobra" (without extension).
		// viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
