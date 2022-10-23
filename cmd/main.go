package main

import (
	"Hermes/internal"

	"github.com/sirupsen/logrus"
)

func main() {
	internal.ListPoints()
	// logger := logrus.New()
	// server := createServerFromConfig(logger)
	// server.Run()
}

func createServerFromConfig(logger *logrus.Logger) *internal.HTTPInstanceAPI {
	return internal.NewHTTPInstanceAPI("127.0.0.1:8000")
}
