package internal

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstanceAPI struct {
	log logrus.FieldLogger

	mongoDatabase *mongo.Database
}

func NewInstanceAPI(log logrus.FieldLogger, mongoDatabase *mongo.Database) *InstanceAPI {
	return &InstanceAPI{
		log:           log,
		mongoDatabase: mongoDatabase,
	}
}

func (a *InstanceAPI) getTrucks(ctx context.Context) {
	var results []bson.M
	collection := a.mongoDatabase.Collection("truck")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println("hello")
		log.Fatal(err)
		// return nil, errors.New("")
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		// return nil, errors.New("")
	}

	log.Println(results)
	for _, result := range results {
		log.Println(result)
	}
}
