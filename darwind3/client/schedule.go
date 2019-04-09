package client

import (
	"github.com/peter-mount/nre-feeds/darwind3"
)

// GetSchedule returns an active schedule or nil
func (c *DarwinD3Client) GetSchedule(rid string) (*darwind3.Schedule, error) {
	msg := &darwind3.Schedule{}
	if found, err := c.get("/schedule/"+rid, msg); err != nil {
		return nil, err
	} else if found {
		return msg, nil
	} else {
		return nil, nil
	}
}
