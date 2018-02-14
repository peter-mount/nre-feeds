package darwintimetable

// GetJourney returns a Journey by making an HTTP call to a remote instance of
// DarwinTimetable
func (c *DarwinTimetableClient ) GetJourney( rid string ) ( *Journey, error ) {
  journey := &Journey{}
  if found, err := c.get( "/journey/" + rid, journey ); err != nil {
    return nil, err
  } else if found {
    return journey, nil
  } else {
    return nil, nil
  }
}
