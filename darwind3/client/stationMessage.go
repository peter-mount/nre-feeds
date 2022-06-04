package client

import (
  "fmt"
  "github.com/peter-mount/nre-feeds/darwind3"
)

// GetStationMessage returns a specific station message
func (c *DarwinD3Client) GetStationMessage(id int64) (*darwind3.StationMessage, error) {
  msg := &darwind3.StationMessage{}
  if found, err := c.get(fmt.Sprintf("/message/%d", id), msg); err != nil {
    return nil, err
  } else if found {
    return msg, nil
  } else {
    return nil, nil
  }
}

// GetStationMessages returns all Station messages
func (c *DarwinD3Client) GetStationMessages() ([]darwind3.StationMessage, error) {
  var msg []darwind3.StationMessage

  _, err := c.get("/messages", &msg)
  return msg, err
}
