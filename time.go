package eaprm

import (
	"fmt"
	"time"
)

// Hour returns a two dimensional array of 1 or zero values, where
// the first dimension of the array is of length 24 and corresponds
// to the hours of the day, and the second dimension of the array is
// the same length as t. Values in the returned array [i,j] will be
// 1 if the hour of day in t[j] equals i, and zero otherwise.
func Hour(t []time.Time) [][]float64 {
	o := make([][]float64, 24)
	for i := range o {
		o[i] = make([]float64, len(t))
	}
	for j, tt := range t {
		o[tt.Hour()][j] = 1
	}
	return o
}

// Month returns a two dimensional array of 1 or zero values, where
// the first dimension of the array is of length 12 and corresponds
// to the months of the year, and the second dimension of the array is
// the same length as t. Values in the returned array [i,j] will be
// 1 if the month in t[j] equals i, and zero otherwise. January is month 0.
func Month(t []time.Time) [][]float64 {
	o := make([][]float64, 12)
	for i := range o {
		o[i] = make([]float64, len(t))
	}
	for j, tt := range t {
		o[int(tt.Month())-1][j] = 1
	}
	return o
}

// Weekend returns an array of 1 or zero values, with 1 values for
// members of t that are on a weekend (i.e., Saturday or Sunday) and
// zero otherwise.
func Weekend(t []time.Time) []float64 {
	o := make([]float64, len(t))
	for j, tt := range t {
		if w := tt.Weekday(); w == time.Saturday || w == time.Sunday {
			o[j] = 1
		}
	}
	return o
}

// Month returns a two dimensional array of 1 or zero values, where
// the first dimension of the array is of length (end - start), where
// start and end are the first and last years considered, and corresponds
// to the year, and the second dimension of the array is
// the same length as t. Values in the returned array [i,j] will be
// 1 if the year in t[j] equals i, and zero otherwise.
func Year(t []time.Time, start, end int) ([][]float64, error) {
	o := make([][]float64, end-start+1)
	for i := range o {
		o[i] = make([]float64, len(t))
	}
	for j, tt := range t {
		y := tt.Year()
		if !(start <= y && y <= end) {
			return nil, fmt.Errorf("eaprm: year %d out of range", y)
		}
		o[y-start][j] = 1
	}
	return o, nil
}
