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
	TruckID        int     `json:"truck_id"`
	OriginLat      float32 `json:"origin_lat"`
	OriginLon      float32 `json:"origin_lon"`
	DestinationLat float32 `json:"destination_lat"`
	DestinationLon float32 `json:"destination_lon"`
}

type SingleLaunchResponse struct {
	DistanceToOrigin float32 `json:"distance_to_origin"`
	TripDistance     float32 `json:"trip_distance"`
}

// Graph
/*
	-
*/
type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Truck
type Truck struct {
	Id              int     `json:"id"`
	FuelConsumption float32 `json:"fuel"`     // Spalanie
	Capacity        int     `json:"capacity"` // Pojemnosc w tonach
	Location        Location
}

// Fleet
type Fleet struct {
	Trucks []Truck
}

// Order
// TODO Implement order struct here
type Order struct{}
