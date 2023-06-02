package forces

import (
	"github.com/ch-braun/go-spring-electrical-layout/layout"
	"gonum.org/v1/gonum/graph"
	gn_layout "gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type Force interface {
	Calculate(g graph.Graph, layoutR2 gn_layout.LayoutR2) map[int64]r2.Vec
}

type ForceStack struct {
	forces []Force

	// RandomizerSeed for random node location generation.
	// This enables reproducible results.
	RandomizerSeed uint64

	// Updates represents the number of iterations that should be run.
	Updates uint

	// remainingUpdates represents the number of updates that remain for this optimization cycle.
	remainingUpdates uint

	// StepSize represents the distance each node may travel for every update.
	StepSize float64

	// CoolingRate represents the decrease of the step size per update.
	CoolingRate float64

	// StopThreshold represents the stopping criteria for the simulation.
	// Once the movements are small enough, the simulation stops.
	StopThreshold float64
}

func NewForceStack(randomizerSeed uint64, updates uint, stepSize float64, coolingRate float64, stopThreshold float64) ForceStack {
	return ForceStack{
		forces:           make([]Force, 0),
		RandomizerSeed:   randomizerSeed,
		Updates:          updates,
		remainingUpdates: updates,
		StepSize:         stepSize,
		CoolingRate:      coolingRate,
		StopThreshold:    stopThreshold,
	}
}

func (s *ForceStack) AddForce(force Force) {
	if s.forces == nil {
		s.forces = make([]Force, 0, 1)
	}
	s.forces = append(s.forces, force)
}

func (s *ForceStack) Update(g graph.Graph, layoutR2 gn_layout.LayoutR2) bool {
	if s.remainingUpdates <= 0 {
		return false
	}

	if !layoutR2.IsInitialized() {
		layout.AssignRandomCoordinates(g, layoutR2, s.RandomizerSeed)
	}

	if s.forces == nil {
		return false
	}

	forcesCombined := make(map[int64]r2.Vec, g.Nodes().Len())

	for _, force := range s.forces {
		fc := force.Calculate(g, layoutR2)

		if len(forcesCombined) == 0 {
			forcesCombined = fc
			continue
		}

		for nodeId, forceCurrentNode := range fc {
			forcesCombined[nodeId] = r2.Add(forcesCombined[nodeId], forceCurrentNode)
		}
	}

	layoutNew := make(map[int64]r2.Vec, g.Nodes().Len())
	stepSize := s.StepSize * math.Pow(1.0-s.CoolingRate, float64(s.Updates-s.remainingUpdates))
	s.remainingUpdates -= 1

	iteratorNodes := g.Nodes()
	for iteratorNodes.Next() {
		node := iteratorNodes.Node()
		nodeId := node.ID()
		var nodeMass float64
		wn, ok := node.(SimpleMassNode)
		if ok {
			nodeMass = wn.Mass
		} else {
			nodeMass = 1.0
		}

		layoutNew[nodeId] = r2.Add(layoutR2.Coord2(nodeId), r2.Scale(stepSize/nodeMass, forcesCombined[nodeId]))
	}

	accumulativeDelta := 0.0
	for nid, vec := range layoutNew {
		accumulativeDelta += calculateVecDistance(layoutR2.Coord2(nid), vec)
		layoutR2.SetCoord2(nid, vec)
	}

	return accumulativeDelta > s.StopThreshold
}
