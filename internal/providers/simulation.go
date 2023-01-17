package providers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Point struct {
}

func (d *DBController) FindNearestNode(location Location) (*Point, error) {
	var nearestLocation Point
	collection := d.db.Collection("points")
	result, _ := collection.Find(context.TODO(),
		bson.D{{
			Key: "location",
			Value: bson.D{{
				Key: "latitude",
				Value: bson.D{{
					Key:   "$gt",
					Value: location.Latitude - 0.01,
				}, {
					Key:   "$lt",
					Value: location.Latitude + 0.01,
				}}}, {
				Key: "longitude",
				Value: bson.D{{
					Key:   "$gt",
					Value: location.Longitude - 0.01,
				}, {
					Key:   "$lt",
					Value: location.Longitude + 0.01,
				}}}}}})
	if err := result.Decode(&nearestLocation); err != nil {
		return nil, err
	}
	return &nearestLocation, nil
}
