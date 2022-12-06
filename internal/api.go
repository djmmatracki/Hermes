package internal

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstanceAPI struct {
	log logrus.FieldLogger

	mongoDatabase *mongo.Database
}

func NewInstanceAPI(log logrus.FieldLogger, mongoDatabase *mongo.Database) *InstanceAPI {
	return &InstanceAPI{
		log:           log,
		mongoDatabase: mongoDatabase,
	}
}

func (a *InstanceAPI) getTrucks(ctx context.Context) ([]Truck, error) {
	var results []Truck
	collection := a.mongoDatabase.Collection("truck")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		a.log.Errorf("error while executing get trucks query")
		return nil, fmt.Errorf("error while executing get trucks query")
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		a.log.Fatal(err)
		return nil, errors.New("Error while")
	}

	return results, nil
}

func (a *InstanceAPI) singleTruckLaunch(truckID int, currentLocation, tripOrigin, tripDestination Location) (*SingleLaunchResponse, error) {
	var distanceToOrigin, distanceToDestination float32
	var truck Truck
	var err error

	collection := a.mongoDatabase.Collection("truck")
	result := collection.FindOne(
		context.TODO(),
		bson.D{{Key: "truck_id", Value: truckID}},
	)

	if err := result.Decode(&truck); err != nil {
		return nil, err
	}

	mainCollection := a.mongoDatabase.Collection("main")

	if distanceToOrigin, err = Astar(mainCollection, currentLocation, tripOrigin); err != nil {
		return nil, err
	}

	if distanceToDestination, err = Astar(mainCollection, tripOrigin, tripDestination); err != nil {
		return nil, err
	}

	return &SingleLaunchResponse{
		TripDistance:     distanceToOrigin + distanceToDestination,
		DistanceToOrigin: distanceToOrigin,
	}, nil
}

func (i *InstanceAPI) simulatedAnneling(trucks []Truck, orders []Order, Nmax int, TStart, TFinal, cooling, k float64) (*TrucksAssignmentSolution, error) {
	assignment := make(TrucksAssignment, 0)

	if len(trucks) != len(orders) {
		return nil, fmt.Errorf("number of trucks and orders should be the same")
	}
	numberOfTrucks := len(trucks)

	for index, truck := range trucks {
		assignment[TruckID(truck.Id)] = OrderID(orders[index].Id)
	}

	orderRevenue := 0
	for _, order in range orders{
		orderRevenue += order.value
	}

	bestAssignment := assignment
	bestF := i.assignmentCost(bestAssignment)
	Temp := TStart
	n := 0

	for Temp > TFinal && n < Nmax {
		pos1 := TruckID(rand.Intn(numberOfTrucks))
		pos2 := TruckID(rand.Intn(numberOfTrucks))

		assignment[pos1], assignment[pos2] = assignment[pos2], assignment[pos1]
		fS := i.assignmentCost(assignment)
		if fS < bestF {
			bestAssignment = assignment
			bestF = fS
		} else {
			delta := fS - bestF
			r := rand.Float64()
			if r < math.Exp(-delta/k*Temp) {
				bestAssignment = assignment
				bestF = fS
			}
		}
		Temp *= cooling
		n++
	}
	// bestAssignment, sum(order.value) - bestF, error
	return &TrucksAssignmentSolution{
		BestTotalIncome: 0,
	}, nil
}

func (i *InstanceAPI) assignmentCost(assingment TrucksAssignment) float64 {
	sum := 0.0
	for _, truckId := range assingment {
		truck, order := getTruck(truckId), getOrder(assignment[truckId])
		sum += truck.FuelConsumption * computeDistance(truck.Location, order.Location)
	}
	return sum
}

// 1/capacity * order.value

// min := 10
// max := 30
// fmt.Println(rand.Intn(max - min) + min)
