package darwingraph

import (
	"encoding/xml"
	"errors"
	"flag"
	"github.com/peter-mount/go-kernel"
	"log"
	"os"
	"sync"
)

// DarwinGraph is a Kernel service which maintains a TiplocGraph and delegate methods to it
type DarwinGraph struct {
	mutex            sync.Mutex // Mutex to protect graph
	graph            *RailGraph // The graph representing the UK rail network
	importFileName   *string    // -import filename to create an initial model
	stationsFileName *string    // -kbstation filename to import data from the NRE Knowledge Base
	cifFileName      *string    // -cif filename to import from an NR CIF file
	cifRouting       *bool      // -no-cif-routing to ignore routing in -cif
	xmlFileName      *string    // -xml filename to load/save the model
	saveModel        *bool      // -save indicates we want to save the model
	tiplocFileName   *string    // -tiplocExport to import from Legolash2o tiploc location map
	nreCacheName     *string    // -nre-schedule-cache copy of schedule cache to import from
}

func (d *DarwinGraph) Name() string {
	return "DarwinGraph"
}

func (d *DarwinGraph) Init(_ *kernel.Kernel) error {
	d.importFileName = flag.String("import", "", "Import tiploc data")
	d.xmlFileName = flag.String("xml", "", "xml filename for the graph")
	d.saveModel = flag.Bool("save", false, "save the model if -xml is set")
	d.stationsFileName = flag.String("kbstation", "", "xml to import KB data into the graph")
	d.cifFileName = flag.String("cif", "", "Network Rail CIF file to import data into the graph")
	d.cifRouting = flag.Bool("cif-routing", false, "With -cif, true to import routing from CIF as well as locations")
	d.tiplocFileName = flag.String("tiploc-location", "", "Import tiploc locations from legolash2o export")
	d.nreCacheName = flag.String("nre-schedule-cache", "", "nre-feeds schedule cache to import from")
	return nil
}

func (d *DarwinGraph) PostInit() error {
	if *d.xmlFileName == "" {
		return errors.New("-xml is required")
	}
	return nil
}

func (d *DarwinGraph) Start() error {
	d.graph = NewRailGraph()

	populate := false

	// Import the model on start
	if *d.importFileName != "" {
		err := d.importFile()
		if err != nil {
			return err
		}
		populate = true
	} else {
		// If not importing the model, load the graph
		err := d.LoadGraph()
		if err != nil {
			return err
		}
	}

	if *d.stationsFileName != "" {
		// Import information from the NRE KB feed
		err := d.importKBStations()
		if err != nil {
			return err
		}
		populate = true
	}

	if *d.cifFileName != "" {
		// Import information from the NRE KB feed
		err := d.importCIF()
		if err != nil {
			return err
		}
		populate = true
	}

	if *d.tiplocFileName != "" {
		err := d.importTiplocLocations()
		if err != nil {
			return err
		}
		populate = true
	}

	if *d.nreCacheName != "" {
		err := d.importNreFeedsCache()
		if err != nil {
			return err
		}
	}

	if populate {
		log.Println("Repopulating stations")
		d.graph.populateStations()
		log.Println("Repopulated stations")

		log.Println("Generating line segments")
		d.graph.generateLineSegments()
		log.Println("Generated line segments")
	}

	// Once started save the current graph (if enabled)
	return d.SaveGraph()
}

func (d *DarwinGraph) LoadGraph() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if *d.xmlFileName == "" {
		return nil
	}

	log.Printf("Restoring graph from %s", *d.xmlFileName)
	f, err := os.Open(*d.xmlFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	e := xml.NewDecoder(f)
	err = e.Decode(d.graph)
	if err != nil {
		return err
	}

	log.Printf("Loaded graph from %s\n", *d.xmlFileName)

	return nil
}

func (d *DarwinGraph) SaveGraph() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !*d.saveModel || *d.xmlFileName == "" {
		return nil
	}

	log.Printf("Persisting graph to %s", *d.xmlFileName)

	f, err := os.Create(*d.xmlFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	e := xml.NewEncoder(f)
	e.Indent("", "  ")
	err = e.Encode(d.graph)
	if err != nil {
		return err
	}

	log.Printf("Persisted graph to %s", *d.xmlFileName)
	return nil
}

// GetTiplocNode returns an existing TiplocNode or nil if it doesn't exist
func (d *DarwinGraph) GetTiplocNode(tiploc string) *TiplocNode {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.graph.GetTiploc(tiploc)
}

// ComputeIfAbsent returns an existing TiplocNode or creates one using the supplied function
func (d *DarwinGraph) ComputeIfAbsent(tiploc string, f func() RailNode) RailNode {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.graph.ComputeIfAbsent(tiploc, f)
}

// LinkTiplocs links two tiplocs together
// Returns the new TiplocEdge or nil if one already exists
func (d *DarwinGraph) LinkTiplocs(a, b string) *TiplocEdge {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.graph.LinkTiplocs(a, b)
}

func (d *DarwinGraph) ForEachNode(f func(node RailNode)) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	nodes := d.graph.graph.Nodes()
	for nodes.Next() {
		f(nodes.Node().(RailNode))
	}
}

func (d *DarwinGraph) ForEachTiplocNode(f func(*TiplocNode)) {
	d.ForEachNode(func(n RailNode) {
		if n != nil && n.NodeType() == NodeTiploc {
			f(n.(*TiplocNode))
		}
	})
}

func (d *DarwinGraph) ForEachStationNode(f func(*StationNode)) {
	d.ForEachNode(func(n RailNode) {
		if n != nil && n.NodeType() == NodeStation {
			f(n.(*StationNode))
		}
	})
}

func (d *DarwinGraph) ForEachTiplocEdge(f func(*TiplocEdge)) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	edges := d.graph.graph.Edges()
	for edges.Next() {
		f(edges.Edge().(*TiplocEdge))
	}
}

func (d *DarwinGraph) ForEachStationEdge(f func(edge *StationEdge)) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.graph.stations != nil {
		edges := d.graph.stations.Edges()
		for edges.Next() {
			f(edges.Edge().(*StationEdge))
		}
	}
}
