package darwingraph

import (
	"encoding/json"
	"github.com/peter-mount/nre-feeds/darwinref"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Legolash = "Legolash2o"
)

func (d *DarwinGraph) importTiplocLocations() error {
	if *d.tiplocFileName == "" {
		return nil
	}

	log.Printf("Importing tiploc locations from %s", *d.tiplocFileName)

	f, err := os.Open(*d.tiplocFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	export := TiplocLocExport{}
	err = json.NewDecoder(f).Decode(&export)
	if err != nil {
		return err
	}

	for _, e := range export.Tiplocs {
		// Remove spaces, there's 1 tiploc in there with it
		tpl := strings.ReplaceAll(e.Tiploc, " ", "")
		n := d.graph.ComputeIfAbsent(tpl, func() RailNode {
			id, _ := strconv.ParseInt(tpl, IdBase, 64)
			return &TiplocNode{
				id:       id,
				Location: darwinref.Location{Tiploc: tpl, Name: e.Name, Crs: e.Details.CRS},
				LocSrc:   Legolash,
			}
		}).(*TiplocNode)

		if n.Crs == "" {
			n.Crs = e.Details.CRS
			n.LocSrc = Legolash
		}

		if n.Name == tpl || n.Name != e.Name {
			n.Name = e.Name
			n.LocSrc = Legolash
		}

		if !isNullIsland(e.Longitude) && !isNullIsland(e.Longitude) {
			n.Lon = e.Longitude
			n.Lat = e.Latitude
			n.LLSrc = Legolash
		}
	}

	log.Printf("Imported tiploc locations from %s", *d.tiplocFileName)

	return nil
}

// TiplocLocExport handles importing an extract of tiplocs to locations posted to
// the OpenRailData-Talk google group.
//
// You can download the json from https://groups.google.com/g/openraildata-talk/c/RI7rumYXM84/m/7ZvumfZyBAAJ
//
type TiplocLocExport struct {
	ExportDate string        `json:"ExportDate"` // Export date
	Tiplocs    []TiplocEntry `json:"Tiplocs"`    // Entries 1 per tiploc
}

// TiplocEntry represents individual entry in the json export.
type TiplocEntry struct {
	Tiploc    string  `json:"Tiploc"`    // Tiploc to map to
	Stanox    int     `json:"Stanox"`    // Stanox code
	Name      string  `json:"Name"`      // Name of location
	Latitude  float32 `json:"Latitude"`  // Latitude
	Longitude float32 `json:"Longitude"` // Longitude
	InTPS     bool    `json:"InTPS"`     //
	InBPlan   bool    `json:"InBPlan"`   //
	Details   struct {
		CRS                string `json:"CRS"`                 // CRS code
		OffNetwork         bool   `json:"OffNetwork"`          // Is tiploc off the network
		Zone               int    `json:"Zone"`                // Travel card zone
		Nalco              int    `json:"Nalco"`               //
		UIC                int    `json:"UIC"`                 //
		ForceLPB           string `json:"ForceLPB"`            //
		CompulsoryStop     bool   `json:"CompulsoryStop"`      //
		TpsStationType     string `json:"TPS_StationType"`     //
		TpsStationCategory string `json:"Tps_StationCategory"` //
		BPlanTimingPoint   string `json:"BPlan_TimingPoint"`   //
	} `json:"Details"`
}
