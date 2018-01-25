package darwind3

import (
  "log"
)

// Processor interface
func (p *TS) Process( tx *Transaction ) error {
  log.Printf(
    "TS rid=%s uid=%s ssd=%s reverse=%s\n",
    p.RID,
    p.UID,
    p.SSD,
    p.ReverseFormation )

  /*
  if p.UR != nil {
    return p.UR.Process( d3, p )
  }
  */
  return nil
}
