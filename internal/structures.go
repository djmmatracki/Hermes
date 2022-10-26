package internal

// Graph
/*
	-
*/
type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type category struct {
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
