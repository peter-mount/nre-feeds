package client

import (
  "github.com/peter-mount/nre-feeds/ldb/service"
)

// GetSchedule returns an active schedule or nil
func (c *DarwinLDBClient) GetService(rid string) (*service.ServiceResult, error) {
  msg := &service.ServiceResult{}
  if found, err := c.get("/service/"+rid, msg); err != nil {
    return nil, err
  } else if found {
    return msg, nil
  } else {
    return nil, nil
  }
}
