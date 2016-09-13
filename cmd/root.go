package cmd

import (
	"os"

	"github.com/Devatoria/admiral/api"
	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Debug
	RootCmd.PersistentFlags().BoolP("debug", "d", true, "Debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	// Address
	RootCmd.PersistentFlags().StringP("address", "a", "127.0.0.1", "API listening address")
	viper.BindPFlag("address", RootCmd.PersistentFlags().Lookup("address"))

	// Port
	RootCmd.PersistentFlags().IntP("port", "p", 3000, "API listening port")
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))

	// Auth
	RootCmd.PersistentFlags().String("auth-issuer", "petito", "Issuer written in bearer token")
	RootCmd.PersistentFlags().Int("auth-token-expiration", 5, "Bearer token expiration time in minutes")
	RootCmd.PersistentFlags().String("auth-certificate", "/certs/server.crt", "Registry certificate path")
	RootCmd.PersistentFlags().String("auth-private-key", "/certs/server.key", "Registry private key path")

	viper.BindPFlag("auth.issuer", RootCmd.PersistentFlags().Lookup("auth-issuer"))
	viper.BindPFlag("auth.token-expiration", RootCmd.PersistentFlags().Lookup("auth-token-expiration"))
	viper.BindPFlag("auth.certificate", RootCmd.PersistentFlags().Lookup("auth-certificate"))
	viper.BindPFlag("auth.private-key", RootCmd.PersistentFlags().Lookup("auth-private-key"))

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

	// Registry
	RootCmd.PersistentFlags().String("registry-address", "http://localhost", "Registry address")
	RootCmd.PersistentFlags().Int("registry-port", 5000, "Registry port")

	viper.BindPFlag("registry.address", RootCmd.PersistentFlags().Lookup("registry-address"))
	viper.BindPFlag("registry.port", RootCmd.PersistentFlags().Lookup("registry-port"))
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
		dbi := db.Instance()
		dbi.AutoMigrate(
			&models.Event{},
			&models.Namespace{},
			&models.Image{},
			&models.Tag{},
			&models.User{},
			&models.Team{},
			&models.TeamNamespaceRight{},
		)

		// Check that issuer exists
		var err error
		if viper.GetString("auth.issuer") == "" {
			panic("Missing auth issuer")
		}

		// Check that certificate file exists
		certificate := viper.GetString("auth.certificate")
		if certificate == "" {
			panic("You must provide a certificate path for auth")
		} else {
			if _, err = os.Stat(certificate); err != nil {
				panic(err)
			}
		}

		// Check that private key file exists
		privateKey := viper.GetString("auth.private-key")
		if privateKey == "" {
			panic("You must provide a private key path for auth")
		} else {
			if _, err = os.Stat(privateKey); err != nil {
				panic(err)
			}
		}

		api.Run(viper.GetString("address"), viper.GetInt("port"))
	},
}
