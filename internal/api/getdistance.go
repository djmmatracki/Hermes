package api

func (a *InstanceAPI) getDistance() {
	// startNode, err := findNearestNode(collection, start)
	// if err != nil {
	// 	return 0, err
	// }

	// stopNode, err := findNearestNode(collection, stop)
	// if err != nil {
	// 	return 0, err
	// }
	// startID := startNode.NodeId
	// stopID := stopNode.NodeId

	// openSet := map[NodeID]string{startID: ""}
	// cameFrom := make(map[NodeID]NodeID)

	// // gScore is cost of the cheapest path from from start to currently known
	// gScore := make(map[NodeID]float32)
	// // fScore is current best guess as how we can get to finish
	// fScore := make(map[NodeID]float32)

	// startNodeNeighbours := startNode.Neighbours
	// for _, neigh_data := range startNodeNeighbours {
	// 	gScore[neigh_data.NeighbourId] = math.MaxFloat32
	// 	fScore[neigh_data.NeighbourId] = math.MaxFloat32
	// }

	// gScore[startID] = 0
	// fScore[startID] = computeDistance(start, stop)

	// for len(openSet) > 0 {
	// 	var adjacency_list_node Record
	// 	current := getLowestNode(fScore, openSet)
	// 	if current == stopID {
	// 		return gScore[current], err
	// 	}

	// 	result := collection.FindOne(
	// 		context.TODO(),
	// 		bson.D{{Key: "node_id", Value: current}},
	// 	)
	// 	err := result.Decode(&adjacency_list_node)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	delete(openSet, current)
	// 	for _, neigh := range adjacency_list_node.Neighbours {
	// 		var tentative_gScore float32 = gScore[current] + neigh.Dist
	// 		if tentative_gScore < gScore[neigh.NeighbourId] {
	// 			cameFrom[neigh.NeighbourId] = current
	// 			gScore[neigh.NeighbourId] = tentative_gScore
	// 			fScore[neigh.NeighbourId] = tentative_gScore + neigh.Dist
	// 			if _, ok := openSet[neigh.NeighbourId]; !ok {
	// 				openSet[neigh.NeighbourId] = ""
	// 			}
	// 		}
	// 	}

	// }
	// return float32(1), nil
}

// func (a *InstanceAPI) findNearestNode(location providers.Location) (*providers., error) {
// 	var location providers.Location
// }
