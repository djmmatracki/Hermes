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
		return nil, errors.New("error while")
	}

	return results, nil
}

func (a *InstanceAPI) getTruck(ctx context.Context, truckID TruckID) (*Truck, error) {
	var truck Truck
	collection := a.mongoDatabase.Collection("truck")
	result := collection.FindOne(ctx, bson.D{{"truck_id", truckID}})

	if err := result.Decode(&truck); err != nil {
		return nil, err
	}
	return &truck, nil
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

func (a *InstanceAPI) simulatedAnneling(orders []Order, Nmax int, TStart, TFinal, cooling, k float32) (*TrucksAssignmentSolution, error) {
	var orderRevenue float32
	assignment := make(TrucksAssignment)

	a.log.Info("getting trucks from mongo db database")
	trucks, err := a.getTrucks(context.Background())
	if err != nil {
		return nil, err
	}

	if len(trucks) != len(orders) {
		return nil, fmt.Errorf("number of trucks and orders should be the same")
	}
	numberOfTrucks := len(trucks)

	a.log.Infof("creating initial assignment for %d", numberOfTrucks)
	for index, truck := range trucks {
		assignment[TruckID(truck.Id)] = orders[index]
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
		pos1 := TruckID(trucks[rand.Intn(numberOfTrucks)].Id)
		pos2 := TruckID(trucks[rand.Intn(numberOfTrucks)].Id)
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
	}, nil
}

func (a *InstanceAPI) costAssignment(orders []Order, assingment TrucksAssignment) (float32, error) {
	var sum float32
	for truckID, order := range assingment {
		truck, err := a.getTruck(context.Background(), truckID)
		if err != nil {
			return 0, err
		}
		sum += computeDistance(truck.Location, order.Location_order) * truck.FuelConsumption
	}
	return sum, nil
}
