package util

import (
  "fmt"
  "testing"
  "time"
)

func testPublicTime_Schedule1() []*PublicTime {
  var ary []*PublicTime
  for i := 10; i<15; i++ {
    t := &PublicTime{}
    t.Parse( fmt.Sprintf( "%02d:%02d", i, i+5) )
    ary = append( ary, t)
  }
  return ary
}

func testPublicTime_Schedule2() []*PublicTime {
  var ary []*PublicTime
  for i := 1; i<10; i++ {
    t := &PublicTime{}
    o := i + 20
    if i > 0 && i < 5 {
      o = i - 5
    }
    t.Parse( fmt.Sprintf( "%02d:%02d", i + o, i+5) )
    ary = append( ary, t)
  }
  return ary
}

// Test Time returns a correct value
func TestPublicTime_TrainTime( ts *testing.T ) {
  start := time.Date( 2018, time.March, 26, 0, 0, 0, 0, London() )

  ary := testPublicTime_Schedule1()
  var times []time.Time
  var first time.Time
  for i, pt := range ary {
    if i == 0 {
      first = pt.Time( start )
      times = append( times, first )
    } else {
      times = append( times, pt.TrainTime( first ) )
    }
  }

  for i, t := range times {
    if i > 0 && !first.Before( t ) {
      ts.Errorf( "Entry %d not in sequence", i )
    }
    first = t
  }
}

// Test Time returns a correct value
func TestPublicTime_TrainTime_Midnight( ts *testing.T ) {
  start := time.Date( 2018, time.March, 26, 0, 0, 0, 0, London() )

  ary := testPublicTime_Schedule2()
  var times []time.Time
  var first time.Time
  for i, pt := range ary {
    if i == 0 {
      first = pt.Time( start )
      times = append( times, first )
    } else {
      times = append( times, pt.TrainTime( first ) )
    }
  }

  for i, t := range times {
    if i > 0 && !first.Before( t ) {
      ts.Errorf( "Entry %d not in sequence", i )
    }
    first = t
  }
}
