package insertion

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"

	"Hermes/internal"

	"github.com/qedus/osmpbf"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	confOptMongoPassword = "MONGO_PASSWORD"
	confOptMongoUser     = "MONGO_USER"
	confOptMongoDatabase = "MONGO_DATABASE"
)

type NodeID int64

type NeighbourData struct {
	// Represents the data about each neighbour
	NeighbourId NodeID  `bson:"neigbour_id"`
	Dist        float64 `bson:"dist"`
}

type Record struct {
	// Represents each record to insert
	NodeId     NodeID          `bson:"node_id"`
	Neighbours []NeighbourData `bson:"neigbours"`
}

func loadData(osmFile string) (map[NodeID]([]NodeID), map[NodeID]internal.Location) {
	/*
		Opens osm.pbf file from folder, decodes it and make Adjacency list (graph) from it and associate Node_ID with connectiong Nodes

		parms:
			None, in future path to osm.pbf file
		outputs:
			returns map_node_nodes - adjacency list in format: map[int]([]int64)
					map_node_LatLon - map format Node_ID: [Nodes_IDs connected to Node_ID]
	*/
	map_node_nodes := make(map[NodeID][]NodeID)
	map_node_LatLon := make(map[NodeID]internal.Location)

	// Read OSM
	f, err := os.Open(osmFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	var nc, wc, rc uint64
	var location internal.Location

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				location = internal.Location{Latitude: v.Lat, Longitude: v.Lon}
				map_node_nodes[NodeID(v.ID)] = []NodeID{}
				map_node_LatLon[NodeID(v.ID)] = location
				nc++
			case *osmpbf.Way:
				// Process Way v.
				for _, i := range v.NodeIDs {
					var actual_values = map_node_nodes[NodeID(i)]
					for _, j := range v.NodeIDs {
						if not_contains(actual_values, NodeID(j)) && i != j {
							actual_values = append(actual_values, NodeID(j))
						}
					}
					map_node_nodes[NodeID(i)] = actual_values
				}
				wc++
			case *osmpbf.Relation:
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	return map_node_nodes, map_node_LatLon
}

func not_contains(s []NodeID, e NodeID) bool {
	for _, a := range s {
		if a == e {
			return false
		}
	}
	return true
}

func insertNodes(collection *mongo.Collection, osmFile string) error {
	/*
		Creates client to connect to datebase and inserts all nodes returned by ListPoints() function
		to our database
	*/
	var dist float64
	map_node_nodes, map_node_LatLon := loadData(osmFile)

	// Inserting
	for nodeId, neighbours := range map_node_nodes {
		// Create an instance of a record
		record := Record{
			NodeId:     nodeId,
			Neighbours: []NeighbourData{},
		}
		// Compute distances from the node to their neihgbours
		for _, neighbourId := range neighbours {
			dist = computeDistance(map_node_LatLon[nodeId], map_node_LatLon[neighbourId])
			record.Neighbours = append(record.Neighbours, NeighbourData{NeighbourId: neighbourId, Dist: dist})
		}

		_, err := collection.InsertOne(context.TODO(), record)
		// id := res.InsertedID

		if err != nil {
			fmt.Printf("Error occured when node inserting %v\n", err)
			continue
		}
		fmt.Printf("properly inserted node %d\n", record.NodeId)
	}
	return nil
}

func computeDistance(startLatLon internal.Location, endLatLon internal.Location) float64 {
	// Function used for calculating
	return math.Sqrt(math.Pow((startLatLon.Latitude-endLatLon.Latitude), 2) + math.Pow((startLatLon.Longitude-endLatLon.Longitude), 2))
}

// func main() {
// 	// Envoke insertion here
// 	viper.SetConfigFile(".env")
// 	viper.ReadInConfig()

// 	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
// 	mongoURI := fmt.Sprintf(
// 		"mongodb+srv://%s:%s@cluster1.yhqlj.mongodb.net/?retryWrites=true&w=majority",
// 		viper.GetString(confOptMongoUser),
// 		viper.GetString(confOptMongoPassword))

// 	clientOptions := options.Client().
// 		ApplyURI(mongoURI).
// 		SetServerAPIOptions(serverAPIOptions)

// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	collection := client.Database(viper.GetString(confOptMongoDatabase)).Collection("main")

// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// collection.DeleteMany(context.TODO(), bson.D{})
// 	for _, osmFile := range os.Args[1:] {
// 		if err := insertNodes(collection, osmFile); err != nil {
// 			fmt.Printf("error while inserting data: %v", err)
// 			return
// 		}
// 	}
// }
