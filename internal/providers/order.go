package providers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *DBController) AddOrder(order *Order) error {
	collection := d.db.Collection("orders")
	order.Id = UID(rand.Intn(10000))
	if _, err := collection.InsertOne(context.TODO(), order); err != nil {
		return err
	}
	return nil
}

func (d *DBController) GetOrders(ctx context.Context) ([]Order, error) {
	var results []Order
	collection := d.db.Collection("orders")
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

func (d *DBController) GetOrder(orderID int64) (*Order, error) {
	var order Order
	collection := d.db.Collection("order")
	result := collection.FindOne(
		context.Background(),
		bson.D{{"order_id", orderID}},
	)
	if err := result.Decode(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (d *DBController) DeleteOrder(orderID int64) error {
	collection := d.db.Collection("orders")
	if _, err := collection.DeleteOne(
		context.Background(),
		bson.D{{"order_id", orderID}},
	); err != nil {
		return err
	}
	return nil
}
