package darwintty

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"io"
	"strings"
)

func (b *Board) Write(w io.Writer) {

	fmt.Fprintf(w, "Live train departures for: %s\n", b.Name)

	// Work out max length of destination and via's
	destLen := len("destination")
	fullLen := 0
	for _, departure := range b.Departures {
		destLen = Max(destLen, len(departure.Destination))
		fullLen = MaxV(fullLen,
			len(departure.LastReport.String()),
			len(departure.TocName()),
			len(departure.Coaches()),
			len(departure.FormationString()),
			len(departure.ToiletStatus()),
		)
	}

	// Expand destLen if fullLen is larger than it and the headers
	headerLen := 2 + 5 + 5 + 3
	if fullLen > (destLen + headerLen) {
		destLen = Max(destLen, fullLen-headerLen)
	}
	maxLen := Max(destLen+headerLen, fullLen)

	fmt1 := fmt.Sprintf("%s %%-%d.%ds %%2.2s %%5.5s %%5.5s %s\n", vertical, destLen, destLen, vertical)

	hdr := fmt.Sprintf(fmt.Sprintf("%s %%-%d.%ds Pl Deprt Exptd %s", vertical, destLen, destLen, vertical), "Destination")

	sep := "%s" + horiz + strings.Repeat(horiz, maxLen) + horiz + "%s\n"
	fmt2 := fmt.Sprintf("%s %%-%d.%ds %s\n", vertical, maxLen, maxLen, vertical)

	fmt.Fprintf(w, sep, topLeft, topRight)
	fmt.Fprintln(w, hdr)

	switch len(b.Departures) {
	case 0:
	default:
		for _, departure := range b.Departures {
			fmt.Fprintf(w, sep, midLeft, midRight)

			var supp []string

			if departure.Reason != "" {
				supp = append(supp, telstar.Split(departure.Reason, maxLen)...)
			}

			supp = append(supp,
				departure.TocName(),
				departure.Coaches(),
				departure.FormationString(),
				departure.ToiletStatus(),
				departure.LastReport.String(),
			)

			fmt.Fprintf(w, fmt1, departure.Destination, departure.Plat, departure.Depart, departure.Expected)

			for _, s := range supp {
				if s != "" {
					fmt.Fprintf(w, fmt2, s)
				}
			}
		}
	}
	fmt.Fprintf(w, sep, bottomLeft, bottomRight)
}
