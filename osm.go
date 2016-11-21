package eaprm

import (
	"fmt"
	"io"

	"github.com/ctessum/geom"
	"github.com/ctessum/geom/encoding/osm"
)

// OSM extracts information from an OpenStreetMap database.
type OSM struct {
	f io.ReadSeeker
}

// NewOSM creates a new OSM object where f is a connection to an
// OpenStreetMap database.
func NewOSM(f io.ReadSeeker) *OSM {
	return &OSM{
		f: f,
	}
}

// FeatureFunc is a class of functions that takes and array of
// OpenStreetMap features and and array of points to relate
// them to and returns a relation between all of the features
// and each of the points.
type FeatureFunc func([]*osm.GeomTags, []geom.Point) ([]float64, error)

type geomType int

const (
	point geomType = iota
	line
	poly
	collection
)

func dominantType(gt []*osm.GeomTags) (geomType, error) {
	var points, lines, polys, collections int
	for _, g := range gt {
		switch g.Geom.(type) {
		case geom.PointLike:
			points++
		case geom.Linear:
			lines++
		case geom.Polygonal:
			polys++
		case geom.GeometryCollection:
			collections++
		default:
			return -1, fmt.Errorf("invalid geometry type %#v", g.Geom)
		}
	}
	if points >= lines && points >= polys && points >= collections {
		return point, nil
	}
	if lines > points && lines >= polys && lines >= collections {
		return line, nil
	}
	if polys > points && polys > lines && polys >= collections {
		return poly, nil
	}
	return collection, nil
}
