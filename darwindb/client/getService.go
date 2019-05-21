package client

import "github.com/peter-mount/nre-feeds/darwindb"

func (c *DarwinDBClient) GetService(rid string) (bool, darwindb.ServiceDetail, error) {
	msg := darwindb.ServiceDetail{}

	found, err := c.get("/service/"+rid, &msg)

	return found, msg, err
}
