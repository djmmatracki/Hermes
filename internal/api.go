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

func (a *InstanceAPI) getOrders(ctx context.Context) ([]Order, error) {
	var results []Order
	collection := a.mongoDatabase.Collection("order")
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

func (a *InstanceAPI) singleTruckLaunch(truckID int, origin, trip_destiantion, destination Location) (*SingleLaunchResponse, error) {
	// var distanceToOrigin, distanceTo
	var truck Truck
	collection := a.mongoDatabase.Collection("truck")
	result := collection.FindOne(
		context.TODO(),
		bson.D{{Key: "truck_id", Value: truckID}},
	)

	err := result.Decode(&truck)
	if err != nil {
		return nil, err
	}

	collection_node := a.mongoDatabase.Collection("main")

	origin_to_dest, _ := A_star(collection_node, origin, trip_destiantion)
	dest_to_finish, _ := A_star(collection_node, trip_destiantion, destination)
	// origin, destination
	return &SingleLaunchResponse{
		TripDistance:     origin_to_dest + dest_to_finish,
		DistanceToOrigin: origin_to_dest,
	}, nil
}
