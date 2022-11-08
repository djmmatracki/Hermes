package main

import (
	"context"
	"fmt"

	"Hermes/insertion"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	confOptMongoPassword = "MONGO_PASSWORD"
	confOptMongoUser     = "MONGO_USER"
	confOptMongoDatabase = "MONGO_DATABASE"
)

func main() {
	// Envoke insertion here
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
	collection := client.Database(viper.GetString(confOptMongoDatabase)).Collection("main")

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	var result insertion.Record

	err = collection.FindOne(
		context.TODO(),
		bson.D{{"node_id", 1158813632}},
	).Decode(&result)

	fmt.Println(result.Neighbours)
	// collection.DeleteMany(context.TODO(), bson.D{})
}
