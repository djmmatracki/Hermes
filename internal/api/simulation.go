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

func checkForValue(order OrderID, assignment map[TruckID]OrderID) bool {
	for _, value := range assignment {
		if value == order {
			return true
		}
	}
	return false
}

func checkAllowanceCapacity(trucks []providers.Truck, orders []providers.Order, assignment TrucksAssignment) bool {
	for truckID, orderID := range assignment {
		for _, truck := range trucks {
			for _, order := range orders {
				if truck.ID == providers.UID(truckID) && order.Id == providers.UID(orderID) {
					if truck.Capacity < order.Capacity {
						return false
					}
				}
			}
		}
	}
	return true
}

func (a *InstanceAPI) SimulatedAnneling(Nmax int, TStart, TFinal, cooling, k float64) (*TrucksAssignmentSolution, error) {
	var order1 OrderID
	var order2 OrderID
	assignment := make(TrucksAssignment)

	if Nmax <= 0 {
		Nmax = 5
	}
	if TStart <= 0 {
		TStart = 10
	}
	if TFinal <= 0 {
		TFinal = 2
	}
	if TStart <= TFinal {
		TStart = TFinal + 1
	}
	if cooling <= 0 {
		cooling = 0.9
	}

	trucks, err := a.GetTrucks(context.Background())
	if err != nil {
		a.log.Error("Error while getting trucks")
		return nil, err
	}

	orders, err := a.GetOrders(context.Background())
	if err != nil {
		a.log.Error("Error while getting orders")
		return nil, err
	}

	numberOfTrucks := len(trucks)
	numberOfOrders := len(orders)

	for index, truck := range trucks {
		if index >= numberOfOrders {
			assignment[TruckID(truck.ID)] = 0
		} else {
			assignment[TruckID(truck.ID)] = OrderID(orders[index].Id)
		}
	}

	bestAssignment := assignment
	bestF, err := a.costAssignment(orders, bestAssignment)
	if err != nil {
		a.log.Error("Error while executing cost assignment")
		return nil, err
	}

	Temp := TStart
	n := 0

	for Temp > TFinal && n < Nmax {
		pos1 := TruckID(trucks[rand.Intn(numberOfTrucks)].ID)
		pos2 := TruckID(trucks[rand.Intn(numberOfTrucks)].ID)
		if numberOfTrucks >= numberOfOrders {
			assignment[pos1], assignment[pos2] = assignment[pos2], assignment[pos1]
		} else {
			for {
				order1 = OrderID(orders[rand.Intn(numberOfOrders)].Id)
				order2 = OrderID(orders[rand.Intn(numberOfOrders)].Id)
				if !(checkForValue(order1, assignment) || checkForValue(order2, assignment)) {
					assignment[pos1], assignment[pos2] = order1, order2
					break
				}
			}
		}

		if checkAllowanceCapacity(trucks, orders, assignment) {
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
	}

	return &TrucksAssignmentSolution{
		BestTotalIncome: -bestF,
		Assignment:      bestAssignment,
	}, nil
}

func (a *InstanceAPI) costAssignment(orders []providers.Order, assingment TrucksAssignment) (float64, error) {
	var sum float64
	for truckID, orderID := range assingment {
		truck, err := a.GetTruck(context.Background(), int64(truckID))
		if err != nil {
			a.log.Error("Error while getting truck")
			return 0, err
		}
		if orderID == 0 {
			continue
		}

		order, err := a.GetOrder(int64(orderID))
		if err != nil {
			a.log.Error("Error while getting order")
			return 0, err
		}
		sum += math.Exp(computeDistance(truck.Location, order.Origin)*(-0.3)) * order.Value
	}
	return -sum, nil
}
