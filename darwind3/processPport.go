package darwind3

import "github.com/etcd-io/bbolt"

type Transaction struct {
	// The root
	pport *Pport
	//
	d3 *DarwinD3
}

// Processor interface used by some types used when processing a message and
// updating our internal state
type Processor interface {
	Process(*Transaction) error
}

type KBProcessor interface {
	Process() error
}

func (d *DarwinD3) ProcessUpdate(p *Pport, f func(*Transaction) error) error {
	t := &Transaction{
		pport: p,
		d3:    d,
	}
	defer t.close()
	return f(t)
}

func (t *Transaction) close() {
	if !t.pport.TS.IsZero() {
		_ = t.d3.UpdateBulkAware(func(tx *bbolt.Tx) error {
			return PutMeta(tx, "ts", t.pport.TS)
		})
	}
}
