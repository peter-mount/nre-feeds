package client

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/darwindb"
	"strings"
	"time"
)

func (c *DarwinDBClient) GetServices(crs string, date time.Time) (bool, darwindb.StationServices, error) {
	var msg darwindb.StationServices

	found, err := c.get(
		fmt.Sprintf("/services/%s/%s/%02d",
			strings.ToLower(crs),
			date.Format("20060102"),
			date.Hour(),
		),
		&msg)

	return found, msg, err
}
