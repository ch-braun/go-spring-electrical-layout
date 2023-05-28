package layout

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
	"log"
	"math"
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
}

func nodeDistance(id1 int64, id2 int64, layoutR2 layout.LayoutR2) float64 {
	vec1 := layoutR2.Coord2(id1)
	vec2 := layoutR2.Coord2(id2)

	return math.Sqrt((vec2.X-vec1.X)*(vec2.X-vec1.X) + (vec2.Y-vec1.Y)*(vec2.Y-vec1.Y))
}

func (s SpringR2) Update(g graph.Graph, layoutR2 layout.LayoutR2) bool {
	if s.Updates <= 0 {
		return false
	}
	s.Updates--

	if !layoutR2.IsInitialized() {
		assignInitialCoordinates(g, layoutR2, s.RandomizeLocations, s.RandomizerSeed, s.OptimalDistance)
	}

	iteratorNodesOuter, iteratorNodesInner := g.Nodes(), g.Nodes()

	for iteratorNodesOuter.Next() {
		n1 := iteratorNodesOuter.Node()
		for iteratorNodesInner.Next() {
			n2 := iteratorNodesInner.Node()

			if n1.ID() == n2.ID() {
				continue
			}
			distance := nodeDistance(n1.ID(), n2.ID(), layoutR2)
			forceRepulsive := -s.RepulsionStrength * (s.OptimalDistance * s.OptimalDistance) /
				math.Abs(distance)

			forceAttractive := (distance * distance) / s.OptimalDistance

			forceCombinedX := (forceRepulsive + forceAttractive) / (math.Abs(distance)) * (layoutR2.Coord2(n2.ID()).X - layoutR2.Coord2(n1.ID()).X)
			forceCombinedY := (forceRepulsive + forceAttractive) / (math.Abs(distance)) * (layoutR2.Coord2(n2.ID()).Y - layoutR2.Coord2(n1.ID()).Y)

			log.Printf("%v %v", forceCombinedX, forceCombinedY)
		}
	}

	return true
}
