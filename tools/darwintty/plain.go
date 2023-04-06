package darwintty

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"time"
)

func (b *Board) Write(w render.Builder) {

	w = w.Printf("Live train departures for: %s", b.Name).
		NewLine()

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

	fmt1 := fmt.Sprintf(" %%-%d.%ds %%2.2s %%5.5s %%5.5s ", destLen, destLen)

	fmt2 := fmt.Sprintf(" %%-%d.%ds ", maxLen, maxLen)

	w = w.Print(topLeft).
		Repeat(horiz, maxLen+2).
		Print(topRight).
		NewLine().
		Print(vertical).
		White().
		Printf(fmt.Sprintf(" %%-%d.%ds Pl Deprt Exptd ", destLen, destLen), "Destination").
		End().
		Print(vertical).
		NewLine()

	switch len(b.Departures) {
	case 0:
	default:
		for _, departure := range b.Departures {
			w = w.Print(midLeft).Repeat(horiz, maxLen+2).Print(midRight).
				NewLine()

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

			w = w.Print(vertical).
				White().
				Printf(fmt1, departure.Destination, departure.Plat, departure.Depart, departure.Expected).
				End().
				Print(vertical).
				NewLine()

			for _, s := range supp {
				if s != "" {
					w = w.Print(vertical).
						Yellow().
						Printf(fmt2, s).
						End().
						Print(vertical).
						NewLine()
				}
			}
		}
	}
	w = w.Print(bottomLeft).Repeat(horiz, maxLen+2).Print(bottomRight).
		NewLine().
		NewLine().
		Printf("Generated: %s", b.Date.Format(time.RFC3339)).
		NewLine()
}
