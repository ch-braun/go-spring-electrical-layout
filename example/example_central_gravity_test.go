package main

import (
	"github.com/ch-braun/go-spring-electrical-layout/force"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
	"testing"
)

func TestCentralGravity(t *testing.T) {
	graph := NewExampleGraph()

	gravity := force.NewCentralGravityR2(r2.Vec{X: 0, Y: 0}, 10.0)

	forceStack := force.NewForceStack(42, 10, 1.0, 0.4, 0.2)

	forceStack.AddForce(gravity)

	// Make a force optimizer with the target graph and update function.
	optimizer := layout.NewOptimizerR2(graph, forceStack.Update)

	for optimizer.Update() {
	}
}
