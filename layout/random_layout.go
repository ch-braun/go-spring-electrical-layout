package layout

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
)

// AssignRandomCoordinates assigns an x and y value to each node.
// These values are randomly chosen from within the value range:
// [-5.0*g.Nodes().Len(), 5.0*g.Nodes().Len()]
func AssignRandomCoordinates(g graph.Graph, layoutR2 layout.LayoutR2, randomizerSeed uint64) {
	iteratorNodes := g.Nodes()

	// Initialize pseudo-random number generator based on provided seed value.
	prng := rand.New(rand.NewSource(randomizerSeed))

	randomValueRange := float64(g.Nodes().Len()) * 10.0

	// Generate for x and y a random float in range [0,1) and
	// scale it by the maximum random range 10.0 * g.Nodes().Len().
	// Center the values by subtracting 0.5 * maximum random range.
	for iteratorNodes.Next() {
		n := iteratorNodes.Node()
		vecR2 := r2.Vec{
			X: prng.Float64()*randomValueRange - 0.5*randomValueRange,
			Y: prng.Float64()*randomValueRange - 0.5*randomValueRange,
		}
		layoutR2.SetCoord2(n.ID(), vecR2)
	}
}
