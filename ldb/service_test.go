package ldb

import (
  "encoding/json"
  "fmt"
  "github.com/peter-mount/nre-feeds/util"
  "log"
  "sort"
  "testing"
)

const serviceJson = "{" +
    "\"rid\":\"%s\"," +
    "\"ssd\":\"20180102\"," +
    "\"location\":{" +
      "\"type\":\"TP\"," +
      "\"tiploc\":\"%s\"," +
      "\"displaytime\":\"%s\"," +
      "\"timetable\":{" +
        "\"time\":\"%s\"," +
        "\"pta\":\"%s\"" +
      "}," +
      "\"forecast\":{" +
        "\"time\":\"%s\"," +
        "\"arr\":{" +
          "\"at\":\"%s\"" +
        "}" +
      "}" +
    "}" +
  "}"

// create new schedule
// rid schedule running id
// tpl tiploc for location
// tm time to use for the arrival
func service_new( t *testing.T, rid, tpl string, tm string ) *Service {
  j := fmt.Sprintf( serviceJson,
    rid,
    tpl,
    tm,
    tm, tm,
    tm, tm,
  )
  v := &Service{}
  if err := json.Unmarshal( []byte(j), v ); err != nil {
    t.Errorf( "Failed to unmarshal json: %v\n\n%s", err, j )
  }
  return v
}

// tests service_new used in the other tests
func TestService_JSON_PARSE( t *testing.T ) {
  a := service_new( t, "12345", "MSTONEE", "01:02" )
  if a == nil {
    t.Errorf( "No schedule returned" )
  }
  if a.RID != "12345" {
    t.Errorf( "Invalid RID" )
  }
}

func service_testBool( t *testing.T, m string, e bool, f func() bool ) {
  v := f()
  if v != e {
    t.Errorf( "%s: got %v expected %v", m, v, e )
  }
}

// Test Service.Compare() works correctly
func TestService_Compare( t *testing.T ) {

  a := service_new( t, "12345", "MSTONEE", "01:02" )
  b := service_new( t, "67890", "MSTONEE", "02:03" )

  // a < b
  service_testBool( t, "a.Compare(b)", true, func() bool {
    return a.Compare( b )
  })

  // b > b
  service_testBool( t, "b.Compare(a)", false, func() bool {
    return b.Compare( a )
  })

  // a = a & b = b
  service_testBool( t, "a.Compare(a)", false, func() bool {
    return a.Compare( a )
  })
  service_testBool( t, "b.Compare(b)", false, func() bool {
    return b.Compare( b )
  })

}

// Test that slices sort correctly
func service_timesSlice( t *testing.T ) []*Service {
  var ary []*Service
  var times = [...]string {
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
    "19:50",
    "19:14",
    "19:18",
    "19:55",
    "19:39",
    "19:33",
    "19:37",
    "20:14",
    "19:32",
    "19:40",
    "19:47",
    "19:52",
    "19:25",
  }

  for _,s := range times {
    ary = append( ary, service_new( t, s, s, s ) )
  }
  return ary
}

func pary( l string, a []*Service ) {
  var ary []string
  for _, av := range a {
    ary = append( ary, av.Location.Time.String() )
  }
  log.Println( l, ary )
}

// Tests we sort Service's correctly.
// Done as the rest service seems to get them wrong at times
func TestService_SliceStable( t *testing.T ) {

  ary := service_timesSlice( t )

  pary( "before sort", ary )

  sort.SliceStable( ary, func( i, j int ) bool {
    return ary[ i ].Compare( ary[ j ] )
  } )

  pary( "after sort", ary )

  var l *util.WorkingTime
  for i, v := range ary {
    if i > 0 && v.Location.Time.Before( l ) {
      t.Errorf( "Element %d not in correct place. Last %v Got %v", i, l.String(), v.Location.Time.String() )
    }
    l = &v.Location.Time
  }

}
