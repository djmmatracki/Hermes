package internal

import (
	"context"
	"math"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// func contains(s []int64, e int64) bool {
// 	for _, a := range s {
// 		if a == e {
// 			return true
// 		}
// 	}
// 	return false
// }

// func reconstruct_path(came_from map[int64]int64, current int64, stop_id int64, start_id int64) []int64 {

// 	var total_path []int64
// 	total_path = append(total_path, current)
// 	var came_from_keys []int64
// 	for key := range came_from {
// 		came_from_keys = append(came_from_keys, key)
// 	}

// 	for {
// 		if contains(came_from_keys, current) && current != stop_id {
// 			current := came_from[current]
// 			total_path = append(total_path, current)
// 		} else {
// 			break
// 		}
// 	}
// 	total_path = append(total_path, start_id)
// 	return total_path
// }

func computeDistance(startLatLon Location, endLatLon Location) float32 {
	// Function used for calculating
	return float32(math.Sqrt(math.Pow((float64(startLatLon.Latitude-endLatLon.Latitude)), 2) + math.Pow((float64(startLatLon.Longitude-endLatLon.Longitude)), 2)))
}

func A_star(call *mongo.Collection, start Location, stop Location) (float32, error) {
	/*
		Classic A_star algorithm, compute the dystans between two points, knowing there's position [Lat, Lon]
	*/

	start_id, _ := find_nearest_node_id(call, start)
	stop_id, _ := find_nearest_node_id(call, stop)

	// Create opensets that will be our stack where we will put our neighbour nodes
	openSet := []NodeID{start_id}
	// cameFrom will map our route that we takes
	cameFrom := make(map[NodeID]NodeID)

	// gScore is cost of the cheapest path from from start to currently known
	gScore := make(map[NodeID]float32)
	// fScore is current best guess as how we can get to finish
	fScore := make(map[NodeID]float32)

	var adjacency_list_node Record
	result := call.FindOne(
		context.TODO(),
		bson.D{{Key: "node_id", Value: start_id}},
	)
	err := result.Decode(&adjacency_list_node)
	if err != nil {
		return 0, err
	}

	// []NeighbourData
	adjacency_list := adjacency_list_node.Neighbours

	// Makes all cheapest costs to infinity and heuristics update
	for _, neigh_data := range adjacency_list {
		gScore[neigh_data.NeighbourId] = math.MaxFloat32
		fScore[neigh_data.NeighbourId] = math.MaxFloat32
	}

	// Update the cost to to first location to 0 (we are here) and fScore to distance now
	gScore[start_id] = 0
	fScore[start_id] = computeDistance(start, stop) // // FROM BASE heuristic!!!!!

	// While we can visit neighbour node
	for len(openSet) > 0 {
		var adjacency_list_node Record
		// From fScore we take node with the lowest value of: road to node + distance from node to destination
		current := get_lowest_node(fScore, openSet)
		if current == stop_id {
			return gScore[current], err
		}

		result := call.FindOne(
			context.TODO(),
			bson.D{{Key: "node_id", Value: current}},
		)
		err := result.Decode(&adjacency_list_node)
		if err != nil {
			return 0, err
		}

		// Delete actual visiting node
		openSet = find_index_remove(openSet, current)
		for _, neigh := range adjacency_list_node.Neighbours {
			var tentative_gScore float32 = gScore[current] + neigh.Dist // FROM BASE heuristic!!!!!
			if tentative_gScore < gScore[neigh.NeighbourId] {
				cameFrom[neigh.NeighbourId] = current
				gScore[neigh.NeighbourId] = tentative_gScore
				fScore[neigh.NeighbourId] = tentative_gScore + neigh.Dist // FROM BASE heuristic!!!!!
				if not_in_slice(neigh.NeighbourId, openSet) {
					openSet = append(openSet, neigh.NeighbourId)
				}
			}
		}

	}
	return float32(1), nil
}

func not_in_slice(value NodeID, slice []NodeID) bool {
	for _, v := range slice {
		if value == v {
			return false
		}
	}
	return true
}

func find_index_remove(slice []NodeID, curr NodeID) []NodeID {
	var index int = 0

	for i, v := range slice {
		if v == curr {
			index = i
		}
	}
	var list_out1 []NodeID = slice[:index]
	var list_out2 []NodeID = slice[index+1:]
	return append(list_out1, list_out2...)
}

func get_lowest_node(fs map[NodeID]float32, available_nodes []NodeID) NodeID {
	var min_value float32 = math.MaxFloat32
	var pos NodeID = -1
	for _, v := range available_nodes {
		if fs[v] <= min_value {
			min_value = fs[v]
			pos = v
		}
	}
	return pos
}

func find_nearest_node_id(call *mongo.Collection, pos Location) (NodeID, error) {
	var location Record
	result, _ := call.Find(context.TODO(), bson.D{{Key: "location", Value: bson.D{{Key: "latitude", Value: bson.D{{Key: "$gt", Value: pos.Latitude - 0.01}, {Key: "$lt", Value: pos.Latitude + 0.01}}}, {Key: "longitude", Value: bson.D{{Key: "$gt", Value: pos.Longitude - 0.01}, {Key: "$lt", Value: pos.Longitude + 0.01}}}}}})
	err := result.Decode(&location)
	if err != nil {
		return -1, err
	}
	return location.NodeId, nil
}

func choose_best_truck(fleet Fleet, order Order, call *mongo.Collection) Truck {
	/*
		Choose the best truck for order knowing it priority
	*/
	priority := calculate_priority(fleet, order, call)
	var chosen_truck Truck
	var max_prio float32 = 0
	for key, val := range priority {
		if val > max_prio {
			chosen_truck = key
		}
	}
	return chosen_truck
}

func calculate_priority(fleet Fleet, order Order, call *mongo.Collection) map[Truck]float32 {
	/*
		Calculating priority of taking the task by fleet
	*/
	var distance float32
	var time_to_complete float32

	var gauss_time float32
	var gauss_cap float32
	var priority map[Truck]float32 = make(map[Truck]float32)

	for _, truck := range fleet.Trucks {
		// Check if truck have enough capacity
		if truck.Capacity >= order.Capacity {
			// Get the distances
			distance_to_get, _ := A_star(call, truck.Location, order.Location_order)
			distance_to_put, _ := A_star(call, order.Location_order, order.Location_to_deliver)
			// Whole distance
			distance = distance_to_get + distance_to_put

			// Calculate estimated time for arrive
			time_to_complete_assigment := distance / Max_speed_truck

			// Time to complete assignment: Time to complete assignment - time now, in time.Duration format
			time_for_assigment := time.Until(order.Time_delivery)
			// Estimated hours to complete assignment
			time_to_complete = float32(time_for_assigment.Hours())

			// Calculating priority
			gauss_time = gaussian(time_to_complete, time_to_complete_assigment, 1)      // Greater sigma means more spread distribution
			gauss_cap = gaussian(float32(truck.Capacity), float32(order.Capacity), 0.6) // Lower sigma means more focused around mean
			priority[truck] = gauss_cap + gauss_time
		}
	}
	return priority
}

func gaussian(x float32, mean float32, sigma float32) float32 {
	// Calculate normal distribution
	return float32(1 / (sigma * float32(math.Sqrt(2*3.14))) * float32(math.Exp((float64(-1)/2)*math.Pow((float64(x-mean))/float64(sigma), 2))))
}
