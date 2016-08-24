package cmd

import (
	"github.com/Devatoria/admiral/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Address
	RootCmd.PersistentFlags().StringP("address", "a", "127.0.0.1", "API listening address")
	viper.BindPFlag("address", RootCmd.PersistentFlags().Lookup("address"))

	// Port
	RootCmd.PersistentFlags().IntP("port", "p", 3000, "API listening port")
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
}

// RootCmd is the admiral root command, launching the API and all the stuff
var RootCmd = &cobra.Command{
	Use:   "admiral",
	Short: "Admiral is a Docker Registry admininistration and authentication system",
	Run: func(cmd *cobra.Command, args []string) {
		api.Run(viper.GetString("address"), viper.GetInt("port"))
	},
}
