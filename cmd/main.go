package main

import (
	"Hermes/internal"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	server := createServerFromConfig(logger, ":8000")
	server.Run()
}

func createServerFromConfig(logger *logrus.Logger, bind string) *internal.HTTPInstanceAPI {
	return internal.NewHTTPInstanceAPI(bind)
}
