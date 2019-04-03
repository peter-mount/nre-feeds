package client

import (
	"github.com/peter-mount/nre-feeds/darwinref"
)

func (c *DarwinRefClient) GetTiplocs(tpl []string) ([]*darwinref.Location, error) {
	ary := make([]*darwinref.Location, 0)

	if found, err := c.post("/tiploc", tpl, &ary); err != nil {
		return nil, err
	} else if found {
		return ary, nil
	} else {
		return nil, nil
	}
}

func (c *DarwinRefClient) GetTiplocsMapKeys(m map[string]interface{}) ([]*darwinref.Location, error) {
	var tpl []string
	for k, _ := range m {
		tpl = append(tpl, k)
	}
	return c.GetTiplocs(tpl)
}

func (c *DarwinRefClient) GetTiploc(tpl string) (*darwinref.Location, error) {
	res := &darwinref.Location{}

	if found, err := c.get("/tiploc/"+tpl, &res); err != nil {
		return nil, err
	} else if found {
		return res, nil
	} else {
		return nil, nil
	}
}
