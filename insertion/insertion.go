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

<<<<<<< HEAD:scripts/insertion.go
func loadData(osmFile string) (map[internal.NodeID]([]internal.NodeID), map[internal.NodeID]internal.Location) {
=======
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
>>>>>>> b4991dfb4428be6fda0c7cb4115e5f50700f9bca:insertion/insertion.go
	/*
		Opens osm.pbf file from folder, decodes it and make Adjacency list (graph) from it and associate Node_ID with connectiong Nodes

		parms:
			None, in future path to osm.pbf file
		outputs:
			returns map_node_nodes - adjacency list in format: map[int]([]int64)
					map_node_LatLon - map format Node_ID: [Nodes_IDs connected to Node_ID]
	*/
	map_node_nodes := make(map[internal.NodeID][]internal.NodeID)
	map_node_LatLon := make(map[internal.NodeID]internal.Location)

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
				location = internal.Location{Latitude: float32(v.Lat), Longitude: float32(v.Lon)}
				map_node_nodes[internal.NodeID(v.ID)] = []internal.NodeID{}
				map_node_LatLon[internal.NodeID(v.ID)] = location
				nc++
			case *osmpbf.Way:
				// Process Way v.
				for _, i := range v.NodeIDs {
					if check_for_valuable_information(v.Tags, internal.Useful_tags) {
						var actual_values = map_node_nodes[internal.NodeID(i)]
						for _, j := range v.NodeIDs {
							if not_contains(actual_values, internal.NodeID(j)) && i != j {
								actual_values = append(actual_values, internal.NodeID(j))
							}
						}
						map_node_nodes[internal.NodeID(i)] = actual_values
					}
				}
				wc++
			case *osmpbf.Relation:
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}

	}
	fmt.Println(nc, wc, rc)
	return map_node_nodes, map_node_LatLon
}

func not_contains(s []internal.NodeID, e internal.NodeID) bool {
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
	var dist float32
	map_node_nodes, map_node_LatLon := loadData(osmFile)

	// Inserting
	for nodeId, neighbours := range map_node_nodes {
		// Create an instance of a record
<<<<<<< HEAD:scripts/insertion.go
		record := internal.Record{
			NodeId:     nodeId,
			Neighbours: []internal.NeighbourData{},
=======
		record := Record{
			NodeId:     nodeId,
			Neighbours: []NeighbourData{},
>>>>>>> b4991dfb4428be6fda0c7cb4115e5f50700f9bca:insertion/insertion.go
		}
		// Compute distances from the node to their neihgbours
		for _, neighbourId := range neighbours {
			dist = computeDistance(map_node_LatLon[nodeId], map_node_LatLon[neighbourId])
<<<<<<< HEAD:scripts/insertion.go
			record.Neighbours = append(record.Neighbours, internal.NeighbourData{NeighbourId: neighbourId, Dist: dist})
=======
			record.Neighbours = append(record.Neighbours, NeighbourData{NeighbourId: neighbourId, Dist: dist})
>>>>>>> b4991dfb4428be6fda0c7cb4115e5f50700f9bca:insertion/insertion.go
		}

		_, err := collection.InsertOne(context.TODO(), record)
		// id := res.InsertedID

		if err != nil {
			// fmt.Printf("Error occured when node inserting %v\n", err)
			continue
		}
		// fmt.Printf("properly inserted node %d\n", record.NodeId)
	}
	return nil
}

func computeDistance(startLatLon internal.Location, endLatLon internal.Location) float32 {
	// Function used for calculating
	return float32(math.Sqrt(math.Pow((float64(startLatLon.Latitude-endLatLon.Latitude)), 2) + math.Pow((float64(startLatLon.Longitude-endLatLon.Longitude)), 2)))
}

func contains(s string, e []string) bool {
	for _, b := range e {
		if s == b {
			return true
		}
	}
	return false
}

func check_for_valuable_information(tags_map map[string]string, tags_useful []string) bool {
	for _, v := range tags_map {
		if contains(v, tags_useful) {
			return true
		}
	}
	return false
}

<<<<<<< HEAD:scripts/insertion.go
func main() {
	// loadData("greater-london-latest.osm.pbf")
	// Envoke insertion here
	loadData("greater-london-latest.osm.pbf")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster1.yhqlj.mongodb.net/?retryWrites=true&w=majority",
		viper.GetString(internal.ConfOptMongoUser),
		viper.GetString(internal.ConfOptMongoPassword))
=======
// func main() {
// 	// Envoke insertion here
// 	viper.SetConfigFile(".env")
// 	viper.ReadInConfig()

// 	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
// 	mongoURI := fmt.Sprintf(
// 		"mongodb+srv://%s:%s@cluster1.yhqlj.mongodb.net/?retryWrites=true&w=majority",
// 		viper.GetString(confOptMongoUser),
// 		viper.GetString(confOptMongoPassword))
>>>>>>> b4991dfb4428be6fda0c7cb4115e5f50700f9bca:insertion/insertion.go

// 	clientOptions := options.Client().
// 		ApplyURI(mongoURI).
// 		SetServerAPIOptions(serverAPIOptions)

<<<<<<< HEAD:scripts/insertion.go
	client, err := mongo.Connect(context.Background(), clientOptions)
	collection := client.Database(viper.GetString(internal.ConfOptMongoDatabase)).Collection("main")
=======
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	collection := client.Database(viper.GetString(confOptMongoDatabase)).Collection("main")
>>>>>>> b4991dfb4428be6fda0c7cb4115e5f50700f9bca:insertion/insertion.go

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
