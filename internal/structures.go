package internal

var Useful_tags = []string{"trunk", "primary", "motorway", "secondary"}

const (
	ConfOptMongoPassword = "MONGO_PASSWORD"
	ConfOptMongoUser     = "MONGO_USER"
	ConfOptMongoDatabase = "MONGO_DATABASE"
)

type NodeID int64

type NeighbourData struct {
	// Represents the data about each neighbour
	NeighbourId NodeID  `bson:"neigbour_id"`
	Dist        float32 `bson:"dist"`
}

type Record struct {
	// Represents each record to insert
	NodeId     NodeID          `bson:"node_id"`
	Neighbours []NeighbourData `bson:"neigbours"`
}

type SingleLaunchRequest struct {
	TruckID        int     `json:"truck_id" validate:"min=0, max=1000000000000, nonnil"`
	OriginLat      float32 `json:"origin_lat" validate:"min=-90, max=90, nonnil"`
	OriginLon      float32 `json:"origin_lon" validate:"min=0, max=180, nonnil"`
	DestinationLat float32 `json:"destination_lat" validate:"min=-90, max=90, nonnil"`
	DestinationLon float32 `json:"destination_lon" validate:"min=0, max=180, nonnil"`
}

type SingleLaunchResponse struct {
	DistanceToOrigin float32 `json:"distance_to_origin"`
	TripDistance     float32 `json:"trip_distance"`
}

type Location struct {
	Latitude  float32 `json:"latitude" bson:"latitude"`
	Longitude float32 `json:"longitude" bson:"longitude"`
}

// Truck
type Truck struct {
	Id              int      `json:"truck_id" bson:"truck_id"`
	FuelConsumption float32  `json:"fuel" bson:"fuel"`         // Spalanie
	Capacity        int      `json:"capacity" bson:"capacity"` // Pojemnosc w tonach
	Location        Location `json:"location" bson:"location"`
}

// Fleet
type Fleet struct {
	Trucks []Truck
}

// Order
// TODO Implement order struct here
type Order struct{}
