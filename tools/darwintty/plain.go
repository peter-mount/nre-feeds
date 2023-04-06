package darwintty

import (
	"bytes"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"io"
	"strings"
)

type TextVisitor struct {
	out         io.Writer
	topLeft     string
	horiz       string
	topRight    string
	vertical    string
	midLeft     string
	midCross    string
	midRight    string
	bottomLeft  string
	bottomRight string

	destLen int
	maxLen  int
	sep     string
	fmt1    string
	fmt2    string
	hdr     string
}

func (t *TextVisitor) serve(r *rest.Rest, b *Board) error {
	var out bytes.Buffer

	t.out = &out

	t.VisitBoard(b)

	r.ContentType("text/plain").
		Value(out.Bytes())

	return nil
}

func (t *TextVisitor) VisitBoard(b *Board) {

	fmt.Fprintf(t.out, "Live train departures for: %s\n", b.Name)

	// Work out max length of destination and via's
	t.destLen = 0
	fullLen := 0
	for _, departure := range b.Departures {
		t.destLen = Max(t.destLen, len(departure.Destination))
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
	if fullLen > (t.destLen + headerLen) {
		t.destLen = Max(t.destLen, fullLen-headerLen)
	}
	t.maxLen = Max(t.destLen+headerLen, fullLen)

	t.fmt1 = fmt.Sprintf("%s %%-%d.%ds %%2.2s %%5.5s %%5.5s %s\n", t.vertical, t.destLen, t.destLen, t.vertical)

	t.hdr = fmt.Sprintf(fmt.Sprintf("%s %%-%d.%ds Pl Deprt Exptd %s", t.vertical, t.destLen, t.destLen, t.vertical), "Destination")

	t.sep = "%s" + t.horiz + strings.Repeat(t.horiz, t.maxLen) + t.horiz + "%s\n"
	t.fmt2 = fmt.Sprintf("%s %%-%d.%ds %s\n", t.vertical, t.maxLen, t.maxLen, t.vertical)

	fmt.Fprintf(t.out, t.sep, t.topLeft, t.topRight)
	fmt.Fprintln(t.out, t.hdr)

	switch len(b.Departures) {
	case 0:
	default:
		for _, departure := range b.Departures {
			fmt.Fprintf(t.out, t.sep, t.midLeft, t.midRight)
			t.VisitDeparture(b, departure)
		}
	}
	fmt.Fprintf(t.out, t.sep, t.bottomLeft, t.bottomRight)
}

func (t *TextVisitor) VisitDeparture(b *Board, departure *Departure) {

	var supp []string

	if departure.Reason != "" {
		supp = append(supp, telstar.Split(departure.Reason, t.maxLen)...)
	}

	supp = append(supp,
		departure.TocName(),
		departure.Coaches(),
		departure.FormationString(),
		departure.ToiletStatus(),
		departure.LastReport.String(),
	)

	fmt.Fprintf(t.out, t.fmt1, departure.Destination, departure.Plat, departure.Depart, departure.Expected)

	for _, s := range supp {
		if s != "" {
			fmt.Fprintf(t.out, t.fmt2, s)
		}
	}
}

func (s Server) servePlain(r *rest.Rest, b *Board) error {
	v := TextVisitor{
		topLeft:     "+",
		horiz:       "-",
		topRight:    "+",
		vertical:    "|",
		midLeft:     "+",
		midCross:    "+",
		midRight:    "+",
		bottomLeft:  "+",
		bottomRight: "+",
	}
	return v.serve(r, b)
}

func (s Server) serveAnsi(r *rest.Rest, b *Board) error {
	v := TextVisitor{
		topLeft:     "┌",
		horiz:       "─",
		topRight:    "┐",
		vertical:    "│",
		midLeft:     "├",
		midCross:    "┼",
		midRight:    "┤",
		bottomLeft:  "└",
		bottomRight: "┘",
	}
	return v.serve(r, b)
}
