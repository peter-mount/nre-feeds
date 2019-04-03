// darwind3 handles the real time push port feed
package darwind3

import "github.com/peter-mount/nre-feeds/bin"

type DarwinD3 struct {
	Timetable    string
	EventManager *DarwinEventManager
	cache        cache
	Messages     *StationMessages
	FeedStatus   FeedStatus
	Config       *bin.Config
}

// OpenDB opens a DarwinReference database.
func (r *DarwinD3) OpenDB(dbFile string, em *DarwinEventManager) error {
	r.FeedStatus.d3 = r
	r.EventManager = em
	r.Messages = NewStationMessages(dbFile)

	return r.cache.initCache(dbFile)
}
