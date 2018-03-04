package darwinref

// GetVias makes a batch lookup of one or more ViaResolveRequest's and returns
// a map of matched Via's.
// The result will only contain those entries that were matched.
func (c *DarwinRefClient) GetVias( request map[string]*ViaResolveRequest ) ( map[string]*Via, error ) {

  response := make( map[string]*Via )

  if found, err := c.post( "/via", request, &response ); err != nil {
    return nil, err
  } else if found {
    return response, nil
  } else {
    return nil, nil
  }
}
