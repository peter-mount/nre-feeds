package client

import (
  "github.com/peter-mount/nre-feeds/darwinref"
)

// GetStations returns all Location's with a CRS code
func (c *DarwinRefClient) GetStations() ( []*darwinref.Location, error ) {
  ary := make( []*darwinref.Location, 0 )

  if found, err := c.get( "/crs", &ary ); err != nil {
    return nil, err
  } else if found {
    return ary, nil
  } else {
    return nil, nil
  }
}

func (c *DarwinRefClient) GetCrs( crs string ) ( *darwinref.CrsResponse, error ) {
  res := &darwinref.CrsResponse{}

  if found, err := c.get( "/crs/" + crs , &res ); err != nil {
    return nil, err
  } else if found {
    return res, nil
  } else {
    return nil, nil
  }
}
