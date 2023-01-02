package providers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *DBController) GetTrucks(ctx context.Context) ([]Truck, error) {
	var results []Truck
	collection := d.db.Collection("truck")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		d.log.Errorf("error while executing get trucks query")
		return nil, fmt.Errorf("error while executing get trucks query")
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		d.log.Error(err)
		return nil, errors.New("error while getting trucks")
	}
	return results, nil
}

func (d *DBController) AddTruck(truck Truck) error {
	collection := d.db.Collection("truck")
	truck.ID = UID(rand.Intn(10000))
	if _, err := collection.InsertOne(context.TODO(), truck); err != nil {
		return err
	}
	return nil
}

func (d *DBController) GetTruck(truckID int64) (*Truck, error) {
	var truck Truck
	collection := d.db.Collection("truck")
	result := collection.FindOne(
		context.Background(),
		bson.D{{"truck_id", truckID}},
	)
	if err := result.Decode(&truck); err != nil {
		return nil, err
	}
	return &truck, nil
}

func (d *DBController) DeleteTruck(truckID int64) error {
	collection := d.db.Collection("truck")
	if _, err := collection.DeleteOne(
		context.Background(),
		bson.D{{"truck_id", truckID}},
	); err != nil {
		return err
	}
	return nil
}
