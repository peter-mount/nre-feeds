package darwind3

import (
	"github.com/peter-mount/nre-feeds/util"
	"time"
)

// Defines the expected Train order at a platform
type TrainOrder struct {
	Order int `json:"order" xml:"order,attr"`
	// The platform number where the train order applies
	Platform string `json:"plat,omitempty" xml:"plat,attr,omitempty"`
	// This is the TS time from Darwin so we keep a copy of when this struct
	// was sent to us
	Date time.Time `json:"date,omitempty"`
}

// The trainOrder as received from darwin
type trainOrderWrapper struct {
	// The tiploc where the train order applies
	Tiploc string `xml:"tiploc,attr"`
	// The CRS code of the station where the train order applies
	CRS string `xml:"crs,attr"`
	// The platform number where the train order applies
	Platform string `xml:"platform,attr"`
	// The Train orders to set
	Set *trainOrderData `xml:"set"`
	// Clear the current train order
	Clear bool `xml:"clear"`
}

// Defines the sequence of trains making up the train order
type trainOrderData struct {
	// The first train in the train order.
	First *trainOrderItem `xml:"first"`
	// The second train in the train order.
	Second *trainOrderItem `xml:"second"`
	// The third train in the train order.
	Third *trainOrderItem `xml:"third"`
}

// Describes the identifier of a train in the train order
type trainOrderItem struct {
	// For trains in the train order where the train is the Darwin timetable,
	// it will be identified by its RID
	RID string `xml:"rid"`
	// One or more scheduled times to identify the instance of the location in
	// the train schedule for which the train order is set.
	Times util.CircularTimes `xml:"times"`
	// Where a train in the train order is not in the Darwin timetable,
	// a Train ID (headcode) will be supplied
	TrainId string `xml:"trainID"`
}
