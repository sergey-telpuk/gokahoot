package main

import (
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/server"
	"github.com/sergey-telpuk/gokahoot/services"
	"gopkg.in/gormigrate.v1"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()
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

	sDB.GetConn().AutoMigrate(
		&models.Test{},
		&models.Question{},
		&models.Answer{},
		&models.Answer{},
		&models.Game{},
		&models.Player{},
		&models.PlayerAnswer{},
		&models.ChatMessage{},
	)

	m := gormigrate.New(sDB.GetConn(), gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: guuid.New().String(),
			Migrate: func(tx *gorm.DB) error {
				test := &models.Test{
					UUID: guuid.New().String(),
					Name: "Test1",
				}
				if err := tx.Create(test).Error; err != nil {
					return err
				}

				question := &models.Question{
					UUID:        guuid.New().String(),
					Text:        "FFFFF",
					TestID:      test.ID,
					RightAnswer: 0,
				}

				if err := tx.Create(question).Error; err != nil {
					return err
				}

				for i := 0; i < 4; i++ {
					m := &models.Answer{
						Text:       "Answer",
						Sequential: i,
						QuestionID: question.ID,
					}

					if err := tx.Create(m).Error; err != nil {
						return err
					}
				}

				game := &models.Game{
					Code:    guuid.New().String(),
					TestID:  test.ID,
					Status:  0,
					Players: nil,
				}

				if err := tx.Create(game).Error; err != nil {
					return err
				}

				for i := 0; i < 4; i++ {
					m := &models.Player{
						UUID:   guuid.New().String(),
						Name:   "Player",
						GameID: game.ID,
					}

					if err := tx.Create(m).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}
