package api

import (
	"Hermes/internal/providers"
	"context"
)

func (a *InstanceAPI) GetCities(ctx context.Context) ([]providers.City, error) {
	return a.dbController.GetCities(ctx)
}
