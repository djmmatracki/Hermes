package main

import (
	"context"
	"fmt"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

type Node struct {
	Latitude       float64
	Longitude      float64
	ConnectedNodes map[int64]string
}

func main() {
	association := make(map[int64]*Node)
	file, err := os.Open("./poland.osm.pbf")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := osmpbf.New(context.Background(), file, 20)
	defer scanner.Close()

	for scanner.Scan() {
		switch o := scanner.Object().(type) {
		case *osm.Node:
			if node, ok := association[int64(o.ID)]; !ok {
				association[int64(o.ID)] = &Node{
					Latitude:       o.Lat,
					Longitude:      o.Lon,
					ConnectedNodes: make(map[int64]string, 0),
				}
			} else {
				node.Latitude = o.Lat
				node.Longitude = o.Lon
			}
		case *osm.Way:
			for index1, node1 := range o.Nodes {
				for index2, node2 := range o.Nodes {
					if index1 != index2 {
						if n1, ok := association[int64(node1.ID)]; ok {
							n1.ConnectedNodes[int64(node2.ID)] = ""
						} else {
							association[int64(node1.ID)] = &Node{
								ConnectedNodes: map[int64]string{
									int64(node2.ID): "",
								},
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("Done!!!")
}
