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

func testWorkingTime_Schedule1() []*WorkingTime {
  var ary []*WorkingTime
  for i := 10; i<15; i++ {
    t := &WorkingTime{}
    t.Parse( fmt.Sprintf( "%02d:%02d:00", i, i+5) )
    ary = append( ary, t)
  }
  return ary
}

func testWorkingTime_Schedule2() []*WorkingTime {
  var ary []*WorkingTime
  for i := 1; i<10; i++ {
    t := &WorkingTime{}
    o := i + 20
    if i > 0 && i < 5 {
      o = i - 5
    }
    t.Parse( fmt.Sprintf( "%02d:%02d:00", i + o, i+5) )
    ary = append( ary, t)
  }
  return ary
}

func testGeneratePublicTimes( start time.Time, ary []*PublicTime ) []time.Time {
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
  return times
}

func testGenerateWorkingTimes( start time.Time, ary []*WorkingTime ) []time.Time {
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
  return times
}

func testTimesInSequence( ts *testing.T, times []time.Time ) {
  var first time.Time
  for i, t := range times {
    if i > 0 && !first.Before( t ) {
      ts.Errorf( "Entry %d not in sequence", i )
    }
    first = t
  }
}

func testDay( f func( int, time.Month) ) {

  // This tests across the clocks going forward on Mar 25 2018
  for day := 23; day <= 26; day ++ {
    f( day, time.March )
  }

  // This tests across the clocks going back on Oct 28 2018
  for day := 27; day <= 29; day ++ {
    f( day, time.October )
  }

}

// Test Time returns a correct value
func TestPublicTime_TrainTime( ts *testing.T ) {
  testDay( func( day int, month time.Month ) {
    start := time.Date( 2018, month, day, 0, 0, 0, 0, London() )

    ary := testPublicTime_Schedule1()

    times := testGeneratePublicTimes( start, ary )

    testTimesInSequence( ts, times )
  } )
}

// Test Time returns a correct value
func TestPublicTime_TrainTime_Midnight( ts *testing.T ) {
  testDay( func( day int, month time.Month ) {
    start := time.Date( 2018, month, day, 0, 0, 0, 0, London() )

    ary := testPublicTime_Schedule2()

    times := testGeneratePublicTimes( start, ary )

    testTimesInSequence( ts, times )
  } )
}

// Test Time returns a correct value
func TestWorkingTime_TrainTime( ts *testing.T ) {
  testDay( func( day int, month time.Month ) {
    start := time.Date( 2018, month, day, 0, 0, 0, 0, London() )

    ary := testWorkingTime_Schedule1()

    times := testGenerateWorkingTimes( start, ary )

    testTimesInSequence( ts, times )
  } )
}

// Test Time returns a correct value
func TestWorkingTime_TrainTime_Midnight( ts *testing.T ) {
  testDay( func( day int, month time.Month ) {
    start := time.Date( 2018, month, day, 0, 0, 0, 0, London() )

    ary := testWorkingTime_Schedule2()

    times := testGenerateWorkingTimes( start, ary )

    testTimesInSequence( ts, times )
  } )
}
