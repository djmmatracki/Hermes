package api

import (
	"Hermes/internal/providers"
	"context"

	"gopkg.in/validator.v2"
)

type TruckID int64

func (a *InstanceAPI) GetTrucks(ctx context.Context) ([]providers.Truck, error) {
	return a.dbController.GetTrucks(ctx)
}

func (a *InstanceAPI) GetTruck(ctx context.Context, truckID int64) (*providers.Truck, error) {
	return a.dbController.GetTruck(truckID)
}

func (a *InstanceAPI) AddTruck(truck providers.Truck) error {
	err := validator.Validate(truck)
	if err != nil {
		return err
	}
	return a.dbController.AddTruck(truck)
}

func (a *InstanceAPI) DeleteTruck(truckID int64) error {
	return a.dbController.DeleteTruck(truckID)
}
