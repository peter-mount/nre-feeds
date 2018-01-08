// Decode's a PportTimetableRef into a DarwinReference which is more suitable
// for use as it map's the data accordingly

package darwinref

import (
  "log"
)

func (s *PportTimetableRef) Decode() *DarwinReference {
  var d *DarwinReference = new( DarwinReference )

  d.Tiploc = make( map[string]*Location )
  d.Crs = make( map[string][]*Location )
  for _, loc := range s.Locations {

    if _, exists := d.Tiploc[ loc.Tiploc ]; exists {
      log.Println( "Tiploc", loc.Tiploc, "duplicated" )
    } else {
      d.Tiploc[ loc.Tiploc ] = loc
    }

    if loc.Crs != "" {
      d.Crs[ loc.Crs ] = append( d.Crs[ loc.Crs ], loc )
    }

  }

  d.Toc = make( map[string]*Toc )
  for _, toc := range s.Toc {
    if _, exists := d.Toc[ toc.Toc ]; exists {
      log.Println( "Toc", s.Toc, "duplicated" )
    } else {
      d.Toc[ toc.Toc ] = toc
    }
  }

  d.LateRunningReasons = make( map[int]string )
  for _, reason := range s.LateRunningReasons {
    d.LateRunningReasons[ reason.Code ] = reason.Text
  }

  d.CancellationReasons = make( map[int]string )
  for _, reason := range s.CancellationReasons {
    d.CancellationReasons[ reason.Code ] = reason.Text
  }

  d.CISSource = make( map[string]string )
  for _, cis := range s.CISSource {
    d.CISSource[ cis.Code ] = cis.Name
  }

  return d
}
