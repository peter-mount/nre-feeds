package client

import (
  "github.com/peter-mount/nre-feeds/ldb/service"
)

// GetSchedule returns an active schedule or nil
func (c *DarwinLDBClient) GetSchedule(crs string) (*service.StationResult, error) {
  msg := &service.StationResult{}
  if found, err := c.get("/boards/"+crs, msg); err != nil {
    return nil, err
  } else if found {
    return msg, nil
  } else {
    return nil, nil
  }
}
