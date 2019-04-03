package util

import (
	"math"
	"strings"
)

const (
	skey_alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	skey_base     = int64(len(skey_alphabet))
	skey_nil      = int64(86400)
)

// Encode int64 to base62 string.
func skey_encode(n int64) string {
	if n == 0 {
		return "0"
	}

	b := make([]byte, 0, 512)
	for n > 0 {
		r := math.Mod(float64(n), float64(skey_base))
		n /= skey_base
		b = append([]byte{skey_alphabet[int(r)]}, b...)
	}
	return string(b)
}

// Encode a working time to base64. Nil or -1 will be set to 86400 or "mty".
func skey_encodeWT(d []string, w *WorkingTime) []string {
	i := skey_nil
	if w != nil {
		i = int64(w.Get())
	}
	if i < 0 || i > 86400 {
		i = skey_nil
	}
	return append(d, skey_encode(i))
}

// GenerateScheduleKey generates a unique key for the given tiploc and working schedule
func GenerateScheduleKey(tpl string, wta *WorkingTime, wtd *WorkingTime, wtp *WorkingTime) string {
	var d []string

	// Just include tpl as it's shorter than encoded
	d = append(d, tpl)
	d = skey_encodeWT(d, wta)
	d = skey_encodeWT(d, wtd)
	d = skey_encodeWT(d, wtp)
	return strings.Join(d, "")
}

// ScheduleKey returns a unique key for this CircularTimes
func (c *CircularTimes) ScheduleKey(tpl string) string {
	return GenerateScheduleKey(tpl, c.Wta, c.Wtd, c.Wtp)
}
