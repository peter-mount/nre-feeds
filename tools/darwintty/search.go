package darwintty

import (
	"bytes"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/darwinref/client"
	"strings"
)

func (s *Server) search(r *rest.Rest) error {
	name := r.Var("name")

	c := client.DarwinRefClient{Url: "https://ref.prod.a51.li"}
	results, err := c.Search(name)
	if err != nil {
		return err
	}

	h := "Station name"
	l := len(h)
	for _, result := range results {
		l = Max(l, len(result.Label))
	}
	fmt1 := fmt.Sprintf("%%-%d.%ds %%s\n", l, l)

	uriPrefix := fmt.Sprintf("https://%s/",
		r.Request().Host)

	var out bytes.Buffer
	fmt.Fprintf(&out, "Search for %q returned %d results\n\n", name, len(results))

	if len(results) > 0 {
		fmt.Fprintf(&out, fmt1, "Station", "Url")
		fmt.Fprintln(&out, strings.Repeat("─", l+len(uriPrefix)+3+1))
		for _, result := range results {
			fmt.Fprintf(&out, fmt1, result.Label, uriPrefix+strings.ToLower(result.Crs))
		}
	}

	r.ContentType("text/plain").
		Value(out.Bytes())

	return nil
}
