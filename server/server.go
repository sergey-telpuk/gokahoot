package server

import (
	"errors"
	"fmt"
	"github.com/sergey-telpuk/gokahoot/di"
	"os"
)

const defaultPort = "8080"

func Run(di *di.DI) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	s := &HttpServer{di}

	if err := s.Run(port); err != nil {
		return errors.New(fmt.Sprintf("Http server was failed %v", err))
	}

	return nil
}
