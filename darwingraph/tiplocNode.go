package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/darwinref"
	"github.com/peter-mount/nre-feeds/util"
	"strconv"
)

// TiplocNode is a node on our generated map
type TiplocNode struct {
	id                 int64   // ID of this node
	darwinref.Location         // Location entry
	LocSrc             string  // Source of Location
	Lat                float32 // Latitude of point
	Lon                float32 // Longitude of point
	LLSrc              string  // Source of Lat/Lon
}

func (n TiplocNode) ID() int64 {
	return n.id
}

func (n TiplocNode) String() string {
	return n.Location.Tiploc
}

// isNullIsland returns true if the coordinate is within 1 second of 0.
func isNullIsland(v float32) bool {
	return v >= -0.0002777778 && v <= 0.0002777778
}

func (n *TiplocNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "id"}, strconv.FormatInt(n.id, 36)).
		AddAttribute(xml.Name{Local: "tpl"}, n.Tiploc).
		AddAttributeIfSet(xml.Name{Local: "crs"}, n.Crs).
		AddAttributeIfSet(xml.Name{Local: "name"}, n.Name).
		AddAttributeIfSet(xml.Name{Local: "toc"}, n.Toc).
		AddBoolAttributeIfSet(xml.Name{Local: "station"}, n.Station).
		AddAttribute(xml.Name{Local: "locSrc"}, n.LocSrc).
		// Coordinates only if they are not at null point island
		If(!isNullIsland(n.Lon) && !isNullIsland(n.Lat)).
		AddFloat32Attribute(xml.Name{Local: "lat"}, n.Lat).
		AddFloat32Attribute(xml.Name{Local: "lon"}, n.Lon).
		AddAttributeIfSet(xml.Name{Local: "llSrc"}, n.LLSrc).
		EndIf().
		Build()
}

func (n *TiplocNode) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		var err error
		var f float64

		switch attr.Name.Local {
		case "id":
			n.id, err = strconv.ParseInt(attr.Value, 36, 64)
		case "tpl":
			n.Tiploc = attr.Value
		case "crs":
			n.Crs = attr.Value
		case "name":
			n.Name = attr.Value
		case "toc":
			n.Toc = attr.Value
		case "station":
			if attr.Value == "true" {
				n.Station = true
			} else {
				n.Station = false
			}
		case "lat":
			f, err = strconv.ParseFloat(attr.Value, 32)
			n.Lat = float32(f)
		case "lon":
			f, err = strconv.ParseFloat(attr.Value, 32)
			n.Lon = float32(f)
		case "llSrc":
			n.LLSrc = attr.Value
		case "locSrc":
			n.LocSrc = attr.Value
		}

		if err != nil {
			return err
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch token.(type) {
		case xml.EndElement:
			return nil
		}
	}
}
