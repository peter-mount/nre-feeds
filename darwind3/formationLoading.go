package darwind3

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"time"
)

// Loading data for an individual location in a schedule linked to a formation.
// Added in v16 2018-12-18 rttiPPTFormations_v1 Loading
type Loading struct {
	// The unique identifier of the formation data.
	// minLength 1, maxLength 20
	Fid string `json:"fid" xml:"fid"`
	// RTTI unique Train ID
	RID string `json:"rid" xml:"rid"`
	// TIPLOC where the loading data applies.
	Tiploc string `json:"tpl" xml:"tpl"`
	// Loading data for an individual coach in the formation.
	// If no loading data is provided for a coach in the formation then it
	// should be assumed to have been cleared.
	Loading []*CoachLoadingData `json:"loading" xml:"loading"`
	// attrbuteGroup CircularTimes
	Times util.CircularTimes `json:"time"`
	// This is the TS time from Darwin so we keep a copy of when this struct
	// was sent to us
	Date time.Time `json:"date,omitempty"`
}

func (s *Loading) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "fid":
			s.Fid = attr.Value

		case "rid":
			s.RID = attr.Value

		case "tpl":
			s.Tiploc = attr.Value
		}
	}

	// Parse CircularTimes attributes
	s.Times.UnmarshalXMLAttributes(start)

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elem interface{}
			switch tok.Name.Local {
			case "loading":
				cld := &CoachLoadingData{}
				s.Loading = append(s.Loading, cld)
				elem = cld

			default:
				if err := decoder.Skip(); err != nil {
					return err
				}
			}

			if elem != nil {
				if err := decoder.DecodeElement(elem, &tok); err != nil {
					return err
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}

// The CoachData & CoachLoadingData/LoadingValue complexTypes in rttiPPTFormations_v1
// We share the object to keep things simple.
type CoachLoadingData struct {
	// Data for an individual coach in a formation.
	// The number/identifier for this coach, e.g. "A" or "12".
	// minLength 1 maxLength 2
	CoachNumber string `json:"coachNumber" xml:"coachNumber,attr"`
	// The loading of this coach as a percentage 0...100
	Loading int `json:"loading" xml:",chardata"`
	// The class of the coach, e.g. "First" or "Standard".
	// CoachData only
	CoachClass string `json:"coachClass,omitempty" xml:"coachClass,attr,omitempty"`
	// The source of the loading data.
	// CoachLoadingData/LoadingValue only
	Src string `json:"src,omitempty" xml:"src,attr,omitempty"`
	// The RTTI instance ID of the src (if any).
	// CoachLoadingData/LoadingValue only
	// Length 2
	SrcInst string `json:"srcInst,omitempty" xml:"srcInst,attr,omitempty"`
	// The availability of a toilet in this coach.
	// E.g. "Unknown", "None" , "Standard" or "Accessible".
	// Note that other values may be supplied in the future without a schema change.
	// If no toilet availability is supplied then it should be assumed to be "Unknown".
	Toilet *Toilet `json:"toilet" xml:"toilet,omitempty"`
}
