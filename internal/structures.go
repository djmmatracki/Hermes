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

type Location struct {
	Latitude  float32 `json:"latitude" bson:"latitude" validate:"min=-90, max=90, nonnil"`
	Longitude float32 `json:"longitude" bson:"longitude" validate:"min=0, max=180, nonnil"`
}

// Truck
type Truck struct {
	Id              int      `json:"truck_id" bson:"truck_id" validate:"min=0, max=1000000000000, nonnil"`
	FuelConsumption float32  `json:"fuel" bson:"fuel" validate:"min=0, max=100, nonnil, nonzero"`
	Capacity        int      `json:"capacity" bson:"capacity" validate:"min=0, max=30, nonnil, nonzero"`
	Location        Location `json:"location" bson:"location"`
}

// Fleet
type Fleet struct {
	Trucks []Truck
}

// Order
// TODO Implement order struct here
type Order struct{}
