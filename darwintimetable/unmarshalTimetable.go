// Reference timetable
package darwintimetable

import (
	"encoding/json"
	"encoding/xml"
	bolt "github.com/etcd-io/bbolt"
	"log"
	"time"
)

func (t *DarwinTimetable) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return t.internalUpdate(func(tx *bolt.Tx) error {
		return t.unmarshalXML(tx, decoder, start)
	})
}

func (t *DarwinTimetable) unmarshalXML(tx *bolt.Tx, decoder *xml.Decoder, start xml.StartElement) error {
	journeyCount := 0
	assocCount := 0

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "timetableID":
			t.timetableId = attr.Value
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:

			switch tok.Name.Local {
			case "Journey":
				var j *Journey = &Journey{}

				if err = decoder.DecodeElement(j, &tok); err != nil {
					return err
				}

				if err, updated := t.addJourney(j); err != nil {
					return err
				} else if updated {
					journeyCount++
				}

			case "Association":
				var a *Association = &Association{}

				err = decoder.DecodeElement(a, &tok)
				if err != nil {
					return err
				}

				err = t.addAssociation(a)
				if err != nil {
					return err
				}

				assocCount++

			default:
				log.Println("Unknown element", tok.Name.Local)
			}

		case xml.EndElement:

			log.Println("Journey's", journeyCount)
			log.Println("Association's", assocCount)

			// Finally update the meta data
			t.importDate = time.Now()

			b, err := json.Marshal(t)
			if err != nil {
				return err
			}

			return tx.Bucket([]byte("Meta")).Put([]byte("DarwinTimetable"), b)
		}
	}
}
