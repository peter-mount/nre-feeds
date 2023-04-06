package darwintty

import (
	"bytes"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"strings"
)

type Server struct {
	Server *rest.Server `kernel:"inject"`
}

func (s *Server) Start() error {
	s.Server.Handle("/", s.home).Methods("GET")
	s.Server.Handle("/search/{name}", s.search).Methods("GET")
	s.Server.Handle("/{crs}", s.get).Methods("GET")

	return nil
}

func (s *Server) respond(r *rest.Rest, b render.Builder) error {
	var a []byte

	switch {
	case IsPlainTextAgent(r.GetHeader("User-Agent")):
		r.ContentType("text/plain")
		a = b.BuildAnsi()

	default:
		r.ContentType("text/html")

		var out bytes.Buffer

		out.Write([]byte(strings.Join([]string{
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
			//"font-size: 75%;",
			"}",
			"a, a:visited {color:gold;text-decoration:none;}",
			"code {white-space: pre;}",
			"span { display: inline-block; }",
			".col1 {color:red;}",
			".col2 {color:lightgreen;}",
			".col3 {color:yellow;}",
			".col4 {color:blue;}",
			".col5 {color:magenta;}",
			".col6 {color:cyan;}",
			".col7 {color:white;}",
			"</style>",
			"</head>",
			"<body class=\"\">",
			"<pre>",
		}, "\n")))

		b.BuildHtml(&out)

		out.Write([]byte(strings.Join([]string{
			"</pre>",
			"</body>",
			"</html>",
		}, "\n")))

		a = out.Bytes()
	}

	r.CacheNoCache().
		AccessControlAllowOrigin("*").
		Value(a)
	return nil
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
