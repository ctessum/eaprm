package eaprm

import (
	"github.com/ctessum/geom"
	"github.com/ctessum/geom/encoding/osm"
	"github.com/ctessum/geom/index/rtree"
)

// Buffer returns a FeatureFunc that sums the count (for point features),
// length (for linear features), or area (for polygonal features) of a
// group of features within the given radius (in degrees)
// of each of the point locations of interest.
func Buffer(radius float64) FeatureFunc {
	return func(features []*osm.GeomTags, points []geom.Point) ([]float64, error) {
		typ, err := osm.DominantType(features)
		if err != nil {
			return nil, err
		}
		switch typ {
		case osm.Point:
			return Normalize(pointBuffer(radius, features, points)), nil
		case osm.Line:
			return Normalize(lineBuffer(radius, features, points)), nil
		case osm.Poly:
			return Normalize(polyBuffer(radius, features, points)), nil
		case osm.Collection:
			return nil, nil
		default:
			panic("invalid type")
		}
	}
}

// Segments is the number of line segments used to represent a circular
// buffer. See github.com/ctessum/geom.Point.Buffer for more information.
var Segments = 20

// pointBuffer counts the number of points in features that are within
// the specified redius of each point. It returns an array of the counts
// for each point.
func pointBuffer(radius float64, features []*osm.GeomTags, points []geom.Point) []float64 {
	featureIndex := rtree.NewTree(25, 50)
	for _, f := range features {
		switch f.Geom.(type) {
		case geom.PointLike:
			p := f.Geom.Points()
			for i := 0; i < f.Geom.Len(); i++ {
				featureIndex.Insert(p())
			}
		}
	}
	o := make([]float64, len(points))
	for i, p := range points {
		buf := p.Buffer(radius, Segments)
		for _, fI := range featureIndex.SearchIntersect(buf.Bounds()) {
			if w := fI.(geom.Point).Within(buf); w == geom.Inside || w == geom.OnEdge {
				o[i] += 1.0
			}
		}
	}
	return o
}

// lineBuffer calculates the length of lines in features that are within
// the specified redius of each point. It returns an array of the lengths
// for each point.
func lineBuffer(radius float64, features []*osm.GeomTags, points []geom.Point) []float64 {
	featureIndex := rtree.NewTree(25, 50)
	for _, f := range features {
		switch f.Geom.(type) {
		case geom.Linear:
			featureIndex.Insert(f.Geom.(geom.Linear))
		}
	}
	o := make([]float64, len(points))
	for i, p := range points {
		buf := p.Buffer(radius, Segments)
		for _, fI := range featureIndex.SearchIntersect(buf.Bounds()) {
			isect := fI.(geom.Linear).Clip(buf)
			if isect != nil {
				o[i] += isect.Length()
			}
		}
	}
	return o
}

// polyBuffer calculates the area of polygons in features that are within
// the specified redius of each point. It returns an array of the areas
// for each point.
func polyBuffer(radius float64, features []*osm.GeomTags, points []geom.Point) []float64 {
	featureIndex := rtree.NewTree(25, 50)
	for _, f := range features {
		switch f.Geom.(type) {
		case geom.Polygonal:
			featureIndex.Insert(f.Geom.(geom.Polygonal))
		}
	}
	o := make([]float64, len(points))
	for i, p := range points {
		buf := p.Buffer(radius, Segments)
		for _, fI := range featureIndex.SearchIntersect(buf.Bounds()) {
			isect := buf.Intersection(fI.(geom.Polygonal))
			if isect != nil {
				o[i] += isect.Area()
			}
		}
	}
	return o
}
