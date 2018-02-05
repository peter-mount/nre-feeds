package darwind3

type Transaction struct {
  // The root
  pport       *Pport
  //
  d3          *DarwinD3
}

// Processor interface used by some types used when processing a message and
// updating our internal state
type Processor interface {
  Process( *Transaction ) error
}

type KBProcessor interface {
  Process() error
}

func (d *DarwinD3) ProcessUpdate( p *Pport, f func( *Transaction ) error ) error {
  t := &Transaction{
    pport: p,
    d3: d,
  }
  return f( t )
}
