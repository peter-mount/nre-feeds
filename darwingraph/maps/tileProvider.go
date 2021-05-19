package maps

import (
	sm "github.com/flopp/go-staticmaps"
)

type tileProvider struct {
	Name      string
	Title     string
	Generator func() *sm.TileProvider
}

// List of supported tile providers. These require no API keys but you must check them first
// to see what licensing usages apply.
var tileProviders = []tileProvider{
	// map.lu is our own map server
	{"maplu", "map.lu osm 2018", NewTileProviderMapLU},
	// NE2 is public domain & we have that map tiled.
	// Disabled currently as static maps only supports slippy maps
	//{"ne2", "map.lu Natural Earth 2 10m", NewTileProviderNaturalEarth2},

	// Provided by staticMaps, check these on their T&C's before using
	{"osm", "openstreetmap.org", sm.NewTileProviderOpenStreetMaps},
	{"wikimedia", "Wikimedia", sm.NewTileProviderWikimedia},
	{"stamen-toner", "Stamen Terrain", sm.NewTileProviderStamenTerrain},
}

func NewTileProviderMapLU() *sm.TileProvider {
	t := new(sm.TileProvider)
	t.Name = "maplu"
	t.Attribution = "Maps & Data (c) map.lu, openstreetmap.org, Network Rail, NRE & contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://s%[1]s.map.lu/osm201810/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

func NewTileProviderNaturalEarth2() *sm.TileProvider {
	t := new(sm.TileProvider)
	t.Name = "ne2"
	t.Attribution = "Maps & Data (c) map.lu, Natural Earth, Network Rail, NRE & contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://s%[1]s.map.lu/NaturalEarth2_10m/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}
