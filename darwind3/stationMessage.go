package darwind3

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"time"
)

type StationMessage struct {
	ID uint64 `json:"id" xml:"id,attr"`
	// The message
	Message string `json:"message" xml:"message"`
	// CRS codes for the stations this message applies
	Station []string `json:"station" xml:"stations>station"`
	// The category of message
	Category string `json:"category" xml:"category,attr"`
	// The severity of the message
	Severity int `json:"severity" xml:"severity,attr"`
	// Whether the train running information is suppressed to the public
	Suppress bool `json:"suppress,omitempty" xml:"suppress,attr,omitempty"`
	// Usually this is the date we insert into the db but here we use the TS time
	// as returned from darwin
	Date time.Time `json:"date,omitempty" xml:"date,attr,omitempty"`
	// URL to this entity
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// Bytes returns the message as an encoded byte slice
func (s *StationMessage) Bytes() ([]byte, error) {
	b, err := json.Marshal(s)
	return b, err
}

// ScheduleFromBytes returns a schedule based on a slice or nil if none
func StationMessageFromBytes(b []byte) *StationMessage {
	if b == nil {
		return nil
	}

	message := &StationMessage{}
	err := json.Unmarshal(b, message)
	if err != nil {
		return nil
	}
	return message
}

func (s *StationMessage) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			if i, err := strconv.ParseInt(attr.Value, 10, 64); err != nil {
				return err
			} else {
				s.ID = uint64(i)
			}

		case "cat":
			s.Category = attr.Value

		case "sev":
			if i, err := strconv.Atoi(attr.Value); err != nil {
				return err
			} else {
				s.Severity = i
			}

		case "suppress":
			s.Suppress = attr.Value == "true"
		}
	}

	inMsg := false
	var nest int
	var msg []byte

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "Station":
				for _, attr := range tok.Attr {
					if attr.Name.Local == "crs" {
						s.Station = append(s.Station, attr.Value)
					}
				}

			case "Msg":
				if !inMsg {
					inMsg = true
					nest = 1
				}

			default:
				if inMsg {
					nest++

					msg = append(msg, '<')
					msg = append(msg, []byte(tok.Name.Local)...)
					for _, attr := range tok.Attr {
						msg = append(msg, ' ')
						msg = append(msg, []byte(attr.Name.Local)...)
						msg = append(msg, '=', '"')
						msg = append(msg, []byte(attr.Value)...)
						msg = append(msg, '"')
					}
					msg = append(msg, '>')
				} else {
					if err := decoder.Skip(); err != nil {
						return err
					}
				}
			}

		case xml.CharData:
			if inMsg {
				msg = append(msg, tok...)
			}

		case xml.EndElement:
			switch tok.Name.Local {
			case "Station":

			case "Msg":
				inMsg = false
				s.Message = string(msg[:])

			case "OW":
				return nil

			default:
				if nest > 0 {
					msg = append(msg, '<', '/')
					msg = append(msg, []byte(tok.Name.Local)...)
					msg = append(msg, '>')
					nest--
				}
			}
		}
	}
}
