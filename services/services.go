package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"valyria-dc/model"
)

var db *gorm.DB

func Start(listen ...string) error {
	g := gin.Default()

	if err := initDb(); err != nil {
		return err
	}

	return g.Run(listen...)
}

func initDb() error {
	vendor := os.Getenv("DB_VENDOR")
	dsn := os.Getenv("DB_DSN")

	var err error

	if vendor == "postgres" {
		log.Println("Using PostgresSQL vendor")
		db, err = gorm.Open(postgres.Open(dsn))
		if err != nil {
			return err
		}
	} else if vendor == "sqlite" {
		log.Println("Using SQLite3 vendor")
		db, err = gorm.Open(sqlite.Open(dsn))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown db vendor: %s", vendor)
	}

	return model.AutoMigrate(db)
}
