package main

import (
	"github.com/ch-braun/go-spring-electrical-layout/forces"
	gonumLayout "gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	graph := simple.NewWeightedDirectedGraph(0.0, 0.0)
	for i := 0; i < 10; i++ {
		graph.AddNode(forces.SimpleMassNode{Node: simple.Node(i), Mass: float64(i + 1)})
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == j {
				continue
			}
			graph.SetWeightedEdge(graph.NewWeightedEdge(forces.SimpleMassNode{Node: simple.Node(i), Mass: float64(i + 1)}, forces.SimpleMassNode{Node: simple.Node(j), Mass: float64(j + 1)}, float64(i+j)))
		}
	}

	spring := forces.NewSpringElectricalR2(10.0, 10.0, 3)

	forceStack := forces.NewForceStack(uint64(42), 100, 0.001, 0.4, 0.2)

	forceStack.AddForce(&spring)

	// Make a forces optimizer with the target graph and update function.
	optimizer := gonumLayout.NewOptimizerR2(graph, forceStack.Update)

	for optimizer.Update() {
	}
}
