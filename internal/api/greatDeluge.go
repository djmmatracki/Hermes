package api

import (
	"context"
	"math/rand"
)

func (a *InstanceAPI) GreatDelugeAlgorithm(waterLevel, rainSpeed, groudLevel float64, numberLoops int) (*TrucksAssignmentSolution, error) {
	var order1 OrderID
	var order2 OrderID
	assignment := make(TrucksAssignment)

	if waterLevel <= 0 {
		waterLevel = 1
	}
	if rainSpeed <= 0 {
		rainSpeed = 1.5
	}
	if groudLevel <= 0 {
		groudLevel = 10
	}
	if groudLevel < waterLevel {
		groudLevel = waterLevel + 10
	}

	trucks, err := a.GetTrucks(context.Background())
	if err != nil {
		return nil, err
	}

	orders, err := a.GetOrders(context.Background())
	if err != nil {
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
	bestF = -bestF
	if err != nil {
		return nil, err
	}

	for groudLevel >= waterLevel {
		l := numberLoops
		for l > 0 {
			// Generating new solution in N(x)
			pos1 := TruckID(trucks[rand.Intn(numberOfTrucks)].ID)
			pos2 := TruckID(trucks[rand.Intn(numberOfTrucks)].ID)
			if numberOfTrucks >= numberOfOrders {
				assignment[pos1], assignment[pos2] = assignment[pos2], assignment[pos1]
			} else {
				for {
					order1 = OrderID(orders[rand.Intn(numberOfOrders)].Id)
					order2 = OrderID(orders[rand.Intn(numberOfOrders)].Id)
					if !(checkForValue(order1, assignment) || checkForValue(order2, assignment)) {
						break
					}
					assignment[pos1], assignment[pos2] = order1, order2
				}
			}

			if checkAllowanceCapacity(trucks, orders, assignment) {
				newF, err := a.costAssignment(orders, assignment)
				if err != nil {
					return nil, err
				}
				newF = -newF
				if newF >= waterLevel {
					bestAssignment = assignment
					bestF = newF
					waterLevel *= rainSpeed
				}
			}
			l--
		}
	}

	return &TrucksAssignmentSolution{
		BestTotalIncome: bestF,
		Assignment:      bestAssignment,
	}, nil
}
