package db

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres adapter
	"github.com/spf13/viper"
)

var instance *gorm.DB
var once sync.Once

// Instance returns the current database instance, and initialize
// a new one if needed using singleton pattern
func Instance() *gorm.DB {
	once.Do(func() {
		for instance == nil {
			log.Println("Connecting to database...")
			db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
				viper.GetString("database.host"),
				viper.GetInt("database.port"),
				viper.GetString("database.user"),
				viper.GetString("database.name"),
				viper.GetString("database.password"),
			))

			if err != nil {
				log.Printf("Failed to connect to database, waiting and retrying: %v\n", err)
				if err := db.Close(); err != nil {
					panic(fmt.Sprintf("Failed to close database connection: %v", err))
				}

				time.Sleep(3 * time.Second)
				continue
			}

			instance = db
		}
	})

	return instance
}

// Exists returns true if the given model exists with the given condition, false otherwise
func Exists(g *gorm.DB, field string, value, model interface{}) bool {
	var count int
	g.Model(model).Where(fmt.Sprintf("%s = ?", field), value).Count(&count)

	if count == 0 {
		return false
	}

	return true
}
