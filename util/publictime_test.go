package util

import (
  "encoding/json"
  "fmt"
  "github.com/peter-mount/golib/codec"
//  "log"
  "testing"
)

const (
  // Used in parsing tests
  pt_time = "12:34"
  pt_timeT = (12*60) + 34
  // Public timetable does not have "00:00" times
  pt_zeroTime = "00:00"
  // Used for comparisons, ranging over midnight
  pt_time2 = "23:10"
  pt_time2T = (23*60)+10

  pt_time3 = "23:25"
  pt_time3T = (23*60)+25

  pt_time4 = "00:10"
  pt_time4T = 10
)

func runHHMM_TimeSeries( t *testing.T, f func(string) bool ) {
  cnt := 0
  for h :=0; h < 24; h++ {
    for m := 0; m < 60; m ++ {
      if f( fmt.Sprintf( "%02d:%02d", h, m ) ) {
        cnt++
        if( cnt>10 ) {
          t.Errorf( "Aborting test after 10 attempts" )
          return
        }
      }
    }
  }
}

// Test PublicTime parses a time correctly
func TestPublicTime_New( t *testing.T ) {
  tst := func( s string, tt int ) {
    v := NewPublicTime( s )

    if v.IsZero() {
      t.Errorf( "PublicTime was zero" )
    }

    if v.Get() != tt {
      t.Errorf( "PublicTime %s wrong, got %d want %d", s, v.Get(), tt )
    }
  }
  tst( pt_time, pt_timeT )
  tst( pt_time2, pt_time2T )
  tst( pt_time3, pt_time3T )
  tst( pt_time4, pt_time4T )
}

func TestPublicTime_Parse( t *testing.T ) {
  var v PublicTime

  v.Parse( pt_time )
  if v.IsZero() {
    t.Errorf( "PublicTime was zero" )
  }

  if v.Get() != pt_timeT {
    t.Errorf( "PublicTime wrong, got %d want %d", v.Get(), pt_timeT )
  }
}

// Test PublicTime.IsZero() correctly handles "00:00" correctly
func TestPublicTime_IsZero( t *testing.T ) {
  v := NewPublicTime( pt_zeroTime )

  if !v.IsZero() {
    t.Errorf( "PublicTime was not zero" )
  }

  if v.Get() > 0 {
    t.Errorf( "Zero PublicTime wrong, got %d", v.Get() )
  }
}

func TestPublicTime_Compare( t *testing.T ) {
  // Test that a is < b, fail if not
  tst := func( a, b string ) {
    av := NewPublicTime( a )
    bv := NewPublicTime( b )
    got := av.Compare( bv )
    if !got {
      t.Errorf( "Compare %s to %s got %v want %v", a, b, got, false )
    }
  }

  // This fails when it shoudl pass as it's not past midnight
  //tst( time, time2 )

  tst( pt_time2, pt_time3 )
  tst( pt_time2, pt_time4 )
  tst( pt_time3, pt_time4 )
}

func TestPublicTime_Equals( t *testing.T ) {
  tst := func( a, b *PublicTime, e bool ) {
    v := a.Equals( b )
    if v != e {
      t.Errorf( "%v Equals %v failed got %v expected %v", a, b, v, e )
    }
  }

  a := NewPublicTime( pt_time )
  b := NewPublicTime( pt_time )
  c := NewPublicTime( pt_time2 )

  tst( a, b, true )
  tst( a, nil, false )
  tst( a, c, false )
  tst( c, a, false )
}

func TestPublicTime_ReadWrite( t *testing.T ) {
  runHHMM_TimeSeries( t, func( s string ) bool {
    a := NewPublicTime( s )

    encoder := codec.NewBinaryCodec()
    encoder.Write( a )
    if encoder.Error() != nil {
      t.Errorf( "%s failed to encode: %v", s, encoder.Error() )
      return true
    }

    b := encoder.Bytes()

    c := &PublicTime{}
    decoder := codec.NewBinaryCodecFrom( b ).Read( c )
    if decoder.Error() != nil {
      t.Errorf( "%s failed to decode: %v", s, decoder.Error() )
      return true
    }

    if a.IsZero() != c.IsZero() {
      t.Errorf( "%s failed isZero, got %v expected %v", s, c.Get(), a.Get() )
      return true
    } else if !a.Equals( c ) {
      t.Errorf( "%s failed, got %v expected %v", s, c.Get(), a.Get() )
      return true
    }

    return false
  } )
}

func TestPublicTime_JSON( t *testing.T ) {
  runHHMM_TimeSeries( t, func( s string ) bool {
    a := NewPublicTime( s )

    if b, err := json.Marshal( a ); err != nil {
      t.Errorf( "%s failed to marshal: %v", a, err )
      return true
    } else {
      // Check the strings match, they should be a JSON string
      // or null for "00:00"
      as := "\"" + s + "\""
      bs := string( b[:] )
      null := s == "00:00"
      if null {
        as = "null"
      }
      if as != bs {
        t.Errorf( "%s failed, marshal got %s expected %s", s, bs, as )
        return true
      }

      // unmarshal back
      c := &PublicTime{}
      if err := json.Unmarshal( b, c ); err != nil {
        t.Errorf( "%s failed to marshal: %v", a, err )
        return true
      }

      if null && !c.IsZero() {
        t.Errorf( "%s failed, got %v expected true", s, c.IsZero() )
        return true
      }

      if !a.Equals( c ) {
        t.Errorf( "%s failed, got %v expected %v", s, c, a )
        return true
      }
    }

    return false
  } )
}
