package darwintty

import (
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/ldb/client"
	"net/http"
	"strings"
)

func (s *Server) get(r *rest.Rest) error {

	log.Printf("  User-Agent %q", r.GetHeader("User-Agent"))
	log.Printf("     Accepts %q", r.GetHeader("Accepts"))
	log.Printf("Content-Type %q", r.GetHeader("Content-Type"))
	if IsPlainTextAgent(r.GetHeader("User-Agent")) {
		log.Println("Plain user")
	}

	crs := strings.ToUpper(r.Var("crs"))

	cl := client.DarwinLDBClient{Url: "https://ldb.prod.a51.li"}
	result, err := cl.GetSchedule(crs)
	if err != nil {
		return err
	}

	board := NewBoard(result)

	switch {
	case IsPlainTextAgent(r.GetHeader("User-Agent")):
		return s.serveAnsi(r, board)

	default:
		r.Status(http.StatusNotFound)
	}

	return nil
}
