package main

import (
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/di"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/server"
	"log"
	"os"
)

func main() {
	services := di.New()

	migrate(services)

	if err := server.Run(services); err != nil {
		os.Exit(1)
	}
	os.Exit(0)

	defer services.Clean()
}

func migrate(di *di.DI) {
	sDB := di.Container.Get(db.ContainerName).(*db.Db)

	if err := sDB.GetConn().Exec("PRAGMA foreign_keys=ON").Error; err != nil {
		log.Fatalf("Connection was failed %v", err)
	}

	sDB.GetConn().AutoMigrate(
		&models.Test{},
		&models.Question{},
		&models.Answer{},
	)
}