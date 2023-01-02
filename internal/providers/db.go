package providers

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// var Useful_tags = []string{"trunk", "primary", "motorway", "secondary"}
// var Max_speed_truck float32 = 75
type UID int64

// type NodeID int64
// type OrderID int64

// type NeighbourData struct {
// 	NeighbourId NodeID  `bson:"neigbour_id"`
// 	Dist        float32 `bson:"dist"`
// }

// type Record struct {
// 	NodeId     NodeID          `bson:"node_id"`
// 	Location   Location        `bson:"location"`
// 	Neighbours []NeighbourData `bson:"neigbours"`
// }

// type SingleLaunchResponse struct {
// 	DistanceToOrigin float32 `json:"distance_to_origin"`
// 	TripDistance     float32 `json:"trip_distance"`
// }

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude" validate:"min=-90, max=90, nonnil"`
	Longitude float64 `json:"longitude" bson:"longitude" validate:"min=0, max=180, nonnil"`
	Name      string  `json:"name" bson:"name"`
}

type Order struct {
	Id          UID      `json:"order_id" bson:"order_id"`
	Name        string   `json:"name" bson:"name"`
	Origin      Location `json:"origin" bson:"origin"`
	Destination Location `json:"destination" bson:"destination"`
	Value       float64  `json:"value" bson:"value" validate:"min=0, nonnil, nonzero"`
	Capacity    int      `json:"capacity" bson:"capacity" validate:"min=0, nonnil, nonzero"`
}

type Truck struct {
	ID              UID      `json:"truck_id" bson:"truck_id" validate:"min=0, max=1000000, nonnil"`
	Name            string   `json:"name" bson:"name"`
	FuelConsumption float64  `json:"fuel" bson:"fuel" validate:"min=0, max=100, nonnil, nonzero"`
	Capacity        int      `json:"capacity" bson:"capacity" validate:"min=0, max=30, nonnil, nonzero"`
	Location        Location `json:"location" bson:"location"`
}

type City struct {
	Name      string  `json:"name" bson:"name"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

func NewDBController(log logrus.FieldLogger, db *mongo.Database) *DBController {
	return &DBController{
		log: log,
		db:  db,
	}
}

type DBController struct {
	log logrus.FieldLogger
	db  *mongo.Database
}
