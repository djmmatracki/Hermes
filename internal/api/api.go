package api

import (
	"Hermes/internal/providers"

	"github.com/sirupsen/logrus"
)

type InstanceAPI struct {
	log          logrus.FieldLogger
	dbController *providers.DBController
}

func NewInstanceAPI(log logrus.FieldLogger, dbController *providers.DBController) *InstanceAPI {
	return &InstanceAPI{
		log:          log,
		dbController: dbController,
	}
}
