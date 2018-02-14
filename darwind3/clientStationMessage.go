package darwind3

import (
  "fmt"
)

// GetJourney returns a Journey by making an HTTP call to a remote instance of
// DarwinTimetable
func (c *DarwinD3Client ) GetStationMessage( id int ) ( *StationMessage, error ) {
  msg := &StationMessage{}
  if found, err := c.get( fmt.Sprintf( "/message/%d", id ), msg ); err != nil {
    return nil, err
  } else if found {
    return msg, nil
  } else {
    return nil, nil
  }
}
