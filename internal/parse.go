package internal

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
)

func ListPoints() (map[int64]([]int64), map[int64]([]float64)) {
	/*
		Opens osm.pbf file from folder, decodes it and make Adjacency list (graph) from it and associate Node_ID with connectiong Nodes

		parms:
			None, in future path to osm.pbf file
		outputs:
			returns map_node_nodes - adjacency list in format: map[int]([]int64)
					map_node_LatLon - map format Node_ID: [Nodes_IDs connected to Node_ID]
	*/
	map_node_nodes := make(map[int64]([]int64))
	map_node_LatLon := make(map[int64]([]float64))

	f, err := os.Open("greater-london-latest.osm.pbf")
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
	for {

		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				map_node_nodes[int64(v.ID)] = []int64{}
				map_node_LatLon[v.ID] = []float64{v.Lat, v.Lon}
				nc++
			case *osmpbf.Way:
				// Process Way v.
				for _, i := range v.NodeIDs {
					var actual_values = map_node_nodes[int64(i)]
					for _, j := range v.NodeIDs {
						if not_contains(actual_values, int64(j)) && i != j {
							actual_values = append(actual_values, int64(j))
						}
					}
					map_node_nodes[int64(i)] = actual_values
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

func not_contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return false
		}
	}
	return true
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func reconstruct_path(came_from map[int64]int64, current int64, stop_id int64, start_id int64) []int64 {
	var total_path []int64
	total_path = append(total_path, current)
	var came_from_keys []int64
	for key, _ := range came_from {
		came_from_keys = append(came_from_keys, key)
	}

	for {
		if contains(came_from_keys, current) && current != stop_id {
			current := came_from[current]
			total_path = append(total_path, current)
		} else {
			break
		}
	}
	total_path = append(total_path, start_id)
	return total_path
}
func heuritic_a_star(pos1 []float64, pos2 []float64) float64 {
	return math.Abs(pos1[1]-pos2[1]) + math.Abs(pos1[0]-pos2[0])
}

func A_star(start int64, stop int64, adjacency_list map[int64]([]int64), nodes_position map[int64]([]float64)) []int64 {
	/*
		Classic A_star algorithm, compute the dystans between two points, knowing there's position [Lat, Lon]

		Parms:
			- start []float64 - starting position
			- stop []float64  - Stoping position
			- adjecency_list map[int64]([]int64) - map of connections in format; node_id: corresponding nodes_ids connected to key
			- nodes_position map[int64]([]float64) - mapping the node_id to real position in format; Node_id: [Latitute, Longitute]

		Returns:
			- map[int64]int64
	*/

	if _, ok := adjacency_list[1]; ok {
		delete(adjacency_list, 1)
	}
	// Find Nodes_ID knowing there's position
	start_id := start
	stop_id := stop

	// Create opensets that will be our stack where we will put our neighbour nodes
	openSet := []int64{start_id}
	// cameFrom will map our route that we takes
	cameFrom := make(map[int64]int64)

	// gScore is cost of the cheapest path from from start to currently known
	gScore := make(map[int64]float64)
	// fScore is current best guess as how we can get to finish
	fScore := make(map[int64]float64)

	fmt.Println("1")

	// Makes all cheapest costs to infinity and heuristics update
	for k, _ := range adjacency_list {
		gScore[k] = math.MaxFloat64
		fScore[k] = heuritic_a_star(nodes_position[start_id], nodes_position[k])
	}

	fmt.Println("1")
	// Update the cost to to first location to 0 (we are here) and fScore to distance now
	gScore[start_id] = 0
	fScore[start_id] = heuritic_a_star(nodes_position[start_id], nodes_position[stop_id])

	// While we can visit neighbour node
	for len(openSet) > 0 {
		// From fScore we take node
		current := get_lowest_node(fScore, openSet)
		fmt.Println(current)
		if current == stop_id {
			return reconstruct_path(cameFrom, current, stop_id, start_id)
		}

		openSet = find_index_remove(openSet, current)
		for _, v := range adjacency_list[current] {
			var tentative_gScore float64 = gScore[current] + heuritic_a_star(nodes_position[current], nodes_position[v])
			if tentative_gScore < gScore[v] {
				cameFrom[v] = current
				gScore[v] = tentative_gScore
				fScore[v] = tentative_gScore + heuritic_a_star(nodes_position[v], nodes_position[stop_id])
				if not_in_slice(v, openSet) {
					openSet = append(openSet, v)
				}
			}
		}

	}
	return []int64{-1}
}

func not_in_slice(value int64, slice []int64) bool {
	for _, v := range slice {
		if value == v {
			return false
		}
	}
	return true
}

func find_index_remove(slice []int64, curr int64) []int64 {
	var index int = 0

	for i, v := range slice {
		if v == curr {
			index = i
		}
	}
	var list_out1 []int64 = slice[:index]
	var list_out2 []int64 = slice[index+1:]
	return append(list_out1, list_out2...)
}

func get_lowest_node(fs map[int64]float64, available_nodes []int64) int64 {
	var min_value float64 = math.MaxFloat64
	var pos int64 = -1
	for _, v := range available_nodes {
		if fs[v] <= min_value {
			min_value = fs[v]
			pos = v
		}
	}
	return pos
}

// func get_key_value(pos []float64, nodes_position map[int64]([]float64)) int64 {
// 	/*
// 		Function takes value to associate it with the key value and returns the key value for it

// 		parms:
// 			- pos []float64 - position for which key we want to find
// 			- nodes_position map[int64]([]float64) - map in which we want to find key knowing the value of that key

// 		returns:
// 			- int64 - return key value, in our case it will be Node_ID,
// 						when the value is not existing returns -1
// 	*/
// 	var position int64 = 0
// 	for k, v := range nodes_position {
// 		if v == pos {
// 			fmt.Println("udalo sie")
// 			position = k
// 		}
// 	}
// 	return position
// }

/*
How Nodes and Ways are connected:


type Node struct {
	ID   int64
	Lat  float64
	Lon  float64
	Tags map[string]string
	Info Info
}

type Way struct {
	ID      int64
	Tags    map[string]string
	NodeIDs []int64
	Info    Info
}

*/
