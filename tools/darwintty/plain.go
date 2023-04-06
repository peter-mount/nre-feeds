package darwintty

import (
	"bytes"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"strings"
)

func (s Server) servePlain(r *rest.Rest, b *Board) error {

	var out bytes.Buffer

	fmt.Fprintf(&out, "Station: %s\n", b.Name)

	// Work out max length of destination and via's
	destLen := 0
	for _, departure := range b.Departures {
		destLen = Max(destLen, len(departure.Destination))
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

	for _, departure := range b.Departures {

		var supp []string

		if departure.Reason != "" {
			supp = append(supp, telstar.Split(departure.Reason, maxLen)...)
		}

		if departure.Toc != "" || departure.Length > 0 {
			s := ""
			if departure.Toc != "" {
				s = s + departure.Toc + " service "
			}
			if departure.Length > 0 {
				s = s + fmt.Sprintf("Formed of %d coaches", departure.Length)
			}
			supp = append(supp, s)
		}

		if departure.LastReport.Location != "" {
			s := ""
			switch {
			case departure.LastReport.Departed:
				s = " departing"
			case departure.LastReport.At:
				s = " at"
			}
			supp = append(supp, fmt.Sprintf("Last seen%s %s at %s",
				s,
				departure.LastReport.Location,
				departure.LastReport.Time,
			))
		}

		fmt.Fprintf(&out, fmt1, departure.Destination, departure.Plat, departure.Depart, departure.Expected)

		for _, s := range supp {
			fmt.Fprintf(&out, fmt2, s)
		}

		fmt.Fprintln(&out, sep)
	}

	r.ContentType("text/plain").
		Value(out.Bytes())
	return nil
}
