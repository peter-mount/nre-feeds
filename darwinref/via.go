package darwinref

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	bolt "github.com/etcd-io/bbolt"
	"time"
)

// Via text
type Via struct {
	XMLName xml.Name `json:"-" xml:"Via"`
	At      string   `json:"at" xml:"at,attr"`
	Dest    string   `json:"dest" xml:"dest,attr"`
	Loc1    string   `json:"loc1" xml:"loc1,attr"`
	Loc2    string   `json:"loc2,omitempty" xml:"loc2,attr,omitempty"`
	Text    string   `json:"text" xml:"viatext,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
}

// Are two Via's equal
func (v *Via) Equals(o *Via) bool {
	if o == nil {
		return false
	}
	return v.At == o.At && v.Dest == o.Dest && v.Loc1 == o.Loc1 && v.Loc2 == o.Loc2
}

// Key the unique key for this entry
func (v *Via) key() string {
	return fmt.Sprintf("%s %s %s %s", v.At, v.Dest, v.Loc1, v.Loc2)
}

func (v *Via) String() string {
	return "Via[At=" + v.At + ", Dest=" + v.Dest + ", Loc1=" + v.Loc1 + ", Loc2=" + v.Loc2 + ", Text=" + v.Text + "]"
}

// GetToc returns details of a TOC
func (r *DarwinReference) GetVia(tx *bolt.Tx, at string, dest string, loc1 string, loc2 string) (*Via, bool) {
	loc, exists := r.GetViaBucket(tx.Bucket([]byte("DarwinVia")), at, dest, loc1, loc2)
	return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getVia(at string, dest string, loc1 string, loc2 string) (*Via, bool) {
	loc, exists := r.GetViaBucket(r.via, at, dest, loc1, loc2)
	return loc, exists
}

func (r *DarwinReference) GetViaBucket(bucket *bolt.Bucket, at string, dest string, loc1 string, loc2 string) (*Via, bool) {
	key := fmt.Sprintf("%s %s %s %s", at, dest, loc1, loc2)
	b := bucket.Get([]byte(key))

	if b != nil {
		var via *Via = &Via{}
		err := json.Unmarshal(b, via)
		if err == nil {
			return via, true
		}
	}

	return nil, false
}

func (r *DarwinReference) addVia(via *Via) (error, bool) {
	// SnapshotUpdate only if it does not exist or is different
	if old, exists := r.getVia(via.At, via.Dest, via.Loc1, via.Loc2); !exists || !via.Equals(old) {
		via.Date = time.Now()

		b, err := json.Marshal(via)
		if err != nil {
			return err, false
		}

		err = r.via.Put([]byte(via.key()), b)
		if err != nil {
			return err, false
		}

		return nil, true
	}

	return nil, false
}
