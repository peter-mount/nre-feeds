package strings

import (
	"sort"
	"strings"
)

func SortLower(a []string) {
	sort.SliceStable(a, func(i, j int) bool {
		return strings.ToLower(a[i]) < strings.ToLower(a[j])
	})
}
