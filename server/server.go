package server

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/di"
	"log"
	"os"
)

const defaultPort = "8080"

func Run(di *di.DI) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := newRouter(di)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	if err := router.Run(":" + port); err != nil {
		return errors.New(fmt.Sprintf("Http server was failed %v", err))
	}

	return nil
}
