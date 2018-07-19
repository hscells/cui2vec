package cui2vec

import (
	"math"
	"gonum.org/v1/gonum/floats"
)

// softmax normalises a slice of concepts.
func softmax(z []Concept) []Concept {
	zExp := make([]float64, len(z))
	for i := range z {
		zExp[i] = math.Exp(z[i].Value)
	}
	sumZExp := floats.Sum(zExp)
	softmax := make([]Concept, len(z))
	for i := range zExp {
		softmax[i] = Concept{CUI: z[i].CUI, Value: zExp[i] / sumZExp}
	}
	return softmax
}
