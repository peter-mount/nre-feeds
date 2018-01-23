package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "time"
)

// Via text
type Via struct {
  XMLName     xml.Name  `json:"-" xml:"Via"`
  At      string        `json:"at" xml:"at,attr"`
  Dest    string        `json:"dest" xml:"dest,attr"`
  Loc1    string        `json:"loc1" xml:"loc1,attr"`
  Loc2    string        `json:"loc2,omitempty" xml:"loc2,attr,omitempty"`
  Text    string        `json:"text" xml:"viatext,attr"`
  // Date entry was inserted into the database
  Date        time.Time `json:"date" xml:"date,attr"`
  // URL to this entity
  Self        string    `json:"self" xml:"self,attr,omitempty"`
}

// Are two Via's equal
func (v *Via) Equals( o *Via ) bool {
  if o == nil {
    return false
  }
  return v.At == o.At && v.Dest == o.Dest && v.Loc1 == o.Loc1 && v.Loc2 == o.Loc2
}

func (v *Via) Write( c *codec.BinaryCodec ) {
  c.WriteString( v.At ).
    WriteString( v.Dest ).
    WriteString( v.Loc1 ).
    WriteString( v.Loc2 ).
    WriteString( v.Text ).
    WriteTime( v.Date )
}

func (v *Via) Read( c *codec.BinaryCodec ) {
  c.ReadString( &v.At ).
    ReadString( &v.Dest ).
    ReadString( &v.Loc1 ).
    ReadString( &v.Loc2 ).
    ReadString( &v.Text ).
    ReadTime( &v.Date )
}

// SetSelf sets the Self field to match this request
func (v *Via) SetSelf( r *rest.Rest ) {
  if v.Loc2 == "" {
    v.Self = r.Self( fmt.Sprintf( "%s/via/%s/%s/%s", r.Context(), v.At, v.Dest, v.Loc1 ) )
  } else {
    v.Self = r.Self( fmt.Sprintf( "%s/via/%s/%s/%s/%s", r.Context(), v.At, v.Dest, v.Loc1, v.Loc2 ) )
  }
}

// Key the unique key for this entry
func (v *Via) key() string {
  return fmt.Sprintf( "%s %s %s %s", v.At, v.Dest, v.Loc1, v.Loc2 )
}

func (v *Via) String() string {
  return "Via[At=" + v.At +", Dest=" + v.Dest +", Loc1=" + v.Loc1 +", Loc2=" + v.Loc2 +", Text=" + v.Text + "]"
}

// GetToc returns details of a TOC
func (r *DarwinReference) GetVia( tx *bolt.Tx, at string, dest string, loc1 string, loc2 string ) ( *Via, bool ) {
  loc, exists := r.GetViaBucket( tx.Bucket( []byte("DarwinVia") ), at, dest, loc1, loc2 )
  return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getVia( at string, dest string, loc1 string, loc2 string ) ( *Via, bool ) {
  loc, exists := r.GetViaBucket( r.via, at, dest, loc1, loc2 )
  return loc, exists
}

func (t *Via) fromBytes( b []byte ) bool {
  if b != nil {
    codec.NewBinaryCodecFrom( b ).Read( t )
  }
  return t.At != ""
}

func (r *DarwinReference) GetViaBucket( bucket *bolt.Bucket, at string, dest string, loc1 string, loc2 string ) ( *Via, bool ) {
  key := fmt.Sprintf( "%s %s %s %s", at, dest, loc1, loc2 )
  b := bucket.Get( []byte( key ) )

  if b != nil {
    var via *Via = &Via{}
    if via.fromBytes( b ) {
      return via, true
    }
  }

  return nil, false
}

func (r *DarwinReference) addVia( via *Via ) ( error, bool ) {
  // Update only if it does not exist or is different
  if old, exists := r.getVia( via.At, via.Dest, via.Loc1, via.Loc2 ); !exists || !via.Equals( old ) {
    via.Date = time.Now()
    codec := codec.NewBinaryCodec()
    codec.Write( via )
    if codec.Error() != nil {
      return codec.Error(), false
    }

    if err := r.via.Put( []byte( via.key() ), codec.Bytes() ); err != nil {
      return err, false
    }

    return nil, true
  }

  return nil, false
}
