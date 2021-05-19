package darwingraph

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strconv"
)

// importKBStations imports the station data from the NRE KB feed.
// We use the retrieved xml to make things easier than hitting our rest service.
//
// The main portion of this is to set/update the coordinates of the stations
func (d *DarwinGraph) importKBStations() error {
	log.Println("Importing", *d.stationsFileName)

	f, err := os.Open(*d.stationsFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	importXml := ImportStations{d: d.graph}
	err = xml.NewDecoder(f).Decode(&importXml)
	if err != nil {
		return err
	}

	log.Println("Imported", *d.stationsFileName)
	return nil
}

type ImportStations struct {
	d         *RailGraph
	inStation bool
	crs       string
	lat       string
	lon       string
}

func getInnerContent(decoder *xml.Decoder, tok *xml.StartElement) (string, error) {
	var v string
	err := decoder.DecodeElement(&v, tok)
	if err != nil {
		return "", err
	}
	return v, nil
}

func (r *ImportStations) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			if io.EOF == err {
				return nil
			}
			return err
		}
		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "Station":
				r.crs = ""
				r.lon = ""
				r.lat = ""
			case "CrsCode":
				r.crs, err = getInnerContent(decoder, &tok)
				if err != nil {
					return err
				}
			case "Longitude":
				r.lon, err = getInnerContent(decoder, &tok)
				if err != nil {
					return err
				}
			case "Latitude":
				r.lat, err = getInnerContent(decoder, &tok)
				if err != nil {
					return err
				}
			}
		case xml.EndElement:
			switch tok.Name.Local {
			case "StationList":
				return nil
			case "Station":
				if r.lat != "" && r.lon != "" {
					lat, err := strconv.ParseFloat(r.lat, 32)
					if err != nil {
						return err
					}
					lon, err := strconv.ParseFloat(r.lon, 32)
					if err != nil {
						return err
					}
					crs := r.d.GetCrs(r.crs)
					if crs != nil {
						for _, n := range crs.tiploc {
							// Only update if not Null Island & either no source or a previous NreKB entry
							// This prevents invalid points and we don't overwrite custom entries
							if !isNullIsland(float32(lat)) && !isNullIsland(float32(lon)) &&
								n.LLSrc == "" || n.LLSrc == "NreKB" {
								n.Lat = float32(lat)
								n.Lon = float32(lon)
								n.LLSrc = "NreKB"
							}
						}
					}
				}
			}
		}
	}
}
