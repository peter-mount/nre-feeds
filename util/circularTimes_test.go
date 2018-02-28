package util

import (
//  "encoding/json"
//  "fmt"
//  "github.com/peter-mount/golib/codec"
  "log"
  "sort"
  "testing"
)

func circularTimes_testBool( t *testing.T, m string, e bool, f func() bool ) {
  v := f()
  if v != e {
    t.Errorf( "%s: got %v expected %v", m, v, e )
  }
}

// Test CircularTimes.Compare() works correctly
func TestCircularTimes_Compare( t *testing.T ) {

  a := &CircularTimes{
    Pta: NewPublicTime( "01:02" ),
  }
  a.UpdateTime()

  b := &CircularTimes{
    Ptd: NewPublicTime( "02:03" ),
  }
  b.UpdateTime()

  // a < b
  circularTimes_testBool( t, "a.Compare(b)", true, func() bool {
    return a.Compare( b )
  })

  // b > b
  circularTimes_testBool( t, "b.Compare(a)", false, func() bool {
    return b.Compare( a )
  })

  // a = a & b = b
  circularTimes_testBool( t, "a.Compare(a)", false, func() bool {
    return a.Compare( a )
  })
  circularTimes_testBool( t, "b.Compare(b)", false, func() bool {
    return b.Compare( b )
  })

}

// Test that slices sort correctly
func circularTimes_timesSlice() []*CircularTimes {
  var ary []*CircularTimes
  var circularTimes_times = [...]string {
    // Ensure first element IS NOT the lowest value, see test for why
    "09:50",
    "09:14",
    "09:18",
    "09:55",
    "09:39",
    "09:33",
    "09:37",
    "10:14",
    "09:32",
    "09:40",
    "09:47",
    "09:52",
    "09:25",
  }

  for i,s := range circularTimes_times {
    var v *CircularTimes
    switch i%2 {
      case 0:
        v = &CircularTimes{
          Pta: NewPublicTime( s ),
        }
      case 1:
        v = &CircularTimes{
          Ptd: NewPublicTime( s ),
        }
      case 2:
        v = &CircularTimes{
          Wta: NewWorkingTime( s ),
        }
      case 3:
        v = &CircularTimes{
          Wtd: NewWorkingTime( s ),
        }
    }
    v.UpdateTime()
    ary = append( ary, v )
  }
  return ary
}

func pary( l string, a []*CircularTimes ) {
  var ary []string
  for _, av := range a {
    ary = append( ary, av.Time.String() )
  }
  log.Println( l, ary )
}

func TestCircularTimes_SliceStable( t *testing.T ) {

  a := circularTimes_timesSlice()
  pary( "a", a )

  b := circularTimes_timesSlice()
  pary( "b", b )

  sort.SliceStable( b, func( i, j int ) bool {
    return b[ i ].Compare( b[ j ] )
  } )

  pary( "b", b )

  if len( a ) != len( b ) {
    t.Errorf( "Length incorrect, expected %d got %d", len( a ), len( b ) )
  } else {
    // This will work if both same length. If an element is identical then we
    // havent sorted correctly
    for i, av := range a {
      if av.Equals( b[i] ) {
        t.Errorf( "Slice failed, element %d not changed: %v", i, av.Time.String() )
      }
    }
  }

}
