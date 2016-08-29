package db

import (
	"fmt"
	"sync"

	"github.com/Devatoria/admiral/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var instance *gorm.DB
var once sync.Once

// Instance returns the current database instance, and initialize
// a new one if needed using singleton pattern
func Instance() *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.user"),
			viper.GetString("database.name"),
			viper.GetString("database.password"),
		))

		if err != nil {
			panic(err)
		}

		db.AutoMigrate(
			&models.Event{},
			&models.Namespace{},
			&models.Image{},
		)

		instance = db
	})

	return instance
}
