package internal

import (
	"math"
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

	/*
		Opcja 1:
			- Pobieranie całej caly danych na raz przy działaniu programu:
				- start
				- stop

		Opcje 2:
		Pobieranie danych w czasie działania programu
			start, stop zostaje
			fScore dla node_id -> neigbour_id -> dist (cale heuritic_a_star tak)
			adjecency_list -> pobieranie w postaci neigbours (z bazy) odwolanie do elementów neigbour_id


	*/

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

	// Makes all cheapest costs to infinity and heuristics update
	for k := range adjacency_list {
		gScore[k] = math.MaxFloat64
		fScore[k] = heuritic_a_star(nodes_position[k], nodes_position[stop_id])
	}

	// Update the cost to to first location to 0 (we are here) and fScore to distance now
	gScore[start_id] = 0
	fScore[start_id] = heuritic_a_star(nodes_position[start_id], nodes_position[stop_id]) // FROM BASE

	// While we can visit neighbour node
	for len(openSet) > 0 {
		// From fScore we take node with the lowest value of: road to node + distance from node to destination
		current := get_lowest_node(fScore, openSet)
		if current == stop_id {
			return reconstruct_path(cameFrom, current, stop_id, start_id)
		}

		// Delete actual visiting node
		openSet = find_index_remove(openSet, current)
		for _, v := range adjacency_list[current] {
			var tentative_gScore float64 = gScore[current] + heuritic_a_star(nodes_position[current], nodes_position[v]) // FROM BASE
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
