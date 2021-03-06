package darwinref

import (
	"encoding/json"
	"encoding/xml"
	bolt "github.com/etcd-io/bbolt"
	"time"
)

// A reason, shared by LateRunningReasons and CancellationReasons
type Reason struct {
	XMLName   xml.Name `json:"-" xml:"Reason"`
	Code      int      `json:"code" xml:"code,attr"`
	Text      string   `json:"reasontext" xml:"reasontext,attr"`
	Cancelled bool     `json:"canc" xml:"canc,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
}

func (a *Reason) Equals(b *Reason) bool {
	if b == nil {
		return false
	}
	return a.Code == b.Code &&
		a.Text == b.Text
}

// GetToc returns details of a TOC
func (r *DarwinReference) GetLateReason(tx *bolt.Tx, id int) (*Reason, bool) {
	loc, exists := r.GetReasonBucket(tx.Bucket([]byte("DarwinLateReason")), id)
	return loc, exists
}

// GetToc returns details of a TOC
func (r *DarwinReference) GetCancellationReason(tx *bolt.Tx, id int) (*Reason, bool) {
	loc, exists := r.GetReasonBucket(tx.Bucket([]byte("DarwinCancelReason")), id)
	return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getLateReason(id int) (*Reason, bool) {
	loc, exists := r.GetReasonBucket(r.lateRunningReasons, id)
	return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getCancellationReason(id int) (*Reason, bool) {
	loc, exists := r.GetReasonBucket(r.cancellationReasons, id)
	return loc, exists
}

func (t *Reason) FromBytes(b []byte) bool {
	if b != nil {
		err := json.Unmarshal(b, t)
		if err != nil {
			return false
		}
	}
	return t.Code != 0
}

func (r *DarwinReference) GetReasonBucket(bucket *bolt.Bucket, id int) (*Reason, bool) {
	var kb []byte = make([]byte, 3)
	kb[0] = byte(id)
	kb[1] = byte(id >> 8)
	kb[2] = byte(id >> 16)

	b := bucket.Get([]byte(kb))

	if b != nil {
		var toc *Reason = &Reason{}
		if toc.FromBytes(b) {
			return toc, true
		}
	}

	return nil, false
}

func addReason(bucket *bolt.Bucket, r *Reason) error {
	id := r.Code
	var kb []byte = make([]byte, 3)
	kb[0] = byte(id)
	kb[1] = byte(id >> 8)
	kb[2] = byte(id >> 16)

	// Get an existing entry & bail if it's the same
	b := bucket.Get([]byte(kb))
	if b != nil {
		var old *Reason = &Reason{}
		if old.FromBytes(b) && r.Equals(old) {
			return nil
		}
	}

	b, err := json.Marshal(r)

	if err != nil {
		return err
	}

	return bucket.Put(kb, b)
}
