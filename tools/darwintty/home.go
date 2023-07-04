package darwintty

import (
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
)

func (s *Server) home(r *rest.Rest) error {

	hostName := *s.Hostname + "/"

	b := render.New().
		Println("UK Rail Departure Boards for the command line").
		NewLine().
		Link(hostName).
		Println(" is a console-oriented service for displaying").
		Println("the departure boards for UK railway stations").
		Println("using terminal-oriented ANSI sequences for").
		Println("console HTTP client (curl, wget etc)").
		Println("or HTML for web browsers.").
		NewLine().
		Println("To use this service:").
		Print("curl ").Link(hostName+"crs").Println(" where crs is the 3 letter CRS code for a station.").
		NewLine().
		Println("Note: with wget try running it with: ").
		White().Print("wget -q -O - ").Link(hostName+"mde").
		NewLine().
		Println("For example:").
		NewLine().
		Print("curl ").Link(hostName+"chx").Println(" For London Charing Cross").
		Print("curl ").Link(hostName+"chc").Println(" For Charing Cross (Glasgow)").
		Print("curl ").Link(hostName+"lbg").Println(" For London Bridge").
		Print("curl ").Link(hostName+"mde").Println(" For Maidstone East").
		NewLine().
		Println("If you do not know the code then use:").
		Print("curl ").Link(hostName+"search/name").Println(" where name is the place name.").
		NewLine().
		Println("For example:").
		NewLine().
		Print("curl ").Link(hostName+"search/maidstone").NewLine().
		Print("curl ").Link(hostName+"search/staplehurst").NewLine().
		Print("curl ").Link(hostName+"search/london").NewLine().
		Print("curl ").Link(hostName+"search/edin").NewLine().
		NewLine().
		Println("All values of crs or search strings are case insensitive.").
		NewLine().
		Println("You can also browse the station index to get the CRS code:").
		Print("curl ").Link(hostName+"index/").NewLine().
		Println("Will return a list sub pages you can use to find a station").
		Printf("as all %d stations are available.", s.getStationCount()).
		NewLine()

	return s.respond(r, b)
}
