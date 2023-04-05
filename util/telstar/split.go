package telstar

import (
	"strings"
)

func Split(s string, l int) []string {
	s = strings.TrimSpace(s)

	var a []string

	for len(s) > l {
		i := strings.LastIndex(s[:l], " ")
		if i > 0 {
			a = append(a, strings.TrimSpace(s[:i]))
			s = strings.TrimSpace(s[i:])
		} else {
			a = append(a, strings.TrimSpace(s[:l]))
			s = strings.TrimSpace(s[l:])
		}
	}

	s = strings.TrimSpace(s)
	if s != "" {
		a = append(a, s)
	}

	return a
}
