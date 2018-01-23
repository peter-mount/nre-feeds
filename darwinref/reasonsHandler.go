package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
  "sort"
  "strconv"
)

func (dr *DarwinReference) ReasonCancelHandler( r *rest.Rest ) error {
  return dr.reasonHandler( true, r )
}

func (dr *DarwinReference) ReasonLateHandler( r *rest.Rest ) error {
  return dr.reasonHandler( false, r )
}

func (dr *DarwinReference) reasonHandler( cancelled bool, r *rest.Rest ) error {
  id, err := strconv.Atoi( r.Var( "id" ) )
  if err != nil {
    return err
  }

  return dr.View( func( tx *bolt.Tx ) error {
    var reason *Reason
    var exists bool

    if cancelled {
      reason, exists = dr.GetCancellationReason( tx, id )
    } else {
      reason, exists = dr.GetLateReason( tx, id )
    }

    if exists {
      reason.SetSelf( r )
      r.Status( 200 ).Value( reason )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}

func (dr *DarwinReference) AllReasonCancelHandler( r *rest.Rest ) error {
  return dr.allReasonHandler( "DarwinCancelReason", "/reason/cancelled", r )
}

func (dr *DarwinReference) AllReasonLateHandler( r *rest.Rest ) error {
  return dr.allReasonHandler( "DarwinLateReason", "/reason/late", r )
}

type ReasonsResponse struct {
  XMLName     xml.Name  `json:"-" xml:"reasons"`
  Reasons  []*Reason    `json:"reasons,omitempty" xml:"Reason"`
  Self        string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func (dr *DarwinReference) allReasonHandler( bname string, prefix string, r *rest.Rest ) error {
  return dr.View( func( tx *bolt.Tx ) error {
    resp := &ReasonsResponse{}

    if err := tx.Bucket( []byte( bname ) ).ForEach( func( k, v []byte ) error {
      reason := &Reason{}
      if reason.fromBytes( v ) {
        reason.SetSelf( r )
        resp.Reasons = append( resp.Reasons, reason )
      }
      return nil
    }); err != nil {
     return err
    } else {

      // Sort result by code
      sort.SliceStable( resp.Reasons, func( a, b int ) bool {
        ra := resp.Reasons[ a ]
        rb := resp.Reasons[ b ]
        return ra.Code < rb.Code
      })

      resp.Self = r.Self( r.Context() + prefix )
      r.Status( 200 ).Value( resp )
    }

    return nil
  })
}
