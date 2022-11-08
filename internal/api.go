package internal

import (
	"context"
	"errors"

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

func (a *InstanceAPI) getTrucks(ctx context.Context) ([]Truck, error) {
	var results []Truck
	collection := a.mongoDatabase.Collection("truck")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		a.log.Fatal(err)
		return nil, errors.New("")
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		a.log.Fatal(err)
		return nil, errors.New("")
	}

	return results, nil
}

func (a *InstanceAPI) singleTruckLaunch(truckID int, origin, destination Location) (*SingleLaunchResponse, error) {
	// var distanceToOrigin, distanceTo
	collection := a.mongoDatabase.Collection("truck")
	collection.FindOne(
		context.TODO(),
		bson.D{{"truck_id", truckID}},
	)
	// Astar(collection, truck_origin, trip_origin Location) -> dystans
	// Astar(collection, trip_origin, trip_destination Location) -> dystans
	// origin, destination
	return nil, nil
}
