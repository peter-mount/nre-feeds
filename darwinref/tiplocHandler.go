package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
)

func (dr *DarwinReference) TiplocHandler( r *rest.Rest ) error {
  return dr.View( func( tx *bolt.Tx ) error {
    tpl := r.Var( "id" )

    if location, exists := dr.GetTiploc( tx, tpl ); exists {
      location.SetSelf( r )
      r.Status( 200 ).Value( location )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}

func (dr *DarwinReference) TiplocsHandler( r *rest.Rest ) error {

  tiplocs := make( []string, 0 )

  if err := r.Body( &tiplocs ); err != nil {
    return err
  }

  return dr.View( func( tx *bolt.Tx ) error {
    var ary []*Location

    for _, tpl := range tiplocs {
      if location, exists := dr.GetTiploc( tx, tpl ); exists {
        location.SetSelf( r )
        ary = append( ary, location )
      }
    }

    r.Status( 200 ).Value( ary )

    return nil
  })
}
