package darwinref

import (
	"encoding/xml"
)

type CrsResponse struct {
	XMLName xml.Name    `json:"-" xml:"crs"`
	Crs     string      `json:"crs" xml:"crs,attr"`
	Tiploc  []*Location `json:"locations,omitempty" xml:"LocationRef"`
}

type ReasonsResponse struct {
	XMLName xml.Name  `json:"-" xml:"reasons"`
	Reasons []*Reason `json:"reasons,omitempty" xml:"Reason"`
}

type TocsResponse struct {
	XMLName xml.Name `json:"-" xml:"tocs"`
	Toc     []*Toc   `json:"tocs,omitempty" xml:"TocRef"`
}

type SearchResult struct {
	Crs      string  `json:"code"`
	Name     string  `json:"name"`
	Label    string  `json:"label"`
	Score    float64 `json:"score"`
	Distance float64 `json:"distance,omitempty"`
}

// An entry in the request object either received by the rest endpoint but also
// used by DarwinRefClient.GetVias().
type ViaResolveRequest struct {
	// CRS of the location we want to show a via
	Crs string `json:"crs"`
	// Destination tiploc
	Destination string `json:"destination"`
	// Tiplocs of journey after this location to search
	Tiplocs []string `json:"tpls"`
}

func (vr *ViaResolveRequest) AppendTiploc(tiploc string) {
	if tiploc != "" {
		vr.Tiplocs = append(vr.Tiplocs, tiploc)
	}
}
