package services

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"valyria-dc/game"
	"valyria-dc/model"
)

var db *gorm.DB

func Start(listen ...string) error {
	game.OnGameEnd(handleGameEnd)

	g := gin.Default()

	if err := initDb(); err != nil {
		return err
	}

	userEndpoints(g.Group("/api/user"))
	scriptEndpoints(g.Group("/api/script"))
	gameEndpoints(g.Group("/api/games"))
	g.GET("/api/simulator", func(ctx *gin.Context) {
		game.ServeWs(ctx.Writer, ctx.Request)
	})

	g.NoRoute(static.Serve("/", static.LocalFile("frontend", false)))

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
