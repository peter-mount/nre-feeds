package maps

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/peter-mount/nre-feeds/darwingraph"
	"image/color"
	"image/png"
	"io"
)

type MapBuilder struct {
	ctx *sm.Context
}

func NewMapBuilder() *MapBuilder {
	m := &MapBuilder{ctx: sm.NewContext()}
	m.ctx.SetCenter(s2.LatLngFromDegrees(54.413, -3.878))
	m.ctx.SetZoom(6)
	return m
}

func (m *MapBuilder) TileProvider(p *sm.TileProvider) *MapBuilder {
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

func (m *MapBuilder) Render(w io.Writer) error {
	img, err := m.ctx.Render()
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}

func (m *MapBuilder) AddObject(o sm.MapObject) *MapBuilder {
	m.ctx.AddObject(o)
	return m
}

func (m *MapBuilder) AppendStation(s *darwingraph.StationNode) *MapBuilder {
	s.ForEachTiploc(func(t *darwingraph.TiplocNode) {
		if t.HasPosition() {
			m.ctx.AddObject(sm.NewCircle(
				s2.LatLngFromDegrees(float64(t.Lat), float64(t.Lon)),
				color.RGBA{R: 0xff, A: 0xff},
				color.RGBA{R: 0xff, A: 0xff},
				500.0,
				5.0,
			))
		}
	})
	return m
}

func (m *MapBuilder) AppendStationEdge(s *darwingraph.StationEdge) *MapBuilder {
	var points []s2.LatLng
	s.ForEachTiploc(func(t *darwingraph.TiplocNode) {
		if t.HasPosition() {
			points = append(points, s2.LatLngFromDegrees(float64(t.Lat), float64(t.Lon)))
		}
	})
	return m.AddObject(sm.NewPath(points, color.RGBA{A: 0xff}, 1))
}
