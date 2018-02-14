package darwinref

// GetStations returns all Location's with a CRS code
func (c *DarwinRefClient) GetStations() ( []*Location, error ) {
  ary := make( []*Location, 0 )

  if found, err := c.get( "/crs", &ary ); err != nil {
    return nil, err
  } else if found {
    return ary, nil
  } else {
    return nil, nil
  }
}
