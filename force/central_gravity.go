package force

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
)

type CentralGravityR2 struct {

	// Centre represents the centre of this gravitational force.
	Centre r2.Vec

	// GravitationalConstant represents the strength of gravity.
	GravitationalConstant float64
}

func NewCentralGravityR2(
	centre r2.Vec,
	gravitationalConstant float64) *CentralGravityR2 {
	return &CentralGravityR2{
		Centre:                centre,
		GravitationalConstant: gravitationalConstant,
	}
}

func (c *CentralGravityR2) Calculate(g graph.Graph, layoutR2 layout.LayoutR2) map[int64]r2.Vec {
	forces := make(map[int64]r2.Vec, g.Nodes().Len())

	iteratorNodes := g.Nodes()

	for iteratorNodes.Next() {
		node := iteratorNodes.Node()
		nodeId := node.ID()
		nodeVec := layoutR2.Coord2(nodeId)

		vecToCenter := r2.Sub(nodeVec, c.Centre)
		distance2 := r2.Norm2(vecToCenter)

		forces[nodeId] = r2.Scale(c.GravitationalConstant/distance2, vecToCenter)
	}

	return forces
}
