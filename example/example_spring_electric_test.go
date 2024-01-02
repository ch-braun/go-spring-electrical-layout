package main

import (
	"github.com/ch-braun/go-spring-electrical-layout/force"
	"gonum.org/v1/gonum/graph/layout"
	"testing"
)

func TestSpringElectric(t *testing.T) {
	graph := NewExampleGraph()

	spring := force.NewSpringElectricalR2(10.0, 10.0, 3, 1)

	forceStack := force.NewForceStack(uint64(42), 100, 0.001, 0.4, 0.2)

	forceStack.AddForce(spring)

	// Make a force optimizer with the target graph and update function.
	optimizer := layout.NewOptimizerR2(graph, forceStack.Update)

	for optimizer.Update() {
	}
}
