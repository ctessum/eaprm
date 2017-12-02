package eaprm

import (
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
