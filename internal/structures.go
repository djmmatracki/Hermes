package internal

// Graph
/*
	-
*/
type Location struct {
	Latitude  float32
	Longitude float32
}

type category struct {
	fuelConsumption float32 // Spalanie
	capacity        int     //
}

// Truck
type Truck struct {
	Id       int
	Category category
	Location Location
}

// Fleet
type Fleet struct {
	Trucks []Truck
}

// []Truck

// Order
