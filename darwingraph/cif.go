package darwingraph

import (
	"bufio"
	"github.com/peter-mount/nre-feeds/darwinref"
	"log"
	"os"
	"strings"
)

type cifImporter struct {
	s              *bufio.Scanner // Scanner to read CIF from
	d              *TiplocGraph   // Graph being imported into
	includeRouting bool           // true then include schedules to form an initial map
	prevTiploc     string         // Previous tiploc in sequence
	curTiploc      string         // Current tiploc in sequence
}

func (d *DarwinGraph) importCIF() error {
	log.Printf("Importing CIF %s", *d.cifFileName)

	f, err := os.Open(*d.cifFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	i := &cifImporter{
		d:              d.graph,             // Graph to import into
		s:              bufio.NewScanner(f), // Scanner to read from
		includeRouting: *d.cifRouting,       // Flag to include/exclude CIF routing
	}

	err = i.parse()
	if err != nil {
		return err
	}

	log.Printf("Imported CIF %s", *d.cifFileName)

	return nil
}

func (i *cifImporter) parse() error {
	for i.s.Scan() {
		var err error
		line := i.s.Text()
		if len(line) >= 2 {
			switch line[0:2] {
			case "TA":
				// Tiploc Amend
				err = i.parseTiploc(line)
			case "TI":
				// Tiploc Insert
				err = i.parseTiploc(line)
			case "LO":
				// Location Origin - start sequence
				i.prevTiploc = ""
				i.parseLocation(line)
			case "LI":
				// Location Intermediate - link with prevTiploc
				i.parseLocation(line)
			case "LT":
				// Location Terminate - link with prevTiploc
				i.parseLocation(line)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *cifImporter) parseLocation(l string) {
	if i.includeRouting {
		parseStringTrim(l, 2, 7, &i.curTiploc)
		if i.prevTiploc != "" {
			// Link the two
			edge := i.d.Link(i.prevTiploc, i.curTiploc)
			if edge != nil {
				edge.Src = "CIF"
			}
		}
		i.prevTiploc = i.curTiploc
	}
}

func (i *cifImporter) parseTiploc(l string) error {
	var tpl, name, crs string
	parseStringTrim(l, 2, 7, &tpl)
	parseStringTitle(l, 18, 26, &name)
	parseStringTrim(l, 53, 3, &crs)

	n := i.d.ComputeIfAbsent(tpl, func() *TiplocNode {
		return &TiplocNode{
			Location: darwinref.Location{Tiploc: tpl, Name: name, Crs: crs},
			LocSrc:   "CIF",
		}
	})

	// If the name differs & we are not a station then set it.
	// The station flag is uses so we keep NRE's station names rather than NROD's
	if n.Name != name && !n.Station {
		n.Name = name
		n.LocSrc = "CIF"
	}

	// If CIF has a CRS but not the node then set it.
	// This should not happen but just incase
	if n.Crs == "" {
		n.Crs = crs
		n.Station = n.Location.IsPublic()
		i.d.addCrs(crs, tpl)
	}

	return nil
}

func parseString(line string, s int, l int, v *string) int {
	*v = line[s : s+l]
	return s + l
}

func parseStringTrim(line string, s int, l int, v *string) int {
	var st string
	var ret = parseString(line, s, l, &st)
	*v = strings.Trim(st, " ")
	return ret
}

// Parse a string, trim then title it
func parseStringTitle(line string, s int, l int, v *string) int {
	var st string
	var ret = parseStringTrim(line, s, l, &st)
	*v = strings.Title(strings.ToLower(st))
	return ret
}
