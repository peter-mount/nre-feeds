package darwinref

import (
	"encoding/xml"
	"time"
)

// Defines a location, i.e. Station or passing point
type Location struct {
	XMLName xml.Name `json:"-" xml:"LocationRef"`
	// Tiploc of this location
	Tiploc string `json:"tpl" xml:"tpl,attr"`
	// CRS of this station, "" for none
	Crs string `json:"crs,omitempty" xml:"crs,attr,omitempty"`
	// TOC who manages this station
	Toc string `json:"toc,omitempty" xml:"toc,attr,omitempty"`
	// Name of this station
	Name string `json:"locname" xml:"locname,attr"`
	// True if this represents a Station, Bus stop or Ferry Terminal
	// i.e. Crs is present but does not start with X or Z
	Station bool `json:"station" xml:"station,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
}

func (a *Location) Equals(b *Location) bool {
	if b == nil {
		return false
	}
	return a.Tiploc == b.Tiploc &&
		a.Crs == b.Crs &&
		a.Toc == b.Toc &&
		a.Name == b.Name
}

// IsPublic returns true if this Location represents a public station.
// This is defined as having a Crs and one that does not start with X (Network Rail,
// some Bus stations and some Ferry terminals) and Z (usually London Underground).
//
// 2019 June 10 Enable Z for now as Farringdon is known as Farringdon Underground.
// This will expose the underground but better than leave a major station. Hopefully with Crossrail this will revert
// back to the single station.
func (t *Location) IsPublic() bool {
	return t.Crs != "" && t.Crs[0] != 'X' //&& t.Crs[0] != 'Z'
}
