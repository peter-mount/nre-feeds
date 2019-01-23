package darwinref

import (
  "encoding/json"
)

type crsimport struct {
   crs  map[string]map[string]string
   ref *DarwinReference
}

func (r *DarwinReference) newCrsImport() *crsimport {
  return &crsimport{
    crs: make( map[string]map[string]string ),
    ref: r,
  }
}

func (c *crsimport) append( loc *Location ) {
  if loc.Crs != "" {
    if e, ok := c.crs[ loc.Crs ]; ok {
      e[ loc.Tiploc ] = loc.Tiploc
    } else {
      c.crs[ loc.Crs ] = make( map[string]string )
      c.crs[ loc.Crs ][ loc.Tiploc ] = loc.Tiploc
    }
  }
}

func (c *crsimport) write() (error, int) {
  count := 0
  for k,v := range c.crs {
    if err, updated := c.writeCrs( k, v ); err != nil {
      return err, count
    } else if updated {
      count ++
    }
  }

  return nil, count
}

func (c *crsimport) writeCrs( k string, v map[string]string ) (error,bool) {
  // Now look at existing entry & skip if we match
  b := c.ref.crs.Get( []byte( k ) )
  if b != nil {
    var ar []string
    err := json.Unmarshal( b, &ar )
    if err != nil {
      return err, false
    }

    if crsCompare( v, ar ) {
      return nil, false
    }
  }

  // get loaded tiplocs
  var tiplocs []string
  for t,_ := range v {
    tiplocs = append( tiplocs, t )
  }

  b, err := json.Marshal( tiplocs )
  if err != nil {
    return err, false
  }

  if err := c.ref.crs.Put( []byte( k ), b ); err != nil {
    return err, false
  }

  return nil, true
}

func crsCompare( m map[string]string, ar []string ) bool {
  if len( ar ) == 0 {
    return false
  }

  for _, t := range ar {
    if _, exists := m[t]; !exists {
      return false
    }
  }
  return true
}
