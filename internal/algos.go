package internal

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func computeDistance(startLatLon Location, endLatLon Location) float32 {
	return float32(math.Sqrt(math.Pow((float64(startLatLon.Latitude-endLatLon.Latitude)), 2) + math.Pow((float64(startLatLon.Longitude-endLatLon.Longitude)), 2)))
}

func Astar(call *mongo.Collection, start Location, stop Location) (float32, error) {
	startID, _ := find_nearest_node_id(call, start)
	stopID, _ := find_nearest_node_id(call, stop)
	openSet := []NodeID{startID}

	cameFrom := make(map[NodeID]NodeID)

	// gScore is cost of the cheapest path from from start to currently known
	gScore := make(map[NodeID]float32)
	// fScore is current best guess as how we can get to finish
	fScore := make(map[NodeID]float32)

	var adjacency_list_node Record
	result := call.FindOne(
		context.TODO(),
		bson.D{{Key: "node_id", Value: startID}},
	)
	err := result.Decode(&adjacency_list_node)
	if err != nil {
		return 0, err
	}

	adjacency_list := adjacency_list_node.Neighbours
	for _, neigh_data := range adjacency_list {
		gScore[neigh_data.NeighbourId] = math.MaxFloat32
		fScore[neigh_data.NeighbourId] = math.MaxFloat32
	}

	gScore[startID] = 0
	fScore[startID] = computeDistance(start, stop)

	for len(openSet) > 0 {
		var adjacency_list_node Record
		current := get_lowest_node(fScore, openSet)
		if current == stopID {
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

		openSet = find_index_remove(openSet, current)
		for _, neigh := range adjacency_list_node.Neighbours {
			var tentative_gScore float32 = gScore[current] + neigh.Dist
			if tentative_gScore < gScore[neigh.NeighbourId] {
				cameFrom[neigh.NeighbourId] = current
				gScore[neigh.NeighbourId] = tentative_gScore
				fScore[neigh.NeighbourId] = tentative_gScore + neigh.Dist
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

// func choose_best_truck(fleet Fleet, order Order, call *mongo.Collection) Truck {
// 	/*
// 		Choose the best truck for order knowing it priority
// 	*/
// 	priority := calculate_priority(fleet, order, call)
// 	var chosen_truck Truck
// 	var max_prio float32 = 0
// 	for key, val := range priority {
// 		if val > max_prio {
// 			chosen_truck = key
// 		}
// 	}
// 	return chosen_truck
// }

// func calculate_priority(fleet Fleet, order Order, call *mongo.Collection) map[Truck]float32 {
// 	/*
// 		Calculating priority of taking the task by fleet
// 	*/
// 	var distance float32
// 	var time_to_complete float32

// 	var gauss_time float32
// 	var gauss_cap float32
// 	var priority map[Truck]float32 = make(map[Truck]float32)

// 	for _, truck := range fleet.Trucks {
// 		// Check if truck have enough capacity
// 		if truck.Capacity >= order.Capacity {
// 			// Get the distances
// 			distance_to_get, _ := Astar(call, truck.Location, order.Location_order)
// 			distance_to_put, _ := Astar(call, order.Location_order, order.Location_to_deliver)
// 			// Whole distance
// 			distance = distance_to_get + distance_to_put

// 			// Calculate estimated time for arrive
// 			time_to_complete_assigment := distance / Max_speed_truck

// 			// Time to complete assignment: Time to complete assignment - time now, in time.Duration format
// 			time_for_assigment := time.Until(order.Time_delivery)
// 			// Estimated hours to complete assignment
// 			time_to_complete = float32(time_for_assigment.Hours())

// 			// Calculating priority
// 			gauss_time = gaussian(time_to_complete, time_to_complete_assigment, 1)      // Greater sigma means more spread distribution
// 			gauss_cap = gaussian(float32(truck.Capacity), float32(order.Capacity), 0.6) // Lower sigma means more focused around mean
// 			priority[truck] = gauss_cap + gauss_time
// 		}
// 	}
// 	return priority
// }

// func gaussian(x float32, mean float32, sigma float32) float32 {
// 	// Calculate normal distribution
// 	return float32(1 / (sigma * float32(math.Sqrt(2*3.14))) * float32(math.Exp((float64(-1)/2)*math.Pow((float64(x-mean))/float64(sigma), 2))))
// }
