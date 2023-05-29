package main

import (
	"github.com/ch-braun/go-spring-electrical-layout/layout"
	gonumLayout "gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	graph := simple.NewWeightedDirectedGraph(0.0, 0.0)
	for i := 0; i < 10; i++ {
		graph.AddNode(simple.Node(i))
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == j {
				continue
			}
			graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(i), simple.Node(j), float64(i+j)))
		}
	}

	spring := layout.SpringElectricR2{
		RandomizeLocations: true,
		RandomizerSeed:     uint64(42),
		OptimalDistance:    10.0,
		RepulsionStrength:  10.0,
		RepulsionExponent:  3,
		Updates:            10,
		RemainingUpdates:   10,
		StepSize:           0.001,
	}

	// Make a layout optimizer with the target graph and update function.
	optimizer := gonumLayout.NewOptimizerR2(graph, spring.Update)

	for optimizer.Update() {
	}
}
