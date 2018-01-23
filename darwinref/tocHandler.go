package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
)

func (dr *DarwinReference) TocHandler( r *rest.Rest ) error {
  return dr.View( func( tx *bolt.Tx ) error {
    id := r.Var( "id" )

    if toc, exists := dr.GetToc( tx, id ); exists {
      toc.SetSelf( r )
      r.Status( 200 ).Value( toc )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}

type TocsResponse struct {
  XMLName     xml.Name  `json:"-" xml:"tocs"`
  Toc      []*Toc       `json:"tocs,omitempty" xml:"TocRef"`
  Self        string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func (dr *DarwinReference) AllTocsHandler( r *rest.Rest ) error {
  return dr.View( func( tx *bolt.Tx ) error {
    resp := &TocsResponse{}

    if err := tx.Bucket( []byte("DarwinToc") ).ForEach( func( k, v []byte ) error {
      toc := &Toc{}
      if toc.fromBytes( v ) {
        toc.SetSelf( r )
        resp.Toc = append( resp.Toc, toc )
      }
      return nil
    }); err != nil {
     return err
    } else {
      resp.Self = r.Self( r.Context() + "/toc" )
      r.Status( 200 ).Value( resp )
    }

    return nil
  })
}
