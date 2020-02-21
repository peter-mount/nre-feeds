package client

import (
  "github.com/peter-mount/nre-feeds/darwinref"
)

func (c *DarwinRefClient) Search(query string) ([]*darwinref.SearchResult, error) {
  var res []*darwinref.SearchResult

  if found, err := c.get("/search/"+query, &res); err != nil {
    return nil, err
  } else if found {
    return res, nil
  } else {
    return nil, nil
  }
}
