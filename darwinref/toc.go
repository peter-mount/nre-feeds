package darwinref

import (
	"encoding/json"
	"encoding/xml"
	bolt "github.com/etcd-io/bbolt"
	"time"
)

// A rail operator
type Toc struct {
	XMLName xml.Name `json:"-" xml:"TocRef"`
	Toc     string   `json:"toc" xml:"toc,attr"`
	Name    string   `json:"tocname" xml:"tocname,attr"`
	Url     string   `json:"url" xml:"url,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
}

func (a *Toc) Equals(b *Toc) bool {
	if b == nil {
		return false
	}
	return a.Toc == b.Toc &&
		a.Name == b.Name &&
		a.Url == b.Url
}

// GetToc returns details of a TOC
func (r *DarwinReference) GetToc(tx *bolt.Tx, toc string) (*Toc, bool) {
	loc, exists := r.GetTocBucket(tx.Bucket([]byte("DarwinToc")), toc)
	return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getToc(tpl string) (*Toc, bool) {
	loc, exists := r.GetTocBucket(r.toc, tpl)
	return loc, exists
}

func (t *Toc) FromBytes(b []byte) bool {
	if b != nil {
		err := json.Unmarshal(b, t)
		if err != nil {
			return false
		}
	}
	return t.Toc != ""
}

func (r *DarwinReference) GetTocBucket(bucket *bolt.Bucket, tpl string) (*Toc, bool) {
	b := bucket.Get([]byte(tpl))

	if b != nil {
		var toc *Toc = &Toc{}
		if toc.FromBytes(b) {
			return toc, true
		}
	}

	return nil, false
}

func (r *DarwinReference) addToc(toc *Toc) (error, bool) {
	// SnapshotUpdate only if it does not exist or is different
	if old, exists := r.getToc(toc.Toc); !exists || !toc.Equals(old) {
		toc.Date = time.Now()

		b, err := json.Marshal(toc)
		if err != nil {
			return err, false
		}

		if err := r.toc.Put([]byte(toc.Toc), b); err != nil {
			return err, false
		}

		return nil, true
	}

	return nil, false
}
