package main

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/graphql"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/repositories"
	"github.com/sergey-telpuk/gokahoot/services"
	"go.uber.org/dig"
	"log"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	di := BuildContainer()

	migrate(di)

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/query", graphqlHandler(di))
	r.GET("/", playgroundHandler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Http server was failed %v", err)
	}

	defer di.Invoke(func(db *gorm.DB) {
		_ = db.Close()
	})
}

// Defining the Graphql handler
func graphqlHandler(di *dig.Container) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.GraphQL(
		graphql.NewExecutableSchema(
			graphql.Config{
				Resolvers: &graphql.Resolver{Di: di},
			},
		),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
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
