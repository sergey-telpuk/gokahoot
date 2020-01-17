package server

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sergey-telpuk/gokahoot/graphql"
	"go.uber.org/dig"
)

func newRouter(di *dig.Container) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/query", graphqlHandler(di))
	r.GET("/", playgroundHandler())
	return r
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
