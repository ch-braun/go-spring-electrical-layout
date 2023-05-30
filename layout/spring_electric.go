package layout

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type SpringElectricR2 struct {

	// RandomizerSeed for random node location generation.
	// This enables reproducible results.
	RandomizerSeed uint64

	// OptimalDistance represents parameter K in the paper.
	OptimalDistance float64

	// RepulsionStrength represents the regularization parameter C in the paper.
	RepulsionStrength float64

	// RepulsionExponent represents parameter p (p > 0) for the repulsive force to reduce distortion effects.
	// Default should be 1.
	RepulsionExponent uint

	// Updates represents the number of iterations that should be run.
	Updates int

	// remainingUpdates represents the number of updates that remain for this optimization cycle.
	remainingUpdates int

	// StepSize represents the distance each node may travel for every update.
	StepSize float64

	// CoolingRate represents
	CoolingRate float64

	// StopThreshold represents the stopping criteria for the simulation.
	// Once the movements are small enough, the simulation stops.
	StopThreshold float64
}

func NewSpringElectricR2(
	randomizerSeed uint64,
	optimalDistance float64,
	repulsionStrength float64,
	repulsionExponent uint,
	updates int,
	stepSize float64,
	coolingRate float64,
	stopThreshold float64) SpringElectricR2 {
	return SpringElectricR2{
		RandomizerSeed:    randomizerSeed,
		OptimalDistance:   optimalDistance,
		RepulsionStrength: repulsionStrength,
		RepulsionExponent: repulsionExponent,
		Updates:           updates,
		remainingUpdates:  updates,
		StepSize:          stepSize,
		CoolingRate:       coolingRate,
		StopThreshold:     stopThreshold,
	}
}

// assignInitialCoordinates assigns an initial x and y value to each node.
// These values are randomly chosen from within the value range:
// [-0.5*g.Nodes().Len()*optimalDistance, 0.5*g.Nodes().Len()*optimalDistance]
func assignInitialCoordinates(g graph.Graph, layoutR2 layout.LayoutR2, randomizerSeed uint64, optimalDistance float64) {
	iteratorNodes := g.Nodes()

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

func calculateVecDistance(vec1 r2.Vec, vec2 r2.Vec) float64 {
	return math.Sqrt((vec2.X-vec1.X)*(vec2.X-vec1.X) + (vec2.Y-vec1.Y)*(vec2.Y-vec1.Y))
}

func calculateVecDifference(vec1 r2.Vec, vec2 r2.Vec) r2.Vec {
	return r2.Vec{
		X: vec2.X - vec1.X,
		Y: vec2.Y - vec1.Y,
	}
}

func (s *SpringElectricR2) Update(g graph.Graph, layoutR2 layout.LayoutR2) bool {

	if s.remainingUpdates <= 0 {
		return false
	}

	if !layoutR2.IsInitialized() {
		assignInitialCoordinates(g, layoutR2, s.RandomizerSeed, s.OptimalDistance)
	}

	forces := make(map[int64]r2.Vec, g.Nodes().Len())
	layoutNew := make(map[int64]r2.Vec, g.Nodes().Len())
	stepSize := s.StepSize * math.Pow(s.CoolingRate, float64(s.Updates-s.remainingUpdates))
	s.remainingUpdates -= 1

	iteratorNodesOuter, iteratorNodesInner := g.Nodes(), g.Nodes()

	for iteratorNodesOuter.Next() {
		n1 := iteratorNodesOuter.Node()
		id1 := n1.ID()
		forceCombined := r2.Vec{X: 0, Y: 0}
		for iteratorNodesInner.Next() {
			id2 := iteratorNodesInner.Node().ID()

			if id1 == id2 {
				continue
			}
			vec1 := layoutR2.Coord2(id1)
			vec2 := layoutR2.Coord2(id2)

			distance := calculateVecDistance(vec1, vec2)
			vectorDifference := calculateVecDifference(vec1, vec2)

			forceRepulsive := -s.RepulsionStrength * math.Pow(s.OptimalDistance, 1.0+float64(s.RepulsionExponent)) /
				math.Pow(distance, float64(s.RepulsionExponent))
			forceAttractive := (distance * distance) / s.OptimalDistance

			forceCombined.X += (forceRepulsive + forceAttractive) / distance * vectorDifference.X
			forceCombined.Y += (forceRepulsive + forceAttractive) / distance * vectorDifference.Y
		}
		iteratorNodesInner.Reset()

		forces[id1] = forceCombined

		var nodeMass float64
		wn, ok := n1.(SimpleMassNode)
		if ok {
			nodeMass = wn.Mass
		} else {
			nodeMass = 1.0
		}

		layoutNew[id1] = r2.Vec{
			X: layoutR2.Coord2(id1).X + stepSize*forceCombined.X/nodeMass,
			Y: layoutR2.Coord2(id1).Y + stepSize*forceCombined.Y/nodeMass,
		}
	}
	accumulativeDelta := 0.0
	for nid, vec := range layoutNew {
		accumulativeDelta += calculateVecDistance(layoutR2.Coord2(nid), vec)
		layoutR2.SetCoord2(nid, vec)
	}

	return accumulativeDelta > s.StopThreshold
}
