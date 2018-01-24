package darwinrest

import (
  bolt "github.com/coreos/bbolt"
  "darwinref"
  "darwintimetable"
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
)

type result struct {
  XMLName     xml.Name                  `json:"-" xml:"result"`
  RID         string                    `json:"rid" xml:"rid,attr"`
  Journey    *darwintimetable.Journey   `json:"journey" xml:"journey"`
  Locations  *darwinref.LocationMap     `json:"locations" xml:"locations>LocationRef"`
  Self        string                    `json:"self" xml:"self,attr"`
}

// JourneyHandler returns a Journey from the timetable and any reference data
func (rs *DarwinRest) JourneyHandler( r *rest.Rest ) error {
  res := &result{ RID: r.Var( "rid" ) }

  if err := rs.TT.View( func( tx *bolt.Tx ) error {
    if journey, exists := rs.TT.GetJourney( tx, res.RID ); exists {
      journey.SetSelf( r )
      res.Journey = journey
    }
    return nil
  }); err != nil {
    return err
  }

  if res.Journey == nil {
    r.Status( 404 )
    return nil
  }

  res.Locations = darwinref.NewLocationMap()
  if err := rs.Ref.View( func( tx *bolt.Tx ) error {

    var tpls []string
    for _, l := range res.Journey.Schedule {
      tpls = append( tpls, l.Tiploc )
    }

    res.Locations.AddTiplocs( rs.Ref, tx, tpls )

    return nil
  }); err != nil {
    return err
  }

  res.Locations.Self( r )
  res.Self = r.Self( r.Context() + "/journey/" + res.RID )
  r.Status( 200 ).Value( res )
  return nil
}
