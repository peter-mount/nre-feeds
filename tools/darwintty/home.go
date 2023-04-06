package darwintty

import (
	"bytes"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
)

func (s *Server) home(r *rest.Rest) error {

	hostName := r.Request().Host

	var out bytes.Buffer

	fmt.Fprintf(&out,
		"To use this service:\n"+
			"http://%s/crs where crs is the 3 letter CRS code for a station.\n"+
			"\nFor example:\n"+
			"http://%s/chx For London Charing Cross\n"+
			"http://%s/chc For Charing Cross (Glasgow)\n"+
			"http://%s/lbg For London Bridge\n"+
			"http://%s/mde For Maidstone East\n"+
			"\n\nIf you do not know the code then use\n"+
			"http://%s/search/name where name is the place name.\n"+
			"\nFor example:\n\n"+
			"http://%s/search/maidstone\n"+
			"http://%s/search/staplehurst\n"+
			"http://%s/search/london\n"+
			"http://%s/search/edin\n"+
			"\nAll values of crs or search strings are case insensitive.",
		hostName, hostName, hostName, hostName, hostName, hostName,
		hostName, hostName, hostName, hostName,
	)

	return s.respond(r, out.Bytes())
}
