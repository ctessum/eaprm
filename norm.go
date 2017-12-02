package eaprm

import (
	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// Normalize transforms d so that its mean is zero and its
// standard deviation is one.
func Normalize(d []float64) []float64 {
	floats.AddConst(-stat.Mean(d, nil), d)
	floats.Scale(1/stat.StdDev(d, nil), d)
	return d
}
