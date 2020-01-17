package main

import (
	"github.com/jinzhu/gorm"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/server"
	"go.uber.org/dig"
	"log"
	"os"
)

func main() {
	di := server.BotContainers()

	migrate(di)

	if err := server.Run(di); err != nil {
		os.Exit(1)
	}
	os.Exit(0)

	defer di.Invoke(func(db *gorm.DB) {
		_ = db.Close()
	})
}

func migrate(di *dig.Container) {
	err := di.Invoke(func(db *db.Db) {
		if err := db.GetConn().Exec("PRAGMA foreign_keys=ON").Error; err != nil {
			log.Fatalf("Connection was failed %v", err)
		}

		db.GetConn().AutoMigrate(
			&models.Test{},
			&models.Question{},
			&models.Answer{},
		)
	})

	if err != nil {
		log.Fatalf("Connection was failed %v", err)
	}
}
