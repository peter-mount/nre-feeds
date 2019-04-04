package ldb

import (
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwinref"
	"github.com/peter-mount/nre-feeds/darwinref/client"
	"log"
)

// RefreshStations ensures we have all Public stations defined on startup.
// Not doing so incurs a performance hit when a train references it for the
// first time.
func (d *LDB) RefreshStations() {
	err := d.Update(func(tx *bbolt.Tx) error {

		log.Println("LDB: Initialising stations")

		refClient := &client.DarwinRefClient{Url: d.Reference}

		locations, err := refClient.GetStations()
		if err != nil {
			return err
		}

		if locations != nil {
			// Map tiplocs by crs
			m := make(map[string][]*darwinref.Location)
			for _, loc := range locations {
				if s, ok := m[loc.Crs]; ok {
					m[loc.Crs] = append(s, loc)
				} else {
					s = make([]*darwinref.Location, 0)
					m[loc.Crs] = append(s, loc)
				}
			}

			// create a station for each
			for _, l := range m {
				createStation(tx, l)
			}
		}

		log.Printf("LDB: %d stations initialized\n", tx.Bucket([]byte(crsBucket)).Stats().KeyN)
		return nil
	})
	if err != nil {
		log.Println("LDB: Station import failed", err)
	}
}
