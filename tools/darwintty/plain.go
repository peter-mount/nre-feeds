package darwintty

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"time"
)

func (b *Board) Write(w render.Builder) {

	w = w.Printf("Live train departures at %s", b.Date.Format(time.RFC3339)).
		NewLine().
		NewLine()

	// Work out max length of destination and via's
	destLen := len("destination")
	fullLen := 0
	for _, departure := range b.Departures {
		destLen = Max(destLen, len(departure.Destination))
		fullLen = MaxV(fullLen,
			len(b.Name)+8, // Station Name Banner
			len(departure.LastReport.String()),
			len(departure.TocName()),
			len(departure.Coaches()),
			len(departure.FormationString()),
		)
		for _, s := range departure.ToiletStatus() {
			fullLen = Max(fullLen, len(s))
		}
	}

	// Expand destLen if fullLen is larger than it and the headers
	headerLen := 2 + 8 + 8 + 3
	if fullLen > (destLen + headerLen) {
		destLen = Max(destLen, fullLen-headerLen)
	}
	maxLen := Max(destLen+headerLen, fullLen)

	fmt1 := fmt.Sprintf(" %%-%d.%ds ", destLen, destLen)
	fmt2 := fmt.Sprintf(" %%-%d.%ds ", maxLen, maxLen)

	sl0 := len(b.Name)
	sl1 := (maxLen - sl0) >> 1

	w = w.Repeat(" ", sl1).
		Print(topLeft).Repeat(horiz, sl0+2).Print(topRight).
		NewLine().
		Repeat(" ", sl1).
		Print(vertical).Printf(" %s ", b.Name).Print(vertical).
		NewLine()

	w = w.Print(topLeft).
		Repeat(horiz, sl1-1).
		Print(topUpper).
		Repeat(horiz, sl0+2).
		Print(topUpper).
		Repeat(horiz, maxLen-sl1-sl0-1).
		//Repeat(horiz, maxLen+2).
		Print(topRight).
		NewLine()

	w = w.Print(vertical).
		White().
		Printf(fmt.Sprintf(" %%-%d.%ds Pl  Depart  Expected ", destLen, destLen), "Destination").
		End().
		Print(vertical).
		NewLine()

	switch len(b.Departures) {
	case 0:
	default:
		for _, departure := range b.Departures {
			w = w.Print(midLeft).Repeat(horiz, maxLen+2).Print(midRight).
				NewLine()

			w = w.Print(vertical)
			if departure.Cancelled {
				w = w.Red().
					Printf(fmt1, departure.Destination).
					Print("        Cancelled    ")
			} else {
				w = w.White().
					Printf(fmt1, departure.Destination).
					Printf("%2.2s %8.8s ", departure.Plat, departure.Depart)
				if departure.Delayed {
					w = w.End().Red().Print(" Delayed ")
				} else {
					if departure.Depart != departure.Expected {
						w = w.End().
							Yellow().
							Printf("%8.8s ", departure.Expected)
					} else {
						w = w.End().
							Green().
							Print(" On Time ")
					}
				}
			}
			w = w.End().
				Print(vertical).
				NewLine()

			w = renderIf(w, departure.TocName(), fmt2, green)

			if departure.Reason != "" {
				h := yellow
				if departure.Cancelled {
					h = red
				}
				for _, s := range telstar.Split(departure.Reason, maxLen) {
					w = renderIf(w, s, fmt2, h)
				}
			}

			if !departure.Cancelled {
				w = renderIf(w, departure.Coaches(), fmt2, green)
				w = renderIf(w, departure.FormationString(), fmt2, green)
				for _, s := range departure.ToiletStatus() {
					w = renderIf(w, s, fmt2, green)
				}
				w = renderIf(w, departure.LastReport.String(), fmt2, white)
				if departure.Delay > 0 && departure.Delay < 300 {
					w = renderIf(w,
						fmt.Sprintf("Delayed by %d minutes", departure.Delay),
						fmt2, yellow)
				}
			}

		}
	}
	w = w.Print(bottomLeft).Repeat(horiz, maxLen+2).Print(bottomRight).
		NewLine().
		NewLine().
		Print("Data provided by ").
		Link("https://departureboards.mobi").
		NewLine().
		Link("https://area51.dev").
		Print(" and National Rail Enquiries").
		NewLine()
}

func renderIf(w render.Builder, s, f string, cb func(builder render.Builder, s, f string) render.Builder) render.Builder {
	if s != "" {
		return cb(w, s, f)
	}
	return w
}

func white(w render.Builder, s, f string) render.Builder {
	return w.Print(vertical).
		Yellow().
		Printf(f, s).
		End().
		Print(vertical).
		NewLine()
}

func green(w render.Builder, s, f string) render.Builder {
	return w.Print(vertical).
		Green().
		Printf(f, s).
		End().
		Print(vertical).
		NewLine()
}

func yellow(w render.Builder, s, f string) render.Builder {
	return w.Print(vertical).
		Yellow().
		Printf(f, s).
		End().
		Print(vertical).
		NewLine()
}

func red(w render.Builder, s, f string) render.Builder {
	return w.Print(vertical).
		Red().
		Printf(f, s).
		End().
		Print(vertical).
		NewLine()
}
