// darwind3 handles the real time push port feed
package darwind3

import (
	"github.com/peter-mount/filecache"
	"github.com/peter-mount/nre-feeds/bin"
)

type DarwinD3 struct {
	Timetable    string
	EventManager *DarwinEventManager
	Messages     StationMessages
	FeedStatus   FeedStatus
	Config       *bin.Config
	status       string
	colour       string

	Alarms          *filecache.CacheTable
	Associations    *filecache.CacheTable
	Meta            *filecache.CacheTable
	Schedules       *filecache.CacheTable
	StationMessages *filecache.CacheTable
}

// Init opens a DarwinReference database.
func (r *DarwinD3) Init(em *DarwinEventManager) {
	r.EventManager = em
	r.FeedStatus.d3 = r
	r.Messages.d3 = r
}

func (r *DarwinD3) SetStatus(status, colour string) {
	if status == "" {
		r.status = "OK"
	} else {
		r.status = status
	}
	if colour == "" {
		r.colour = "green"
	} else {
		r.colour = colour
	}
}

func (r *DarwinD3) GetStatus() (string, string) {
	return r.status, r.colour
}
