package client

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/darwind3"
)

// GetJourney returns a Journey by making an HTTP call to a remote instance of
// DarwinTimetable
func (c *DarwinD3Client) GetStationMessage(id uint64) (*darwind3.StationMessage, error) {
	msg := &darwind3.StationMessage{}
	if found, err := c.get(fmt.Sprintf("/message/%d", id), msg); err != nil {
		return nil, err
	} else if found {
		return msg, nil
	} else {
		return nil, nil
	}
}
