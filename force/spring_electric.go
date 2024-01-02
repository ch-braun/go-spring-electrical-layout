package force

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/layout"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type SpringElectricalR2 struct {

	// OptimalDistance represents parameter K in the paper.
	OptimalDistance float64

	// RepulsionStrength represents the regularization parameter C in the paper.
	RepulsionStrength float64

	// RepulsionExponent represents parameter p (p > 0) for the repulsive force to reduce distortion effects.
	// Default should be 1. This exponent is part of the denominator. Increasing it reduces the applied force.
	RepulsionExponent uint

	// AttractionExponent is an added parameter for the attractive spring force to regulate its effect.
	// Default should be 1.0. This exponent is part of the numerator. Increasing it increases the applied force.
	AttractionExponent float64

	// Epsilon is the threshold for small forces. Every force below this value will be omitted.
	Epsilon float64
}

func NewSpringElectricalR2(
	optimalDistance float64,
	repulsionStrength float64,
	repulsionExponent uint,
	attractionExponent float64,
	epsilon float64) *SpringElectricalR2 {
	return &SpringElectricalR2{
		OptimalDistance:    optimalDistance,
		RepulsionStrength:  repulsionStrength,
		RepulsionExponent:  repulsionExponent,
		AttractionExponent: attractionExponent,
		Epsilon:            epsilon,
	}
}

func (s *SpringElectricalR2) Calculate(g graph.Graph, layoutR2 layout.LayoutR2) map[int64]r2.Vec {

	forces := make(map[int64]r2.Vec, g.Nodes().Len())

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

			distance := r2.Norm(r2.Sub(vec1, vec2))
			vectorDifference := r2.Sub(vec1, vec2)

			forceRepulsive := -s.RepulsionStrength * math.Pow(s.OptimalDistance, 1.0+float64(s.RepulsionExponent)) /
				math.Pow(distance, float64(s.RepulsionExponent))

			var forceAttractive float64
			if e := g.Edge(id1, id2); e != nil {
				edgeWeight := 1.
				if we, ok := e.(graph.WeightedEdge); ok {
					edgeWeight = we.Weight()
				}
				forceAttractive = edgeWeight * math.Pow(distance, s.AttractionExponent) / s.OptimalDistance
			} else {
				forceAttractive = 0.
			}

			// Apply epsilon gate
			if math.Abs(forceRepulsive) < s.Epsilon {
				forceRepulsive = 0.
			}

			if math.Abs(forceAttractive) < s.Epsilon {
				forceAttractive = 0.
			}

			vectorDifference = r2.Scale((forceRepulsive+forceAttractive)/distance, vectorDifference)

			forceCombined.X += vectorDifference.X
			forceCombined.Y += vectorDifference.Y
		}
		iteratorNodesInner.Reset()

		forces[id1] = forceCombined
	}

	return forces
}
