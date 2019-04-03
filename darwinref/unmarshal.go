package darwinref

import (
	"encoding/json"
	"encoding/xml"
	bolt "github.com/etcd-io/bbolt"
	"log"
	"time"
)

func (r *DarwinReference) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return r.internalUpdate(func(tx *bolt.Tx) error {
		return r.unmarshalXML(tx, decoder, start)
	})
}

func (r *DarwinReference) unmarshalXML(tx *bolt.Tx, decoder *xml.Decoder, start xml.StartElement) error {
	date := time.Now()
	crs := r.newCrsImport()
	tplCount := 0
	tocCount := 0
	viaCount := 0

	r.cisSource = make(map[string]string)

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "timetableId":
			r.timetableId = attr.Value
		}
	}

	// Reason map to write to
	var late bool
	var inReason bool

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "LocationRef":
				var loc *Location = &Location{}
				if err = decoder.DecodeElement(loc, &tok); err != nil {
					return err
				}
				loc.Date = date

				if err, updated := r.addTiploc(loc); err != nil {
					return err
				} else if updated {
					tplCount++
				}

				// Append to CRS map
				crs.append(loc)

			case "TocRef":
				var toc *Toc = &Toc{}
				if err = decoder.DecodeElement(toc, &tok); err != nil {
					return err
				}

				if err, updated := r.addToc(toc); err != nil {
					return err
				} else if updated {
					tocCount++
				}

			case "LateRunningReasons":
				inReason = true
				late = true

			case "CancellationReasons":
				inReason = true
				late = false

			case "Reason":
				if inReason {
					var reason *Reason = &Reason{}
					if err = decoder.DecodeElement(reason, &tok); err != nil {
						return err
					}

					reason.Cancelled = !late
					reason.Date = date
					if late {
						if err = addReason(r.lateRunningReasons, reason); err != nil {
							return err
						}
					} else {
						if err = addReason(r.cancellationReasons, reason); err != nil {
							return err
						}
					}
				}

			case "CISSource":
				var cis *CISSource = &CISSource{}
				if err = decoder.DecodeElement(cis, &tok); err != nil {
					return err
				}
				r.cisSource[cis.Code] = cis.Name

			case "Via":
				var via *Via = &Via{}
				if err = decoder.DecodeElement(via, &tok); err != nil {
					return err
				}

				if err, updated := r.addVia(via); err != nil {
					return err
				} else if updated {
					viaCount++
				}

			default:
				log.Println("Unknown element", tok.Name.Local)
			}

		case xml.EndElement:
			if !inReason {
				log.Printf("Imported %d Tiplocs", tplCount)

				if err, count := crs.write(); err != nil {
					return err
				} else {
					log.Printf("Imported %d CRS", count)
				}

				log.Printf("Imported %d TOC's", tocCount)

				log.Printf("Imported %d Via's", viaCount)

				// Finally update the meta data
				r.importDate = time.Now()

				b, err := json.Marshal(r)
				if err != nil {
					return err
				}
				return tx.Bucket([]byte("Meta")).Put([]byte("DarwinReference"), b)
			}
			inReason = false
		}
	}

}
