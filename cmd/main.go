package main

import (
	"Hermes/internal"
	"fmt"
	// "github.com/sirupsen/logrus"
)

func main() {
	// logger := logrus.New()
	// server := createServerFromConfig(logger, ":8000")
	// server.Run()

	maping_nodes := make(map[int64]([]int64))
	maping_latlon := make(map[int64]([]float64))
	maping_nodes, maping_latlon = internal.ListPoints()

	fmt.Println(maping_nodes[99880])
	fmt.Println(maping_latlon[99882])

	fmt.Println(internal.A_star(99880, 7702634198, maping_nodes, maping_latlon))

}

// func createServerFromConfig(logger *logrus.Logger, bind string) *internal.HTTPInstanceAPI {

// 	return internal.NewHTTPInstanceAPI(bind)
// }
