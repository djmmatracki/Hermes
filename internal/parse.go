package internal

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	for key := range came_from {
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

// CHANGE THE INPUT VALUES
func A_star(call *mongo.Collection, start Location, stop Location) (float32, error) {
	/*
		Classic A_star algorithm, compute the dystans between two points, knowing there's position [Lat, Lon]

		Parms:
			- start []float64 - starting position
			- stop []float64  - Stoping position
			- adjecency_list map[int64]([]int64) - map of connections in format; node_id: corresponding nodes_ids connected to key
			- nodes_position map[int64]([]float64) - mapping the node_id to real position in format; Node_id: [Latitute, Longitute]

		Returns:
			- float32,err - distance, error
	*/

	// Find Nodes_ID knowing there's position
	start_pos := []float32{start.Latitude, start.Longitude}
	stop_pos := []float32{stop.Latitude, stop.Longitude}

	// start_id = find_nearest_node_id(start_pos)

	// Create opensets that will be our stack where we will put our neighbour nodes
	openSet := []int64{start_id}
	// cameFrom will map our route that we takes
	cameFrom := make(map[int64]int64)

	// gScore is cost of the cheapest path from from start to currently known
	gScore := make(map[int64]float64)
	// fScore is current best guess as how we can get to finish
	fScore := make(map[int64]float64)

	var adjacency_list_node Record
	result := call.FindOne(
		context.TODO(),
		bson.D{{"node_id", start_id}},
	)
	err := result.Decode(&adjacency_list_node)
	if err != nil {
		return 0, err
	}

	// []NeighbourData
	adjacency_list := adjacency_list_node.Neighbours

	// Makes all cheapest costs to infinity and heuristics update
	for node_id, _ := range adjacency_list {
		gScore[int64(node_id)] = math.MaxFloat64
		fScore[int64(node_id)] = 1000 // FROM BASE heuristic!!!!!
	}

	// Update the cost to to first location to 0 (we are here) and fScore to distance now
	gScore[start_id] = 0
	fScore[start_id] = 100 // // FROM BASE heuristic!!!!!

	// While we can visit neighbour node
	for len(openSet) > 0 {
		// From fScore we take node with the lowest value of: road to node + distance from node to destination
		current := get_lowest_node(fScore, openSet)
		if current == stop_id {
			return float32(gScore[current]), err
		}

		// Delete actual visiting node
		openSet = find_index_remove(openSet, current)
		for node_id, _ := range adjacency_list {
			var tentative_gScore float64 = gScore[current] + 1 // FROM BASE heuristic!!!!!
			if tentative_gScore < gScore[int64(node_id)] {
				cameFrom[int64(node_id)] = current
				gScore[int64(node_id)] = tentative_gScore
				fScore[int64(node_id)] = tentative_gScore + 1 // FROM BASE heuristic!!!!!
				if not_in_slice(int64(node_id), openSet) {
					openSet = append(openSet, int64(node_id))
				}
			}
		}

	}
	return float32(1), err
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

// find_nearest_node_id <- zrobić
