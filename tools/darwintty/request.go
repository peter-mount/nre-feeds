package darwintty

import (
	"bytes"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/ldb/client"
	"net/http"
	"strings"
)

func (s *Server) get(r *rest.Rest) error {
	crs := strings.ToUpper(r.Var("crs"))

	if len(crs) != 3 {
		r.Status(http.StatusNotFound)
		return s.respond(r, []byte("Not found"))
	}

	cl := client.DarwinLDBClient{Url: "https://ldb.prod.a51.li"}
	result, err := cl.GetSchedule(crs)
	if err != nil {
		return err
	}

	if result == nil {
		r.Status(http.StatusNotFound)
		return s.respond(r, []byte("Not found"))
	}

	board := NewBoard(result)

	var out bytes.Buffer
	board.Write(&out)

	return s.respond(r, out.Bytes())
}
