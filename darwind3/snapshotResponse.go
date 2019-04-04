package darwind3

import (
	"encoding/xml"
	"log"
)

// Snapshot Response
type SR struct {
	XMLName xml.Name `json:"-" xml:"sR"`
	Actions []Processor
}

func (s *SR) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elem Processor
			switch tok.Name.Local {
			case "schedule":
				elem = &Schedule{}

			case "deactivated":
				elem = &DeactivatedSchedule{}

			case "TS":
				elem = &TS{}

			case "trainOrder":
				elem = &trainOrderWrapper{}

			case "OW":
				elem = &StationMessage{}

			case "formationLoading":
				elem = &Loading{}

			case "association":
				elem = &Association{}

			case "trackingID":
				elem = &TrackingID{}

			case "alarm":
				elem = &RttiAlarm{}

			// Unsupported (so far) elements:
			// scheduleFormations
			// trainAlert
			default:
				log.Println("Skipping", tok.Name.Local, tok.Name.Space)
				if err := decoder.Skip(); err != nil {
					return err
				}
			}

			if elem != nil {
				if err := decoder.DecodeElement(elem, &tok); err != nil {
					return err
				}
				s.Actions = append(s.Actions, elem)
			}

		case xml.EndElement:
			return nil
		}
	}
}

// Process this message
func (p *SR) Process(tx *Transaction) error {

	if len(p.Actions) > 0 {
		for _, s := range p.Actions {
			if err := s.Process(tx); err != nil {
				return err
			}
		}
	}

	return nil
}
