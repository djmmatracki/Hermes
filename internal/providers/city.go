package providers

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *DBController) GetCities(ctx context.Context) ([]City, error) {
	var results []City
	collection := d.db.Collection("cities")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		d.log.Errorf("error while executing get orders query")
		return nil, fmt.Errorf("error while executing get orders query")
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		d.log.Error(err)
		return nil, errors.New("error while getting orders")
	}
	return results, nil
}
