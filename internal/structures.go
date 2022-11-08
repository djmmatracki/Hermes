package internal

var Useful_tags = []string{"trunk", "primary"}

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
