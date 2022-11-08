package main

import (
	"Hermes/internal"
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
	server := createServerFromConfig(logger, ":8000")
	server.Run()
}

func createServerFromConfig(logger *logrus.Logger, bind string) *internal.HTTPInstanceAPI {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster1.yhqlj.mongodb.net/?retryWrites=true&w=majority",
		viper.GetString(confOptMongoUser),
		viper.GetString(confOptMongoPassword))

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		logger.WithError(err).Fatal("could not instatiate ")
	}

	instanceAPI := internal.NewInstanceAPI(
		logger.WithField("component", "api"),
		client.Database(viper.GetString(confOptMongoDatabase)),
	)

	return internal.NewHTTPInstanceAPI(bind, logger.WithField("component", "http"), instanceAPI)
}
