package layout

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
)

type SpringR2 struct {
	// RandomizeLocations indicates whether the nodes should
	// initially be placed at random locations or at (0,0).
	RandomizeLocations bool

	// RandomizerSeed for random node location generation.
	// This enables reproducible results.
	RandomizerSeed uint64

	// OptimalDistance represents parameter K in the paper.
	OptimalDistance float64

	// RepulsionStrength represents the regularization parameter C in the paper.
	RepulsionStrength float64

	// Updates represents the number of iterations that should be run.
	Updates int
}

// assignInitialCoordinates assigns an initial x and y value to each node.
// These values are either (0,0) or randomly chosen from within the value range:
// [-0.5*g.Nodes().Len()*optimalDistance, 0.5*g.Nodes().Len()*optimalDistance]
func assignInitialCoordinates(g graph.Graph, layoutR2 layout.LayoutR2, randomizeLocations bool, randomizerSeed uint64, optimalDistance float64) {
	iteratorNodes := g.Nodes()

	// If locations shall not be assigned randomly, take this shortcut.
	if !randomizeLocations {
		for iteratorNodes.Next() {
			n := iteratorNodes.Node()
			vecR2 := r2.Vec{X: 0, Y: 0}
			layoutR2.SetCoord2(n.ID(), vecR2)
		}
		iteratorNodes.Reset()
		return
	}

	// Initialize pseudo-random number generator based on provided seed value.
	prng := rand.New(rand.NewSource(randomizerSeed))

	randomValueRange := float64(g.Nodes().Len()) * optimalDistance

	// Generate for x and y a random float in range [0,1) and
	// scale it by the maximum random range g.Nodes().Len() * optimalDistance.
	// Center the values by subtracting 0.5 * (g.Nodes().Len() * optimalDistance).
	for iteratorNodes.Next() {
		n := iteratorNodes.Node()
		vecR2 := r2.Vec{
			X: prng.Float64()*randomValueRange - 0.5*randomValueRange,
			Y: prng.Float64()*randomValueRange - 0.5*randomValueRange,
		}
		layoutR2.SetCoord2(n.ID(), vecR2)
	}
	iteratorNodes.Reset()
}

func (spring SpringR2) Update(g graph.Graph, layoutR2 layout.LayoutR2) bool {
	if spring.Updates <= 0 {
		return false
	}
	spring.Updates--

	if !layoutR2.IsInitialized() {
		assignInitialCoordinates(g, layoutR2, spring.RandomizeLocations, spring.RandomizerSeed, spring.OptimalDistance)
	}

	return true
}
