package eaprm

import (
	"os"
	"testing"

	"github.com/ctessum/geom"
	"github.com/ctessum/geom/encoding/osm"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

const (
	Xmin = -157.983108
	Xmax = -157.620559
	Ymin = 21.246082
	Ymax = 21.421331
	nx   = 100
	ny   = 100

	mPerDegree = 111.111 * 1000 // meters per degree latitude
)

type GridXYZ struct {
	*mat.Dense
	Xmin, Xmax, Ymin, Ymax float64
}

func (g GridXYZ) Z(c, r int) float64 {
	return g.Dense.At(r, c)
}

func (g GridXYZ) X(c int) float64 {
	cols, _ := g.Dims()
	return g.Xmin + (g.Xmax-g.Xmin)/float64(cols)*float64(c)
}

func (g GridXYZ) Y(r int) float64 {
	_, rows := g.Dims()
	return g.Ymin + (g.Ymax-g.Ymin)/float64(rows)*float64(r)
}

func TestBuffer(t *testing.T) {
	points := testImagePoints(Xmin, Xmax, Ymin, Ymax, nx, ny)

	const (
		figWidth  = 20 * vg.Centimeter
		figHeight = 20 * vg.Centimeter
		figRows   = 3
		figCols   = 3
	)
	img := vgimg.NewWith(vgimg.UseWH(figWidth, figHeight), vgimg.UseDPI(96))
	dc := draw.New(img)
	tiles := draw.Tiles{
		Rows: figRows,
		Cols: figCols,
	}

	f, err := os.Open(os.ExpandEnv("${GOPATH}/src/github.com/ctessum/geom/encoding/osm/testdata/honolulu_hawaii.osm.pbf"))
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	for i, keyVal := range [][]string{{"highway", "primary"}, {"building", "residential"}, {"natural", "tree"}} {
		for j, bufferLen := range []float64{50, 250, 1000} { // meters
			var tags *osm.Data
			tags, err = osm.ExtractTag(f, keyVal[0], keyVal[1])
			if err != nil {
				t.Fatal(err)
			}
			var geomTags []*osm.GeomTags
			geomTags, err = tags.Geom()
			if err != nil {
				t.Fatal(err)
			}
			var vals []float64
			vals, err = Buffer(bufferLen/mPerDegree)(geomTags, points)
			if err != nil {
				t.Fatal(err)
			}
			g := GridXYZ{
				Xmin:  Xmin,
				Xmax:  Xmax,
				Ymin:  Ymin,
				Ymax:  Ymax,
				Dense: mat.NewDense(nx, ny, vals),
			}
			h := plotter.NewHeatMap(g, moreland.ExtendedBlackBody().Palette(100))

			var p *plot.Plot
			p, err = plot.New()
			if err != nil {
				t.Error(err)
			}
			p.Add(h)
			p.Draw(tiles.At(dc, j, i))
		}
	}
	w, err := os.Create("testdata/buffer.png")
	if err != nil {
		t.Error(err)
	}
	cc := vgimg.PngCanvas{Canvas: img}
	if _, err := cc.WriteTo(w); err != nil {
		t.Error(err)
	}
}

// testImagePoints returns the points at which to make calculations
// for a test image.
func testImagePoints(xstart, xend, ystart, yend float64, nx, ny int) []geom.Point {
	dx := (xend - xstart) / float64(nx)
	dy := (yend - ystart) / float64(ny)
	g := make([]geom.Point, nx*ny)
	var i int
	y0 := ystart + dy/2
	x0 := xstart + dx/2
	for yy := y0; yy <= yend; yy += dy {
		for xx := x0; xx <= xend; xx += dx {
			g[i] = geom.Point{X: xx, Y: yy}
			i++
		}
	}
	return g
}
