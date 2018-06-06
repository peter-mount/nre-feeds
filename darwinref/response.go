package darwinref

import (
  "encoding/xml"
)

type CrsResponse struct {
  XMLName     xml.Name            `json:"-" xml:"crs"`
  Crs         string              `json:"crs" xml:"crs,attr"`
  Tiploc   []*Location            `json:"locations,omitempty" xml:"LocationRef"`
  Self        string              `json:"self,omitempty" xml:"self,attr,omitempty"`
}

type ReasonsResponse struct {
  XMLName     xml.Name  `json:"-" xml:"reasons"`
  Reasons  []*Reason    `json:"reasons,omitempty" xml:"Reason"`
  Self        string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

type TocsResponse struct {
  XMLName     xml.Name  `json:"-" xml:"tocs"`
  Toc      []*Toc       `json:"tocs,omitempty" xml:"TocRef"`
  Self        string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

type SearchResult struct {
  Crs       string      `json:"code"`
  Name      string      `json:"name"`
  Label     string      `json:"label"`
  Score     float64     `json:"score"`
  Distance  float64     `json:"distance,omitempty"`
}

// An entry in the request object either received by the rest endpoint but also
// used by DarwinRefClient.GetVias().
type ViaResolveRequest struct {
  // CRS of the location we want to show a via
  Crs           string    `json:"crs"`
  // Destination tiploc
  Destination   string    `json:"destination"`
  // Tiplocs of journey after this location to search
  Tiplocs     []string    `json:"tpls"`
}
