package darwintty

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"strings"
)

func (s *Server) search(r *rest.Rest) error {
	name := r.Var("name")

	results, err := s.refClient.Search(name)
	if err != nil {
		return err
	}

	h := "Station name"
	l := len(h)
	for _, result := range results {
		l = Max(l, len(result.Label))
	}
	fmt1 := fmt.Sprintf("%%-%d.%ds ", l, l)

	uriPrefix := *s.Hostname + "/"

	b := render.New().
		Printf("Search for %q returned %d results", name, len(results)).
		NewLine().
		NewLine()

	if len(results) > 0 {
		b = b.Printf(fmt1, "Station").Println("Url").
			Repeat(horiz, l+len(uriPrefix)+3+1).
			NewLine()

		for _, result := range results {
			b = b.Printf(fmt1, result.Label).
				Link(uriPrefix + strings.ToLower(result.Crs)).
				NewLine()
		}

		b = b.Repeat(horiz, l+len(uriPrefix)+3+1).
			NewLine()
	}

	return s.respond(r, b)
}
