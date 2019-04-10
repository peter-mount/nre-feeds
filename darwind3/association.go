package darwind3

import (
	"encoding/json"
	"encoding/xml"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/util"
	"time"
)

// Type describing an association between schedules
type Association struct {
	// The through, previous working or link-to service
	Main AssocService `json:"main"`
	// The starting, terminating, subsequent working or link-from service
	Assoc AssocService `json:"assoc"`
	// The TIPLOC of the location where the association occurs.
	Tiploc string `json:"tiploc"`
	// Association Category Type: JJ=Join, VV=Split, LK=Linked, NP=Next-Working
	Category string `json:"category"`
	// True if this association is cancelled,
	// i.e. the association exists but will no longer happen.
	Cancelled bool `json:"cancelled,omitempty"`
	// True if this association is deleted,
	// i.e. the association no longer exists.
	Deleted bool `json:"deleted,omitempty"`
	// The schedule associated with this Association.
	// Note, d3 does nothing with this field. It's used in ldb to attach the appropriate schedule
	// in it's response
	Schedule *Schedule `json:"schedule,omitempty"`
	// This is the TS time from Darwin when this Association was updated
	Date time.Time `json:"date,omitempty"`
}

// xs:complexType name="AssocService"
type AssocService struct {
	// RTTI Train ID.
	// Note that since this is an RID, the service must already exist within Darwin.
	RID string `json:"rid"`
	// One or more scheduled times to identify the instance of the location
	// in the train schedule where the association occurs.
	Times util.CircularTimes `json:"timetable"`
	// Location for this entry
	Location *Location `json:"location,omitempty"`
	LocInd   int       `json:"locInd"`
	// The origin of this service
	Origin *Location `json:"origin,omitempty"`
	// The destination of this service
	Destination *Location `json:"destination,omitempty"`
}

func (a *Association) IsJoin() bool {
	return a != nil && a.Category == "JJ"
}

func (a *Association) IsSplit() bool {
	return a != nil && a.Category == "VV"
}

func (a *Association) IsNextTrain() bool {
	return a != nil && a.Category == "NP"
}

func (a *Association) Equals(b *Association) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return a.Tiploc == b.Tiploc &&
		a.Category == b.Category &&
		a.Main.Equals(&b.Main) &&
		a.Assoc.Equals(&b.Assoc)
}

func (a *AssocService) Equals(b *AssocService) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return a.RID == b.RID && a.Times.Equals(&b.Times)
}

func (a *Association) Clone() *Association {
	return &Association{
		Main:      a.Main,
		Assoc:     a.Assoc,
		Tiploc:    a.Tiploc,
		Category:  a.Category,
		Cancelled: a.Cancelled,
		Deleted:   a.Deleted,
		Date:      a.Date,
	}
}

func (d3 *DarwinD3) UpdateAssociations(sched *Schedule) {
	if d3.cache.tx != nil {
		d3.updateAssociations(d3.cache.tx, sched)
	} else {
		_ = d3.View(func(tx *bbolt.Tx) error {
			d3.updateAssociations(tx, sched)
			return nil
		})
	}
}

func (d3 *DarwinD3) updateAssociations(tx *bbolt.Tx, sched *Schedule) {
	assocs := getAssociations(tx, sched.RID)

	if assocs != nil {
		sched.Associations = assocs.Associations

		for _, a := range assocs.Associations {
			for idx, l := range sched.Locations {
				if a.Tiploc == l.Tiploc {
					np := a.Category == "NP"
					if a.Main.RID == sched.RID && (a.Category == "VV" || np) {
						a.Main.Location = l
						a.Main.LocInd = idx
						d3.updateAssociation(a, &a.Assoc, np)
					} else if a.Assoc.RID == sched.RID && (a.Category == "JJ" || np) {
						a.Assoc.Location = l
						a.Assoc.LocInd = idx
						d3.updateAssociation(a, &a.Main, np)
					}
				}
			}
		}
	}
}

func (d3 *DarwinD3) updateAssociation(a *Association, as *AssocService, np bool) {

	// np=true then do not resolve in the get else we could go into an infinite loop
	var s *Schedule
	if np {
		s = d3.GetScheduleNoResolve(as.RID)
	} else {
		s = d3.GetSchedule(as.RID)
	}

	if s != nil {
		//s.UpdateTime()
		as.Origin = s.Origin
		as.Destination = s.Destination
		for _, l := range s.Locations {
			if l.Times.Equals(&as.Times) {
				as.Location = l
				return
			}
		}
	}

	as.Location = nil
}

func (s *Association) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "tiploc":
			s.Tiploc = attr.Value

		case "category":
			s.Category = attr.Value

		case "isCancelled":
			s.Cancelled = attr.Value == "true"

		case "isDeleted":
			s.Deleted = attr.Value == "true"
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

			case "main":
				if err := decoder.DecodeElement(&s.Main, &tok); err != nil {
					return err
				}

			case "assoc":
				if err := decoder.DecodeElement(&s.Assoc, &tok); err != nil {
					return err
				}

			default:
				if err := decoder.Skip(); err != nil {
					return err
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}

func (s *AssocService) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "rid":
			s.RID = attr.Value
		}
	}
	s.Times.UnmarshalXMLAttributes(start)

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch token.(type) {
		case xml.StartElement:
			if err := decoder.Skip(); err != nil {
				return err
			}

		case xml.EndElement:
			return nil
		}
	}
}

// Entry in the AssociationBucket
type Associations struct {
	RID          string         `json:"rid"`
	Associations []*Association `json:"associations"`
}

func AssociationsFromBytes(b []byte) *Associations {
	if b == nil {
		return nil
	}

	associations := &Associations{}
	err := json.Unmarshal(b, associations)
	if err != nil || associations.RID == "" {
		return nil
	}
	return associations
}

func (associations *Associations) Bytes() ([]byte, error) {
	b, err := json.Marshal(associations)
	return b, err
}

func getAssociations(tx *bbolt.Tx, rid string) *Associations {
	b := tx.Bucket([]byte(AssociationBucket)).Get([]byte(rid))
	if b != nil {
		return AssociationsFromBytes(b)
	}
	return nil
}

func (associations *Associations) putAssociations(tx *bbolt.Tx) {
	b, _ := associations.Bytes()
	_ = tx.Bucket([]byte(AssociationBucket)).Put([]byte(associations.RID), b)
}
