package darwintty

import (
	"github.com/peter-mount/go-kernel/v2/rest"
	"regexp"
	"strings"
)

type Server struct {
	Server *rest.Server `kernel:"inject"`
	regexp *regexp.Regexp
}

func (s *Server) Start() error {
	s.Server.Handle("/", s.home).Methods("GET")
	s.Server.Handle("/search/{name}", s.search).Methods("GET")
	s.Server.Handle("/{crs}", s.get).Methods("GET")

	re, err := regexp.Compile(`((https?://)([-a-z0-9\.:/]+))`)
	if err == nil {
		s.regexp = re
	}
	return err
}

func (s *Server) respond(r *rest.Rest, src []byte) error {
	switch {
	case IsPlainTextAgent(r.GetHeader("User-Agent")):
		r.ContentType("text/plain")
	default:
		r.ContentType("text/html")

		var a []byte
		a = append(a, strings.Join([]string{
			"<html>",
			"<head>",
			"<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"/>",
			"<meta charset=\"UTF-8\">",
			"<meta name=\"google\" content=\"notranslate\">",
			"<meta http-equiv=\"Content-Language\" content=\"en\">",
			"<link rel=\"stylesheet\" type=\"text/css\" href=\"https://adobe-fonts.github.io/source-code-pro/source-code-pro.css\">",
			"<style type=\"text/css\">",
			"body {background:black;color:#bbbbbb;}",
			"pre, code {",
			"font-family: \"Source Code Pro\", \"DejaVu Sans Mono\", Menlo, \"Lucida Sans Typewriter\", \"Lucida Console\", monaco, \"Bitstream Vera Sans Mono\", monospace;",
			"font-size: 75%;",
			"}",
			"a, a:visited {color:gold;text-decoration:none;}",
			"code {white-space: pre;}",
			"</style>",
			"</head>",
			"<body class=\"\">",
			"<pre>",
		}, "\n")...)

		a = append(a, s.regexp.ReplaceAllFunc(src, expandLinks)...)

		src = append(a, strings.Join([]string{
			"</pre>",
			"</body>",
			"</html>",
		}, "\n")...)
	}

	r.CacheNoCache().
		AccessControlAllowOrigin("*").
		Value(src)
	return nil
}

func expandLinks(f []byte) []byte {
	var a []byte
	a = append(a, "<a href=\""...)
	a = append(a, f...)
	a = append(a, "\">"...)
	a = append(a, f...)
	a = append(a, "</a>"...)
	return a
}

const (
	topLeft     = "┌"
	horiz       = "─"
	topRight    = "┐"
	vertical    = "│"
	midLeft     = "├"
	midCross    = "┼"
	midRight    = "┤"
	bottomLeft  = "└"
	bottomRight = "┘"
)
