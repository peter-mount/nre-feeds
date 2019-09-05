// Reference timetable
package darwintimetable

import (
	"encoding/xml"
	"time"
)

type Association struct {
	XMLName   xml.Name     `json:"-" xml:"Association"`
	Main      AssocService `json:"main" xml:"main"`
	Assoc     AssocService `json:"assoc" xml:"assoc"`
	Tiploc    string       `json:"tiploc" xml:"tiploc,attr"`
	Category  string       `json:"category" xml:"category,attr"`
	Cancelled bool         `json:"cancelled" xml:"isCancelled,attr"`
	Deleted   bool         `json:"deleted" xml:"isDeleted,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
}

func (a *Association) Equals(b *Association) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return a.Main.Equals(&b.Main) &&
		a.Assoc.Equals(&b.Assoc) &&
		a.Tiploc == b.Tiploc &&
		a.Category == b.Category
}

func (r *DarwinTimetable) addAssociation(a *Association) error {
	err := r.addJourneyAssociation(a, a.Main)
	if err != nil {
		return err
	}

	return r.addJourneyAssociation(a, a.Assoc)
}

type AssocService struct {
	RID string `json:"rid" xml:"rid,attr"`
	Wta string `json:"wta,omitempty" xml:"wta,attr"`
	Wtd string `json:"wtd,omitempty" xml:"wtd,attr"`
	Wtp string `json:"wtp,omitempty" xml:"wtp,attr"`
	Pta string `json:"pta,omitempty" xml:"pta,attr"`
	Ptd string `json:"ptd,omitempty" xml:"ptd,attr"`
}

func (a *AssocService) Equals(b *AssocService) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return a.RID == b.RID &&
		a.Wta == b.Wta &&
		a.Wtd == b.Wtd &&
		a.Wtp == b.Wtp &&
		a.Pta == b.Pta &&
		a.Ptd == b.Ptd
}

func (r *DarwinTimetable) addJourneyAssociation(a *Association, as AssocService) error {
	journey, exists := r.getJourney(as.RID)
	if !exists {
		return nil
	}

	for i, ja := range journey.Associations {
		if ja.Equals(a) {
			journey.Associations[i] = a
			return r.putJourney(journey)
		}
	}

	journey.Associations = append(journey.Associations, a)
	return r.putJourney(journey)
}
