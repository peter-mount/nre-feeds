package darwintty

import (
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"net/http"
	"strings"
)

func (s *Server) get(r *rest.Rest) error {
	crs := strings.ToUpper(r.Var("crs"))

	if len(crs) != 3 {
		r.Status(http.StatusNotFound)
		return s.respond(r, render.New().Println("Not found"))
	}

	result, err := s.ldbClient.GetSchedule(crs)
	if err != nil {
		return err
	}

	if result == nil {
		r.Status(http.StatusNotFound)
		return s.respond(r, render.New().Println("Not found"))
	}

	board := NewBoard(result)

	b := render.New()

	board.Write(b)

	return s.respond(r, b)
}
