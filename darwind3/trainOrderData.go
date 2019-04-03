package darwind3

import (
	"encoding/xml"
)

// Defines the sequence of trains making up the train order
type trainOrderData struct {
	// The first train in the train order.
	First *trainOrderItem `xml:"first"`
	// The second train in the train order.
	Second *trainOrderItem `xml:"second"`
	// The third train in the train order.
	Third *trainOrderItem `xml:"third"`
}

func (s *trainOrderData) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var ptr **trainOrderItem
			switch tok.Name.Local {
			case "first":
				ptr = &s.First

			case "second":
				ptr = &s.Second

			case "third":
				ptr = &s.Third

			default:
				if err := decoder.Skip(); err != nil {
					return err
				}
			}

			if ptr != nil {
				*ptr = &trainOrderItem{}
				if err := decoder.DecodeElement(*ptr, &tok); err != nil {
					return err
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}
