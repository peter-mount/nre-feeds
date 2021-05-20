package darwingraph

import (
	"encoding/json"
	"github.com/peter-mount/nre-feeds/darwind3"
	"log"
	"os"
	"path/filepath"
)

func (d *DarwinGraph) importNreFeedsCache() error {
	log.Printf("Importing NRE-Feeds schedule cache from %s", *d.nreCacheName)
	imp := &nreFeedImport{d: d}
	defer imp.status()

	return filepath.Walk(*d.nreCacheName, imp.process)
}

type nreFeedImport struct {
	d         *DarwinGraph
	fileCount int // Number of files read
}

func (d *nreFeedImport) status() {
	log.Printf("Imported %d schedules", d.fileCount)
}

func (d *nreFeedImport) process(path string, info os.FileInfo, err error) error {
	if err == nil && !info.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Failed to open %s, %s", path, err.Error())
			return err
		}
		defer f.Close()

		sched := darwind3.Schedule{}
		err = json.NewDecoder(f).Decode(&sched)
		if err != nil {
			log.Printf("Failed to parse %s, %s", path, err.Error())
			return err
		}

		// Ignore Bus & Ship services
		if !(sched.TrainId == "0B00" || sched.TrainId == "0S00") {
			err = d.importNreFeedsSchedule(sched)
			if err != nil {
				log.Printf("Failed to import %s, %s", path, err.Error())
				return err
			}
		}

		d.fileCount++
		if (d.fileCount % 10000) == 0 {
			d.status()
		}
	}
	return err
}

func (d *nreFeedImport) importNreFeedsSchedule(s darwind3.Schedule) error {
	var prevLoc *darwind3.Location
	for _, loc := range s.Locations {
		if prevLoc != nil {
			d.d.LinkTiplocs(prevLoc.Tiploc, loc.Tiploc)
		}
		prevLoc = loc
	}
	return nil
}
