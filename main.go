package main

import (
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/server"
	"github.com/sergey-telpuk/gokahoot/services"
	"log"
	"os"
)

func main() {
	s, err := services.New()

	if err != nil {
		log.Fatal(err)
	}

	migrate(s)

	if err := server.Run(s); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

	defer s.Clean()
}

func migrate(di *services.DI) {
	sDB := di.Container.Get(db.ContainerName).(*db.Db)

	if err := sDB.GetConn().Exec("PRAGMA foreign_keys=ON").Error; err != nil {
		log.Fatalf("Connection was failed %v", err)
	}

	sDB.GetConn().AutoMigrate(
		&models.Test{},
		&models.Question{},
		&models.Answer{},
		&models.Game{},
		&models.Player{},
	)
}
