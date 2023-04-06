package darwintty

import (
	"bytes"
	"github.com/peter-mount/go-kernel/v2/cron"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
	refClient "github.com/peter-mount/nre-feeds/darwinref/client"
	ldbClient "github.com/peter-mount/nre-feeds/ldb/client"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"strings"
	"sync"
)

type Server struct {
	Server       *rest.Server      `kernel:"inject"`
	Cron         *cron.CronService `kernel:"inject"`
	ldbClient    ldbClient.DarwinLDBClient
	refClient    refClient.DarwinRefClient
	mutex        sync.Mutex
	stations     map[string][]*darwinref.Location
	stationCount int
}

func (s *Server) Start() error {
	s.ldbClient.Url = "https://ldb.prod.a51.li"
	s.refClient.Url = "https://ref.prod.a51.li"

	s.Server.Handle("/", s.home).Methods("GET")
	s.Server.Handle("/search/{name}", s.search).Methods("GET")
	s.Server.Handle("/index/", s.index).Methods("GET")
	s.Server.Handle("/index/{prefix}", s.index).Methods("GET")
	s.Server.Handle("/{crs}", s.get).Methods("GET")

	// Refresh the station index every hour, at 25 mins past the hour,
	// so we don't bombard the remote service all at once
	if _, err := s.Cron.AddFunc("0 25 * * * *", func() {
		_ = s.refreshIndex()
	}); err != nil {
		return err
	}

	// Finish off by refreshing the index now so we have data on startup
	return s.refreshIndex()
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
	topUpper    = "┴"
	vertical    = "│"
	midLeft     = "├"
	midCross    = "┼"
	midRight    = "┤"
	bottomLeft  = "└"
	bottomRight = "┘"
	bottomLower = "┬"
)
