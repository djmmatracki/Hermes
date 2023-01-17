package api

// func find_index_remove(slice []NodeID, curr NodeID) []NodeID {
// 	var index int = 0

// 	for i, v := range slice {
// 		if v == curr {
// 			index = i
// 		}
// 	}
// 	var list_out1 []NodeID = slice[:index]
// 	var list_out2 []NodeID = slice[index+1:]
// 	return append(list_out1, list_out2...)
// }

// func getLowestNode(fs map[NodeID]float32, available_nodes map[NodeID]string) NodeID {
// 	var min_value float32 = math.MaxFloat32
// 	var pos NodeID = -1
// 	for _, v := range available_nodes {
// 		if fs[v] <= min_value {
// 			min_value = fs[v]
// 			pos = v
// 		}
// 	}
// 	return pos
// }

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
