package main

import (
	"github.com/ch-braun/go-spring-electrical-layout/force"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
	"log"
	"testing"
)

func TestSpringElectricGravity(t *testing.T) {
	graph := NewExampleGraph()

	spring := force.NewSpringElectricalR2(10.0, 10.0, 3, 1)

	gravity := force.NewCentralGravityR2(r2.Vec{X: 0, Y: 0}, 100.0)

	forceStack := force.NewForceStack(uint64(42), 100, 0.001, 0.4, 0.2)

	forceStack.AddForce(spring)
	forceStack.AddForce(gravity)

	// Make a force optimizer with the target graph and update function.
	optimizer := layout.NewOptimizerR2(graph, forceStack.Update)

	for optimizer.Update() {
	}

	iter := graph.Nodes()

	for iter.Next() {
		log.Printf("%v: %v", iter.Node().ID(), optimizer.LayoutNodeR2(iter.Node().ID()).Coord2)
	}
}
