package api

import (
	"Hermes/internal/providers"
	"context"
	"math"
	"math/rand"
)

// func (a *InstanceAPI) SingleTruckLaunch(truckID int, currentLocation, tripOrigin, tripDestination Location) (*SingleLaunchResponse, error) {
// 	var distanceToOrigin, distanceToDestination float32
// 	var truck Truck
// 	var err error

// 	collection := a.mongoDatabase.Collection("truck")
// 	result := collection.FindOne(
// 		context.TODO(),
// 		bson.D{{Key: "truck_id", Value: truckID}},
// 	)

// 	if err := result.Decode(&truck); err != nil {
// 		return nil, err
// 	}

// 	mainCollection := a.mongoDatabase.Collection("main")

// 	if distanceToOrigin, err = Astar(mainCollection, currentLocation, tripOrigin); err != nil {
// 		return nil, err
// 	}

// 	if distanceToDestination, err = Astar(mainCollection, tripOrigin, tripDestination); err != nil {
// 		return nil, err
// 	}

//		return &SingleLaunchResponse{
//			TripDistance:     distanceToOrigin + distanceToDestination,
//			DistanceToOrigin: distanceToOrigin,
//		}, nil
//	}
type TrucksAssignment map[TruckID]OrderID

type TrucksAssignmentSolution struct {
	Assignment      TrucksAssignment `json:"assignment"`
	BestTotalIncome float64          `json:"best_total"`
}

func computeDistance(startLatLon providers.Location, endLatLon providers.Location) float64 {
	return math.Sqrt(math.Pow((float64(startLatLon.Latitude-endLatLon.Latitude)), 2) + math.Pow((float64(startLatLon.Longitude-endLatLon.Longitude)), 2))
}

func (a *InstanceAPI) SimulatedAnneling(Nmax int, TStart, TFinal, cooling, k float64) (*TrucksAssignmentSolution, error) {
	var orderRevenue float64
	assignment := make(TrucksAssignment)

	trucks, err := a.GetTrucks(context.Background())
	if err != nil {
		return nil, err
	}

	orders, err := a.GetOrders(context.Background())
	if err != nil {
		return nil, err
	}

	numberOfTrucks := len(trucks)
	for index, truck := range trucks {
		assignment[TruckID(truck.ID)] = OrderID(orders[index].Id)
	}

	for _, order := range orders {
		orderRevenue += order.Value
	}

	bestAssignment := assignment
	bestF, err := a.costAssignment(orders, bestAssignment)
	if err != nil {
		return nil, err
	}

	Temp := TStart
	n := 0

	for Temp > TFinal && n < Nmax {
		pos1 := TruckID(trucks[rand.Intn(numberOfTrucks)].ID)
		pos2 := TruckID(trucks[rand.Intn(numberOfTrucks)].ID)
		assignment[pos1], assignment[pos2] = assignment[pos2], assignment[pos1]

		fS, err := a.costAssignment(orders, assignment)
		if err != nil {
			return nil, err
		}

		if fS < bestF {
			bestAssignment = assignment
			bestF = fS
		} else {
			delta := fS - bestF
			r := rand.Float32()
			if r < float32(math.Exp(float64(-delta/k*Temp))) {
				bestAssignment = assignment
				bestF = fS
			}
		}
		Temp *= cooling
		n++
	}

	return &TrucksAssignmentSolution{
		BestTotalIncome: orderRevenue - bestF,
		Assignment:      bestAssignment,
	}, nil
}

func (a *InstanceAPI) costAssignment(orders []providers.Order, assingment TrucksAssignment) (float64, error) {
	var sum float64
	for truckID, order := range assingment {
		truck, err := a.GetTruck(context.Background(), int64(truckID))
		if err != nil {
			return 0, err
		}
		order, err := a.GetOrder(int64(order))
		if err != nil {
			return 0, err
		}
		sum += computeDistance(truck.Location, order.Origin) * truck.FuelConsumption
	}
	return sum, nil
}
