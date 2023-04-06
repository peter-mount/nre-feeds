package darwintty

import (
	"github.com/peter-mount/nre-feeds/ldb"
	"github.com/peter-mount/nre-feeds/ldb/service"
)

// GetTiploc returns the Tiploc name or the tiploc if not present
func GetTiploc(r *service.StationResult, tpl string) string {
	entry, _ := r.Tiplocs.Get(tpl)
	if entry != nil {
		return entry.Name
	}
	return tpl
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxV(a int, b ...int) int {
	for _, e := range b {
		a = Max(a, e)
	}
	return a
}

func GetDestName(r *service.StationResult, departure ldb.Service) string {
	loc := departure.Location
	dest := departure.Dest
	destName := dest.Tiploc
	if loc.FalseDestination != "" {
		destName = loc.FalseDestination
	}
	destName = GetTiploc(r, destName)

	if v, exists := r.Via[departure.RID]; exists {
		destName = destName + " " + v.Text
	}

	return destName
}
