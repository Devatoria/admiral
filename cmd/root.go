package cmd

import (
	"github.com/Devatoria/admiral/api"
	"github.com/Devatoria/admiral/db"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Address
	RootCmd.PersistentFlags().StringP("address", "a", "127.0.0.1", "API listening address")
	viper.BindPFlag("address", RootCmd.PersistentFlags().Lookup("address"))

	// Port
	RootCmd.PersistentFlags().IntP("port", "p", 3000, "API listening port")
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))

	// Database
	RootCmd.PersistentFlags().String("database-host", "localhost", "Database host")
	RootCmd.PersistentFlags().Int("database-port", 5432, "Database port")
	RootCmd.PersistentFlags().String("database-user", "admiral", "Database user")
	RootCmd.PersistentFlags().String("database-password", "admiral", "Database password")
	RootCmd.PersistentFlags().String("database-name", "admiral", "Database name")

	viper.BindPFlag("database.host", RootCmd.PersistentFlags().Lookup("database-host"))
	viper.BindPFlag("database.port", RootCmd.PersistentFlags().Lookup("database-port"))
	viper.BindPFlag("database.user", RootCmd.PersistentFlags().Lookup("database-user"))
	viper.BindPFlag("database.password", RootCmd.PersistentFlags().Lookup("database-password"))
	viper.BindPFlag("database.name", RootCmd.PersistentFlags().Lookup("database-name"))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/admiral/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// RootCmd is the admiral root command, launching the API and all the stuff
var RootCmd = &cobra.Command{
	Use:   "admiral",
	Short: "Admiral is a Docker Registry admininistration and authentication system",
	Run: func(cmd *cobra.Command, args []string) {
		// Force database init
		_ = db.Instance()

		api.Run(viper.GetString("address"), viper.GetInt("port"))
	},
}
