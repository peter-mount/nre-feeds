package darwintty

import (
	"bytes"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"strings"
	"time"
)

func (s Server) servePlain(r *rest.Rest, result *service.StationResult) error {
	var out bytes.Buffer

	fmt.Fprintf(&out, "Station: %s\n", GetTiploc(result, result.Station[0]))

	// Work out max length of destination and via's
	destLen := 0
	for _, departure := range result.Services {
		loc := departure.Location
		if !loc.IsDestination() {
			destLen = Max(destLen, len(GetDestName(result, departure)))
		}
	}

	fmt1 := fmt.Sprintf("| %%-%d.%ds %%2.2s %%5.5s %%5.5s |\n", destLen, destLen)

	hdr := fmt.Sprintf(fmt.Sprintf("| %%-%d.%ds Pl Deprt Exptd |", destLen, destLen), "Destination")
	// maxLen is destLen + the headers
	maxLen := destLen + 2 + 5 + 5 + 3
	sep := "+-" + strings.Repeat("-", maxLen) + "-+"
	fmt2 := fmt.Sprintf("| %%-%d.%ds |\n", maxLen, maxLen)

	fmt.Fprintln(&out, sep)
	fmt.Fprintln(&out, hdr)
	fmt.Fprintln(&out, sep)

	for _, departure := range result.Services {
		loc := departure.Location
		if !loc.IsDestination() {

			destName := GetDestName(result, departure)

			var supp []string

			if result.Reasons != nil {
				reason := result.Reasons.Late[departure.LateReason.Reason]
				if loc.Cancelled && departure.CancelReason.Reason != 0 {
					reason = result.Reasons.Cancelled[departure.CancelReason.Reason]
				}
				if reason != nil {
					supp = append(supp, telstar.Split(reason.Text, maxLen)...)
				}
			}

			plat := ""
			if !(loc.Forecast.Platform.CISSuppressed || loc.Forecast.Platform.Suppressed) {
				plat = loc.Forecast.Platform.Platform
			}

			fmt.Fprintf(&out, fmt1,
				destName,
				plat,
				loc.Time.String()[:5],
				loc.Forecast.Time.String()[:5],
			)

			for _, s := range supp {
				fmt.Fprintf(&out, fmt2, s)
			}

			fmt.Fprintln(&out, sep)
		}
	}

	fmt.Fprintf(&out, "Generated: %s\n", result.Date.Format(time.RFC3339))

	b := out.Bytes()
	r.ContentType("text/plain").
		Value(b)
	return nil
}
