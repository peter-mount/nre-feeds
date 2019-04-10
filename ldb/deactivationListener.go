package ldb

import (
	"github.com/peter-mount/nre-feeds/darwind3"
)

// deactivationListener removes Services when a schedule is deactivated
func (d *LDB) deactivationListener(e *darwind3.DarwinEvent) {
	if e.RID != "" {
		d.RemoveSchedule(e.RID)
	}
}
