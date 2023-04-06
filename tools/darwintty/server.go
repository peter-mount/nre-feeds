package darwintty

import (
	"github.com/peter-mount/go-kernel/v2/rest"
)

type Server struct {
	Server *rest.Server `kernel:"inject"`
}

func (s *Server) Start() error {
	s.Server.Handle("/search/{name}", s.search).Methods("GET")
	s.Server.Handle("/{crs}", s.get).Methods("GET")
	return nil
}
