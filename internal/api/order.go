package api

import (
	"Hermes/internal/providers"
	"context"
)

type OrderID int64

func (a *InstanceAPI) AddOrder(order *providers.Order) error {
	return a.dbController.AddOrder(order)
}

func (a *InstanceAPI) GetOrders(ctx context.Context) ([]providers.Order, error) {
	return a.dbController.GetOrders(ctx)
}

func (a *InstanceAPI) GetOrder(orderID int64) (*providers.Order, error) {
	return a.dbController.GetOrder(orderID)
}

func (a *InstanceAPI) DeleteOrder(orderID int64) error {
	return a.dbController.DeleteOrder(orderID)
}
