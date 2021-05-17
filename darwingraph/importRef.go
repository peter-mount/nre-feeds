package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/darwinref"
	"io"
	"log"
	"os"
	"time"
)

func (d *DarwinGraph) importFile() error {
	log.Println("Importing", *d.importFileName)

	f, err := os.Open(*d.importFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	importXml := ImportXml{d: d.graph}
	err = xml.NewDecoder(f).Decode(&importXml)
	if err != nil {
		return err
	}

	log.Printf("Read %d locations", len(d.graph.ids))

	return nil
}

type ImportXml struct {
	d     *TiplocGraph // direct will allow import faster as not locked
	count int
}

func (r *ImportXml) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	date := time.Now()
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
			if tok.Name.Local == "LocationRef" {
				loc := darwinref.Location{}
				if err = decoder.DecodeElement(&loc, &tok); err != nil {
					return err
				}
				loc.Date = date
				loc.Station = loc.IsPublic()

				// Ensure we have an entry
				node := r.d.ComputeIfAbsent(loc.Tiploc, func() *TiplocNode {
					return &TiplocNode{LocSrc: "NreRef"}
				})

				// Update the location to the parsed one
				node.Location = loc
			}
		}
	}

}
