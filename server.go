package main

import (
	"github.com/jinzhu/gorm"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
	"github.com/sergey-telpuk/gokahoot/services"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/sergey-telpuk/gokahoot/graphql"
	"go.uber.org/dig"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	di := BuildContainer()

	migrate(di)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{Di: di}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Http server was failed %v", err)
	}

	defer di.Invoke(func(db *gorm.DB) {
		_ = db.Close()
	})
}

func BuildContainer() *dig.Container {
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
	if err := container.Provide(services.InitTestService); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}
	if err := container.Provide(services.InitQuestionService); err != nil {
		log.Fatalf("Provide container was error, error %v", err)
	}

	return container
}

func migrate(di *dig.Container) {
	err := di.Invoke(func(db *db.Db) {
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
