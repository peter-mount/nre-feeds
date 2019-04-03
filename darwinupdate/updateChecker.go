package darwinupdate

import (
	"time"
)

func (u *DarwinUpdate) ImportRequiredTimetable(v interface{ TimetableId() string }) bool {
	// Import if no TimetableId
	if v.TimetableId() == "" {
		return true
	}

	// Import if TimetableId is older than the current day
	limit := time.Now().Truncate(24 * time.Hour)
	tid, err := time.Parse("20060102150405", v.TimetableId())
	// Error then force import as tid is invalid
	return err != nil || tid.Before(limit)
}
