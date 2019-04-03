package darwinref

import (
	"bytes"
	"encoding/json"
	bolt "github.com/etcd-io/bbolt"
	//  "encoding/xml"
	"github.com/peter-mount/golib/rest"
	"sort"
	"strings"
)

type LocationMap struct {
	m map[string]*Location
}

func NewLocationMap() *LocationMap {
	return &LocationMap{m: make(map[string]*Location)}
}

// AddTiploc adds a Tiploc to the response
func (r *LocationMap) Add(t *Location) {
	if _, ok := r.m[t.Tiploc]; !ok {
		r.m[t.Tiploc] = t
	}
}

// AddTiplocs adds an array of Tiploc's to the response
func (r *LocationMap) AddAll(t []*Location) {
	for _, e := range t {
		r.Add(e)
	}
}

func (r *LocationMap) AddTiploc(dr *DarwinReference, tx *bolt.Tx, t string) {
	if _, ok := r.m[t]; !ok {
		if loc, exists := dr.GetTiploc(tx, t); exists {
			r.m[t] = loc
		}
	}
}

func (r *LocationMap) AddTiplocs(dr *DarwinReference, tx *bolt.Tx, ts []string) {
	bucket := tx.Bucket([]byte("DarwinTiploc"))
	for _, t := range ts {
		if _, ok := r.m[t]; !ok {
			if loc, exists := dr.GetTiplocBucket(bucket, t); exists {
				r.m[t] = loc
			}
		}
	}
}

func (r *LocationMap) Get(n string) (*Location, bool) {
	t, e := r.m[n]
	return t, e
}

// Self sets the Self field to match this request
func (r *LocationMap) Self(rs *rest.Rest) {
	for _, v := range r.m {
		v.Self = rs.Self("/ref/tiploc/" + v.Tiploc)
	}
}

func (r *LocationMap) ForEach(f func(*Location)) {
	for _, v := range r.m {
		f(v)
	}
}

func (t *LocationMap) MarshalJSON() ([]byte, error) {
	// Tiploc sorted by NLC
	var vals []*Location
	for _, v := range t.m {
		vals = append(vals, v)
	}

	sort.SliceStable(vals, func(i, j int) bool {
		return strings.Compare(vals[i].Tiploc, vals[j].Tiploc) < 0
	})

	b := &bytes.Buffer{}
	b.WriteByte('{')

	for i, v := range vals {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		b.WriteString(v.Tiploc)
		b.WriteByte('"')
		b.WriteByte(':')

		if eb, err := json.Marshal(v); err != nil {
			return nil, err
		} else {
			b.Write(eb)
		}
	}

	b.WriteByte('}')
	return b.Bytes(), nil
}
