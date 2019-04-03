package darwinref

import (
	"bytes"
	"encoding/json"
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/golib/rest"
	"sort"
	"strconv"
)

// ReasonMap allows a set of Reasons (either Late or Cancelled) to be built up
// usually from a set of schedules - e.g. ldb.Service
type ReasonMap struct {
	Late      map[int]*Reason
	Cancelled map[int]*Reason
}

func NewReasonMap() *ReasonMap {
	return &ReasonMap{
		Late:      make(map[int]*Reason),
		Cancelled: make(map[int]*Reason),
	}
}

// Add a Reason to the map
// t *Reason
// tx Bolt transaction
// f function to retrieve, usually DarwinReference.GetLateReason
// or DarwinReference.GetCancellationReason
func (r *ReasonMap) Add(id int, canc bool, tx *bolt.Tx, dr *DarwinReference) {
	var m map[int]*Reason
	var f func(*bolt.Tx, int) (*Reason, bool)
	if canc {
		m = r.Cancelled
		f = dr.GetCancellationReason
	} else {
		m = r.Late
		f = dr.GetLateReason
	}

	if _, ok := m[id]; !ok {
		if rr, ok := f(tx, id); ok {
			m[rr.Code] = rr
		}
	}
}

func (r *ReasonMap) AddReason(reason *Reason) {
	if reason != nil && reason.Code > 0 {
		if reason.Cancelled {
			r.Cancelled[reason.Code] = reason
		} else {
			r.Late[reason.Code] = reason
		}
	}
}

// Self sets the Self field to match this request
func (r *ReasonMap) Self(rs *rest.Rest) {
	for _, v := range r.Cancelled {
		v.SetSelf(rs)
	}
	for _, v := range r.Late {
		v.SetSelf(rs)
	}
}

func (r *ReasonMap) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	if len(r.Late) == 0 && len(r.Cancelled) == 0 {
		b.WriteString("null")
	} else {

		b.WriteByte('{')

		sep, err := r.marshalJSON(b, "late", r.Late, false)
		if err != nil {
			return nil, err
		}

		sep, err = r.marshalJSON(b, "cancelled", r.Cancelled, sep)
		if err != nil {
			return nil, err
		}

		b.WriteByte('}')
	}

	return b.Bytes(), nil
}

func (r *ReasonMap) marshalJSON(b *bytes.Buffer, l string, m map[int]*Reason, notFirst bool) (bool, error) {
	if len(m) == 0 {
		return notFirst, nil
	}

	if notFirst {
		b.WriteByte(',')
	}

	b.WriteByte('"')
	b.WriteString(l)
	b.WriteString("\":{")

	var vals []*Reason
	for _, v := range m {
		vals = append(vals, v)
	}

	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].Code < vals[j].Code
	})

	for i, v := range vals {
		if i > 0 {
			b.WriteByte(',')
		}
		// This is a map so the keys MUST be strings
		b.WriteByte('"')
		b.WriteString(strconv.FormatInt(int64(v.Code), 10))
		b.WriteByte('"')
		b.WriteByte(':')

		if eb, err := json.Marshal(v); err != nil {
			return notFirst, err
		} else {
			b.Write(eb)
		}
	}

	b.WriteByte('}')

	return true, nil
}
