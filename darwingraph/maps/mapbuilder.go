package maps

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
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

func (m *MapBuilder) SetTileProvider(p *sm.TileProvider) *MapBuilder {
	m.ctx.SetTileProvider(p)
	return m
}

func (m *MapBuilder) SetSize(width, height int) *MapBuilder {
	m.ctx.SetSize(width, height)
	return m
}

func (m *MapBuilder) Render(w io.Writer) error {
	img, err := m.ctx.Render()
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}
