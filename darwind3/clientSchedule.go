package darwind3

// GetSchedule returns an active Schedule or nil
func (c *DarwinD3Client ) GetSchedule( rid string ) ( *Schedule, error ) {
  msg := &Schedule{}
  if found, err := c.get( "/schedule/" + rid, msg ); err != nil {
    return nil, err
  } else if found {
    return msg, nil
  } else {
    return nil, nil
  }
}
