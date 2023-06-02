package main

import (
	"github.com/ch-braun/go-spring-electrical-layout/force"
	"gonum.org/v1/gonum/graph/simple"
)

func NewExampleGraph() *simple.WeightedDirectedGraph {
	graph := simple.NewWeightedDirectedGraph(0.0, 0.0)
	for i := 0; i < 10; i++ {
		graph.AddNode(force.SimpleMassNode{Node: simple.Node(i), Mass: float64(i + 1)})
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == j {
				continue
			}
			graph.SetWeightedEdge(graph.NewWeightedEdge(force.SimpleMassNode{Node: simple.Node(i), Mass: float64(i + 1)}, force.SimpleMassNode{Node: simple.Node(j), Mass: float64(j + 1)}, float64(i+j)))
		}
	}
	return graph
}
