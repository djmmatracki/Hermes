package internal

import (
	// "Hermes/secret" // jakiś secret_file
	"context"
	// "fmt"
	"log"
	"math"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NeighData struct {
	// Represents the data about each neighbour
	NeighId int
	Dist    float64
}

type Record struct {
	// Represents each record to insert
	NodeId     int64
	Neighbours []NeighData
}

// var uri string = fmt.Sprintf("mongodb://%s:%s@%s:%d", secret.Admin, secret.Passwd, secret.Host, secret.Port)

func InsertNodes(uri string, dbName string, collectionName string, map_node_nodes map[int64]([]int64), map_node_LatLon map[int64](Location)) {
	/*
		Creates client to connect to datebase and inserts all nodes returned by ListPoints() function
		to our database
	*/

	// Connect to the db
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Retrieve a collection from db
	collection := client.Database(dbName).Collection(collectionName)

	// Inserting
	for nodeId, neighbours := range map_node_nodes {
		// Create an instance of a record
		record := Record{
			NodeId:     nodeId,
			Neighbours: []NeighData{},
		}
		// Compute distances from the node to their neihgbours
		for neighbourId := range neighbours {
			record.Neighbours = append(record.Neighbours, NeighData{NeighId: neighbourId, Dist: ComputeDistance(map_node_LatLon[nodeId], map_node_LatLon[int64(neighbourId)])})
		}
		// Insert record to the collection
		_, err := collection.InsertOne(context.TODO(), record)
		// id := res.InsertedID
		if err != nil {
			// fmt.Println("Error occured when %v node inserting", id)
			log.Fatal(err)
		}
	}
}

func ComputeDistance(startLatLon Location, endLatLon Location) float64 {
	// Function used for calculating
	return math.Sqrt(math.Pow((startLatLon.Latitude-endLatLon.Latitude), 2) + math.Pow((startLatLon.Longitude-endLatLon.Longitude), 2))
}
