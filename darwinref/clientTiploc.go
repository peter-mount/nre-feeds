package darwinref

func (c *DarwinRefClient) GetTiplocs( tpl []string ) ( []*Location, error ) {
  ary := make( []*Location, 0 )

  if found, err := c.post( "/tiploc", tpl, &ary ); err != nil {
    return nil, err
  } else if found {
    return ary, nil
  } else {
    return nil, nil
  }
}

func (c *DarwinRefClient) GetTiplocsMapKeys( m map[string]interface{} ) ( []*Location, error ) {
  var tpl []string
  for k, _ := range m {
    tpl = append( tpl, k )
  }
  return c.GetTiplocs( tpl )
}

func (c *DarwinRefClient) GetTiploc( tpl string ) ( *Location, error ) {
  res := &Location{}

  if found, err := c.get( "/tiploc/" + tpl , &res ); err != nil {
    return nil, err
  } else if found {
    return res, nil
  } else {
    return nil, nil
  }
}
