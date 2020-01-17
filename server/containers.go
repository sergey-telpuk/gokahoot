package server

import (
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/repositories"
	"github.com/sergey-telpuk/gokahoot/services"
	"go.uber.org/dig"
	"log"
)

func BotContainers() *dig.Container {
	container := dig.New()

	if err := container.Provide(db.Init); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if err := container.Provide(repositories.InitTestRepository); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if err := container.Provide(repositories.InitQuestionRepository); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if err := container.Provide(repositories.InitAnswerRepository); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if err := container.Provide(services.InitTestService); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if err := container.Provide(services.InitQuestionService); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	return container
}
