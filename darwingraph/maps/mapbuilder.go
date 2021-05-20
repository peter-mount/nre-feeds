package maps

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/peter-mount/nre-feeds/darwingraph"
	"image/color"
	"image/png"
	"io"
)

type MapBuilder struct {
	ctx *sm.Context
}

type MapTask func(builder *MapBuilder)

func NewMapBuilder() *MapBuilder {
	m := &MapBuilder{ctx: sm.NewContext()}
	m.ctx.SetCenter(s2.LatLngFromDegrees(54.413, -3.878))
	m.ctx.SetZoom(6)
	return m
}

func (m *MapBuilder) TileProvider(p *sm.TileProvider) *MapBuilder {
	p.IgnoreNotFound = true
	m.ctx.SetTileProvider(p)
	return m
}

func (m *MapBuilder) Center(lon, lat float64) *MapBuilder {
	m.ctx.SetCenter(s2.LatLngFromDegrees(lat, lon))
	return m
}

func (m *MapBuilder) Size(width, height int) *MapBuilder {
	m.ctx.SetSize(width, height)
	return m
}

func (m *MapBuilder) Zoom(z int) *MapBuilder {
	m.ctx.SetZoom(z)
	return m
}

func (m *MapBuilder) Run(t MapTask) *MapBuilder {
	if t != nil {
		t(m)
	}
	return m
}

func (m *MapBuilder) Render(w io.Writer) error {
	img, err := m.ctx.Render()
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}

type wrapperPath struct {
	p *sm.Path
}

func (o wrapperPath) Bounds() s2.Rect { return o.p.Bounds() }
func (o wrapperPath) ExtraMarginPixels() (float64, float64, float64, float64) {
	return o.p.ExtraMarginPixels()
}
func (o wrapperPath) Draw(gc *gg.Context, trans *sm.Transformer) {
	if len(o.p.Positions) <= 1 {
		return
	}

	first := true

	gc.ClearPath()
	gc.SetLineWidth(o.p.Weight)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	for _, ll := range o.p.Positions {
		x, y := trans.LatLngToXY(ll)
		if x < 0 || y < 0 || int(x) >= gc.Width() || int(y) >= gc.Height() {
			first = true
		} else if first {
			gc.MoveTo(x, y)
			first = false
		} else {
			gc.LineTo(x, y)
		}
	}
	gc.SetColor(o.p.Color)
	gc.Stroke()
}

func (m *MapBuilder) AddObject(o sm.MapObject) *MapBuilder {
	if o != nil {
		if p, ok := o.(*sm.Path); ok {
			m.ctx.AddObject(&wrapperPath{p: p})
		} else {
			m.ctx.AddObject(o)
		}
	}
	return m
}

func (m *MapBuilder) ForEachStationNode(d *darwingraph.DarwinGraph, f func(*MapBuilder, *darwingraph.StationNode)) *MapBuilder {
	d.ForEachStationNode(func(n *darwingraph.StationNode) {
		f(m, n)
	})
	return m
}

func (m *MapBuilder) AppendStation(s *darwingraph.StationNode) *MapBuilder {
	s.ForEachTiploc(func(t *darwingraph.TiplocNode) {
		if t.HasPosition() {
			m.AddObject(sm.NewCircle(
				s2.LatLngFromDegrees(float64(t.Lat), float64(t.Lon)),
				color.RGBA{R: 0xff, A: 0xff},
				color.RGBA{R: 0xff, A: 0xff},
				100.0,
				5.0,
			))
		}
	})
	return m
}

func (m *MapBuilder) ForEachStationEdge(d *darwingraph.DarwinGraph, f func(*MapBuilder, *darwingraph.StationEdge)) *MapBuilder {
	d.ForEachStationEdge(func(e *darwingraph.StationEdge) {
		f(m, e)
	})
	return m
}

func (m *MapBuilder) AppendStationEdge(s *darwingraph.StationEdge) *MapBuilder {
	var a s2.LatLng    // Previous point
	var ar []s2.LatLng // Slice of points forming this segment
	s.ForEachTiploc(func(t *darwingraph.TiplocNode) {
		if t.HasPosition() {
			ll := s2.LatLngFromDegrees(float64(t.Lat), float64(t.Lon))

			// Ignore lines longer than .3 degrees - cuts out some major errors in the data
			if len(ar) > 0 && ll.Distance(a).Degrees() > 0.2 {
				m.appendEdge(ar) // Append existing segment up to this point
				ar = nil         // Start a new slice
			}

			ar = append(ar, ll) // Append new point to segment
			a = ll              // copy for distance check against next point
		}
	})
	m.appendEdge(ar)
	return m
}

func (m *MapBuilder) appendEdge(p []s2.LatLng) {
	if len(p) > 1 {
		m.AddObject(sm.NewPath(p, color.RGBA{A: 0xff}, 1))
	}
}
