package main

import (
	"Hermes/internal/api"
	"Hermes/internal/http"
	"Hermes/internal/providers"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	confOptMongoPassword = "MONGO_PASSWORD"
	confOptMongoUser     = "MONGO_USER"
	confOptMongoDatabase = "MONGO_DATABASE"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	server := createServerFromConfig(logger, ":8000")
	server.Run()
}

func createServerFromConfig(logger *logrus.Logger, bind string) *http.HTTPInstanceAPI {
	viper.AddConfigPath("/config")
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.ReadInConfig()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster1.yhqlj.mongodb.net/?retryWrites=true&w=majority&ssl=true",
		viper.GetString(confOptMongoUser),
		viper.GetString(confOptMongoPassword))

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.WithError(err).Fatal("could not instatiate ")
	}

	dbController := providers.NewDBController(
		logger.WithField("component", "db"),
		client.Database(viper.GetString(confOptMongoDatabase)),
	)
	instanceAPI := api.NewInstanceAPI(
		logger.WithField("component", "api"),
		dbController,
	)

	return http.NewHTTPInstanceAPI(bind, logger.WithField("component", "http"), instanceAPI)
}
