package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/sergey-telpuk/gokahoot/di"
	"github.com/sergey-telpuk/gokahoot/graphql"
	"net/http"
	"time"
)

type HttpServer struct {
	DI *di.DI
}

func (s *HttpServer) Run(port string) error {
	return s.routers(s.DI).Run(":" + port)
}

func (s *HttpServer) routers(di *di.DI) *gin.Engine {
	c := cors.Default()
	r := gin.Default()

	graph := graphqlHandler(di)

	r.Any("/query", gin.WrapH(c.Handler(graph)))
	r.GET("/", gin.WrapH(playgroundHandler()))

	return r
}

// Defining the Graphql handler
func graphqlHandler(di *di.DI) *handler.Server {
	schema := graphql.NewExecutableSchema(
		graphql.Config{
			Resolvers: &graphql.Resolver{Di: di},
		},
	)

	h := handler.New(
		schema,
	)

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})

	h.SetQueryCache(lru.New(1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return h
}

// Defining the Playground handler
func playgroundHandler() http.HandlerFunc {
	return playground.Handler("GraphQL", "/query")
}
