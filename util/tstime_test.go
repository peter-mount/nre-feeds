package util

import (
  "encoding/json"
  "testing"
)

// Test we Marshal/Unmarshal JSON correctly
func TestTSTime_JSON( t *testing.T ) {

  gett := func( s string ) *WorkingTime {
    if s == "00:00:00" {
      return nil
    }
    return NewWorkingTime( s )
  }

  tst := func( f func( v *TSTime, s string ) ) bool {
    return runHHMMSS_TimeSeries( t, func( s string ) bool {
      v := &TSTime{}
      f( v, s )
      return tstime_test( t, v )
    } )
  }

  tst( func( v *TSTime, s string ) {
    v.ET = gett( s )
  } )
  tst( func( v *TSTime, s string ) {
    v.ETMin = gett( s )
  } )
  tst( func( v *TSTime, s string ) {
    v.WET = gett( s )
  } )
  tst( func( v *TSTime, s string ) {
    v.AT = gett( s )
  } )
}

func tstime_testTime( t *testing.T, s string, a *WorkingTime, b *WorkingTime ) bool {
  if a == nil {
    if b != nil {
      t.Errorf( "%s expected nil, got %v", s, b )
      return true
    }
  } else if a.IsZero() != b.IsZero() {
    t.Errorf( "%s expected zero %v, got %v", s, a.IsZero(), b.IsZero() )
    return true
  } else if !a.Equals( b ) {
    t.Errorf( "%s expected %v, got %v", s, a, b )
    return true
  }
  return false
}

func tstime_test( t *testing.T, a *TSTime ) bool {

  if b, err := json.Marshal( a ); err != nil {
    t.Errorf( "%s failed to marshal: %v", a, err )
    return true
  } else {

    // unmarshal back
    c := &TSTime{}
    if err := json.Unmarshal( b, c ); err != nil {
      t.Errorf( "%s failed to marshal: %v", a, err )
      return true
    }

    return tstime_testTime( t, "ET", a.ET, c.ET ) ||
           tstime_testTime( t, "ETMin", a.ETMin, c.ETMin ) ||
           tstime_testTime( t, "WET", a.WET, c.WET ) ||
           tstime_testTime( t, "AT", a.AT, c.AT )
  }

}
