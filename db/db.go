package db

import (
	"fmt"
	"sync"

	"github.com/Devatoria/admiral/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var instance *gorm.DB
var once sync.Once

// Instance returns the current database instance, and initialize
// a new one if needed using singleton pattern
func Instance() *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.name"),
		))

		if err != nil {
			panic(err)
		}

		db.AutoMigrate(
			&models.Event{},
			&models.EventRequest{},
			&models.EventSource{},
			&models.EventTarget{},
		)

		instance = db
	})

	return instance
}
