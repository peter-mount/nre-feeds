package darwinref

import (
  "encoding/json"
  "testing"
)

// Basic test to ensure we can add entries to a map
func TestReasonMap_AddReason( t *testing.T ) {
  var r *ReasonMap = NewReasonMap()

  r.AddReason( &Reason{
    Code: 1,
    Cancelled: false,
  })

  r.AddReason( &Reason{
    Code: 2,
    Cancelled: true,
  })

}

// common code used by JSON tests
func testReasonMap_JSON( t *testing.T, r *ReasonMap, expected string ) {
  if b, err := json.Marshal( r ); err != nil {
    t.Error( err )
  } else {
    s := string(b[:])

    if s != expected {
      t.Errorf( "Invalid json generated\nGot: %s\nExpected: %s", s, expected )
    }
  }
}

// Marshalling an empty map shound return null
func TestReasonMap_JSON_No_entries( t *testing.T ) {
  var r *ReasonMap = NewReasonMap()

  testReasonMap_JSON( t, r, "null" )
}

// Test we only include late entries if no cancellations
func TestReasonMap_JSON_Late( t *testing.T ) {
  var r *ReasonMap = NewReasonMap()

  r.AddReason( &Reason{
    Code: 1,
    Cancelled: false,
  })

  testReasonMap_JSON( t, r,
    "{\"late\":{\"1\":{\"code\":1,\"reasontext\":\"\",\"canc\":false,\"date\":\"0001-01-01T00:00:00Z\",\"self\":\"\"}}}" )
}

// Test we only include cancellations when no late entries
func TestReasonMap_JSON_Cancelled( t *testing.T ) {
  var r *ReasonMap = NewReasonMap()

  r.AddReason( &Reason{
    Code: 2,
    Cancelled: true,
  })

  testReasonMap_JSON( t, r,
    "{\"cancelled\":{\"2\":{\"code\":2,\"reasontext\":\"\",\"canc\":true,\"date\":\"0001-01-01T00:00:00Z\",\"self\":\"\"}}}" )
}

// Test that json includes both maps
func TestReasonMap_JSON_Late_and_Cancelled( t *testing.T ) {
  var r *ReasonMap = NewReasonMap()

  r.AddReason( &Reason{
    Code: 1,
    Cancelled: false,
  })

  r.AddReason( &Reason{
    Code: 2,
    Cancelled: true,
  })

  testReasonMap_JSON( t, r,
    "{\"late\":{\"1\":{\"code\":1,\"reasontext\":\"\",\"canc\":false,\"date\":\"0001-01-01T00:00:00Z\",\"self\":\"\"}},\"cancelled\":{\"2\":{\"code\":2,\"reasontext\":\"\",\"canc\":true,\"date\":\"0001-01-01T00:00:00Z\",\"self\":\"\"}}}" )
}
