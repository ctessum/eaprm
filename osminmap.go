package eaprm

import (
	"github.com/ctessum/geom"
	"github.com/ctessum/geom/encoding/osm"
	"github.com/gonum/floats"
	"github.com/spatialmodel/inmap"
	"github.com/spatialmodel/inmap/sr"
)

// InMAPPol represents pollutant names in the InMAP SR matrix.
type InMAPPol string

// These are the pollutants accepted in the InMAP SR matrix.
const (
	PrimaryPM25 = "PrimaryPM25"
	PNH4        = "pNH4"
	PNO3        = "pNO3"
	PSO4        = "pSO4"
	SOA         = "SOA"
)

// InMAP returns a FeatureFunc that calculates PM2.5 concentrations at
// each point of interest caused by ground-level emissions from a group
// of features, where each feature emits at a rate equal to 1 for point
// features, the feature length for line features, and the feature area
// for polygon features.
func InMAP(r *sr.Reader, pol InMAPPol) FeatureFunc {
	// TODO: match projections
	return func(features []*osm.GeomTags, points []geom.Point) ([]float64, error) {
		typ, err := dominantType(features)
		if err != nil {
			return nil, err
		}
		var emis []*inmap.EmisRecord
		switch typ {
		case point:
			emis = pointEmis(features, pol)
		case line:
			emis = lineEmis(features, pol)
		case poly:
			emis = polyEmis(features, pol)
		case collection:
			return nil, nil
		default:
			panic("invalid type")
		}
		o := make([]float64, len(r.Geometry()))
		for _, e := range emis {
			conc, err := r.Concentrations(e)
			if err != nil {
				return nil, err
			}
			floats.Add(o, conc.TotalPM25())
		}
		return o, nil
	}
}

func pointEmis(features []*osm.GeomTags, pol InMAPPol) []*inmap.EmisRecord {
	e := make([]*inmap.EmisRecord, 0, len(features))
	for _, f := range features {
		switch f.Geom.(type) {
		case geom.PointLike:
			e = append(e, newEmis(f.Geom, pol, 1))
		}
	}
	return e
}

func lineEmis(features []*osm.GeomTags, pol InMAPPol) []*inmap.EmisRecord {
	e := make([]*inmap.EmisRecord, 0, len(features))
	for _, f := range features {
		switch f.Geom.(type) {
		case geom.Linear:
			l := f.Geom.(geom.Linear)
			e = append(e, newEmis(f.Geom, pol, l.Length()))
		}
	}
	return e
}

func polyEmis(features []*osm.GeomTags, pol InMAPPol) []*inmap.EmisRecord {
	e := make([]*inmap.EmisRecord, 0, len(features))
	for _, f := range features {
		switch f.Geom.(type) {
		case geom.Polygonal:
			p := f.Geom.(geom.Polygonal)
			e = append(e, newEmis(f.Geom, pol, p.Area()))
		}
	}
	return e
}

func newEmis(g geom.Geom, pol InMAPPol, val float64) *inmap.EmisRecord {
	switch pol {
	case PrimaryPM25:
		return &inmap.EmisRecord{
			Geom: g,
			PM25: val,
		}
	case PNH4:
		return &inmap.EmisRecord{
			Geom: g,
			NH3:  val,
		}
	case PNO3:
		return &inmap.EmisRecord{
			Geom: g,
			NOx:  val,
		}
	case PSO4:
		return &inmap.EmisRecord{
			Geom: g,
			SOx:  val,
		}
	case SOA:
		return &inmap.EmisRecord{
			Geom: g,
			VOC:  val,
		}
	default:
		panic("invalid pollutant")
	}
}
