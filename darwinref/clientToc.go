package darwinref

// GetToc retrieve a Toc by its code
func (c *DarwinRefClient) GetToc( toc string ) ( *Toc, error ) {
  res := &Toc{};

    if found, err := c.get( "/toc/" + toc, &res ); err != nil {
      return nil, err
    } else if found {
      return res, nil
    } else {
      return nil, nil
    }
}

// AddToc adds a Toc to a TocMap
func (c *DarwinRefClient) AddToc( m *TocMap, toc string ) {
  if toc != "" {
    if _, exists := m.Get( toc ); !exists {
      if t, _ := c.GetToc( toc ); t != nil {
        m.Add( t )
      }
    }
  }
}
