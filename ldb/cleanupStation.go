package ldb

// Cleans up a station removing old schedules
func (s *Station) cleanup() {
	/* TODO reimplement station cleanup
	now := time.Now()
	day := now.Add(-2 * time.Hour)

	s.Update(func() error {
		for rid, service := range s.Services {
			if service.Timestamp().Before(day) {
				statistics.Incr("ldb.clean")
				delete(s.Services, rid)
			}
		}
		return nil
	})

	*/
}

// Cleanup removes any old schedules still in memory for each station
/* TODO reimplement Cleanup
func (st *Stations) Cleanup() {
	// Get set of current stations
	var crs []*Station
	st.Update(func() error {
		for _, s := range st.crs {
			crs = append(crs, s)
		}
		return nil
	})

	// Cleanup each one
	for _, s := range crs {
		s.cleanup()
	}
}

*/
