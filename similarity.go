package cui2vec

import (
	"math"
	"fmt"
	"gonum.org/v1/gonum/floats"
)

// product returns a vector of element-wise products of two input vectors.
func product(x, y []float64) ([]float64, error) {
	if len(x) != len(y) {
		return nil, fmt.Errorf("x and y have unequal lengths: %d / %d", len(x), len(y))
	}

	p := make([]float64, len(x))
	for i := range x {
		p[i] = x[i] * y[i]
	}
	return p, nil
}

// dotProduct returns the dot product of two vectors.
func dotProduct(x, y []float64) (float64, error) {
	p, err := product(x, y)
	if err != nil {
		return 0, err
	}
	return floats.Sum(p), nil
}

// norm returns the vector norm.  Use pow = 2.0 for Euclidean.
func norm(x []float64, pow float64) float64 {
	s := 0.0

	for _, v := range x {
		s += math.Pow(v, pow)
	}

	return math.Pow(s, 1/pow)
}

// Cosine returns the cosine similarity between two vectors.
func Cosine(x, y []float64) (float64, error) {
	d, err := dotProduct(x, y)
	if err != nil {
		return 0, err
	}

	xNorm := norm(x, 2.0)
	yNorm := norm(y, 2.0)

	return d / (xNorm * yNorm), nil
}
