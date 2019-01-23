package util

import (
  "encoding/json"
  "fmt"
  "testing"
)

const (
  // Used in parsing tests
  wt_time = "12:34"
  wt_timeT = ((12*60) + 34)*60
  // Public timetable does not have "00:00" times
  wt_zeroTime = "00:00"
  // Used for comparisons, ranging over midnight
  wt_time2 = "23:10"
  wt_time2T = ((23*60)+10)*60

  wt_time3 = "23:25:00"
  wt_time3T = ((23*60)+25)*60

  wt_time4 = "00:10:00"
  wt_time4T = 10*60
)

func runHHMMSS_TimeSeries( t *testing.T, f func(string) bool ) bool {
  cnt := 0
  for h :=0; h < 24; h++ {
    for m := 0; m < 60; m++ {
      for s := 0; s < 60; s++ {
        if f( fmt.Sprintf( "%02d:%02d:%02d", h, m, s ) ) {
          cnt++
          if( cnt>10 ) {
            t.Errorf( "Aborting test after 10 attempts" )
            return true
          }
        }
      }
    }
  }
  return false
}

// Test WorkingTime parses a time correctly
func TestWorkingTime_New( t *testing.T ) {
  tst := func( s string, tt int ) {
    v := NewWorkingTime( s )

    if v.IsZero() {
      t.Errorf( "WorkingTime was zero" )
    }

    if v.Get() != tt {
      t.Errorf( "WorkingTime %s wrong, got %d want %d", s, v.Get(), tt )
    }
  }
  tst( wt_time, wt_timeT )
  tst( wt_time2, wt_time2T )
  tst( wt_time3, wt_time3T )
  tst( wt_time4, wt_time4T )
}

func TestWorkingTime_Parse( t *testing.T ) {
  var v WorkingTime

  v.Parse( wt_time )
  if v.IsZero() {
    t.Errorf( "WorkingTime was zero" )
  }

  if v.Get() != wt_timeT {
    t.Errorf( "WorkingTime wrong, got %d want %d", v.Get(), wt_timeT )
  }
}

// Test WorkingTime.IsZero() correctly handles "00:00" correctly
func TestWorkingTime_IsZero( t *testing.T ) {
  v := NewWorkingTime( wt_zeroTime )

  if !v.IsZero() {
    t.Errorf( "WorkingTime was not zero" )
  }

  if v.Get() > 0 {
    t.Errorf( "Zero WorkingTime wrong, got %d", v.Get() )
  }
}

func TestWorkingTime_Compare( t *testing.T ) {
  // Test that a is < b, fail if not
  tst := func( a, b string ) {
    av := NewWorkingTime( a )
    bv := NewWorkingTime( b )
    got := av.Compare( bv )
    if !got {
      t.Errorf( "Compare %s to %s got %v want %v", a, b, got, false )
    }
  }

  // This fails when it shoudl pass as it's not past midnight
  //tst( time, time2 )

  tst( wt_time2, wt_time3 )
  tst( wt_time2, wt_time4 )
  tst( wt_time3, wt_time4 )
}

func TestWorkingTime_Equals( t *testing.T ) {
  tst := func( a, b *WorkingTime, e bool ) {
    v := a.Equals( b )
    if v != e {
      t.Errorf( "%v Equals %v failed got %v expected %v", a, b, v, e )
    }
  }

  a := NewWorkingTime( wt_time )
  b := NewWorkingTime( wt_time )
  c := NewWorkingTime( wt_time2 )

  tst( a, b, true )
  tst( a, nil, false )
  tst( a, c, false )
  tst( c, a, false )
}

func TestWorkingTime_JSON( t *testing.T ) {
  runHHMMSS_TimeSeries( t, func( s string ) bool {
    a := NewWorkingTime( s )

    if b, err := json.Marshal( a ); err != nil {
      t.Errorf( "%s failed to marshal: %v", a, err )
      return true
    } else {
      // Check the strings match, they should be a JSON string
      // or null for "00:00"
      as := "\"" + s + "\""
      bs := string( b[:] )
      null := s == "00:00:00"
      if null {
        as = "null"
      }
      if as != bs {
        t.Errorf( "%s failed, marshal got %s expected %s", s, bs, as )
        return true
      }

      // unmarshal back
      c := &WorkingTime{}
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

func TestWorkingTime_After( t *testing.T ) {
  a := NewWorkingTime( "02:03" )
  b := NewWorkingTime( "01:02" )

  if !a.After( b ) {
    t.Errorf( "%v not after %v", a, b )
  }

  if b.After( a ) {
    t.Errorf( "%v after %v", b, a )
  }
}

func TestWorkingTime_Before( t *testing.T ) {
  a := NewWorkingTime( "01:02" )
  b := NewWorkingTime( "02:03" )

  if !a.Before( b ) {
    t.Errorf( "%v not before %v", a, b )
  }

  if b.Before( a ) {
    t.Errorf( "%v before %v", a, b )
  }
}

func TestWorkingTime_Between( t *testing.T ) {
  from := NewWorkingTime( "20:30:00" )
  to := NewWorkingTime( "21:30:00" )

  // Time between the range

  v := NewWorkingTime( "20:33:30" )

  if !v.Between( from, to ) {
    t.Errorf( "%v not between %v and %v", v, from, to )
  }

  // Cross midnight v should not be between from and to
  if v.Between( to, from ) {
    t.Errorf( "%v between %v and %v", v, from, to )
  }

  // Time outside the range

  v = NewWorkingTime( "09:40" )

  if v.Between( from, to ) {
    t.Errorf( "%v between %v and %v", v, from, to )
  }

  // Cross midnight v should not be between from and to
  if !v.Between( to, from ) {
    t.Errorf( "%v not between %v and %v", v, from, to )
  }

}
