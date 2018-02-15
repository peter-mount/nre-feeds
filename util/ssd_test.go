package util

import (
  //"encoding/json"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "testing"
)

func runDate_TimeSeries( t *testing.T, f func(string) bool ) bool {
  for y := 2000; y <= 2040; y++ {
    for m := 1; m <=12; m++ {
      // We'll test 1..28 ignore other days
      for d := 1; d <= 28; d++ {
        s := fmt.Sprintf( "%04d-%02d-%02d", y, m, d )
        if f(s) {
          return true
        }
      }
    }
  }
  return false
}

func TestSSD_Parse( t *testing.T ) {
  runDate_TimeSeries( t, func(s string) bool {
    ssd := &SSD{}
    ssd.Parse( s )
    if ssd.String() != s {
      t.Errorf( "SSD expected %s got %s", s, ssd.String() )
      return true
    }
    return false
  })
}

func TestSSD_ReadWrite( t *testing.T ) {
  runDate_TimeSeries( t, func(s string) bool {
    a := &SSD{}
    a.Parse( s )

    encoder := codec.NewBinaryCodec()
    encoder.Write( a )
    if encoder.Error() != nil {
      t.Errorf( "%s failed to encode: %v", s, encoder.Error() )
      return true
    }

    b := encoder.Bytes()

    c := &SSD{}
    decoder := codec.NewBinaryCodecFrom( b ).Read( c )
    if decoder.Error() != nil {
      t.Errorf( "%s failed to decode: %v", s, decoder.Error() )
      return true
    }

    if !a.Equals( c ) {
      t.Errorf( "%s failed, got %v expected %v", s, c, a )
      return true
    }

    return false
  } )
}
