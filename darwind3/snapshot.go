package darwind3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	bolt "github.com/etcd-io/bbolt"
	"github.com/jlaffaye/ftp"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"sync"
	"time"
)

func (fs *FeedStatus) loadSnapshot(ts time.Time, m *sync.Mutex) error {
	m.Lock()
	defer m.Unlock()

	return fs.d3.ftpClient(func(con *ftp.ServerConn) error {
		err := fs.processSnapshot(con)
		if err != nil {
			return err
		}
		return nil
	})
}

func (fs *FeedStatus) processSnapshot(con *ftp.ServerConn) error {
	log.Println("Looking for latest snapshot")

	entries, err := con.List("snapshot")
	if err != nil {
		return err
	}

	var entry *ftp.Entry
	for _, e := range entries {
		if e.Name == "snapshot.gz" {
			entry = e
		}
	}

	if entry == nil {
		return errors.New("Not found a snapshot")
	}

	n := "snapshot/" + entry.Name

	err = fs.d3.GetMeta("snapshot", &fs.snapshotTime)
	if err != nil {
		return err
	}
	log.Println("Dates", fs.snapshotTime.Format(util.HumanDateTime), entry.Time.Format(util.HumanDateTime))

	if !entry.Time.After(fs.snapshotTime) {
		log.Println("Not retrieving", n, "as not newer than last one")
		return nil
	}

	log.Println("Retrieving", n)
	r, err := con.Retr(n)
	if err != nil {
		return err
	}
	defer r.Close()

	// Disable remote timetable resolution for the duration then run on one single tx
	oldTT := fs.d3.Timetable
	fs.d3.Timetable = ""
	err = fs.d3.BulkUpdate(func(tx *bolt.Tx) error {
		log.Println("Importing snapshot")

		gr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}

		lc := 0
		scanner := bufio.NewScanner(gr)
		for scanner.Scan() {
			ln := scanner.Bytes()
			err = scanner.Err()
			if err != nil {
				return err
			}

			p := &Pport{}
			r := bytes.NewReader(ln)
			err := xml.NewDecoder(r).Decode(p)
			if err != nil {
				return err
			}

			err = p.Process(fs.d3)
			if err != nil {
				return err
			}

			lc++
			if (lc % 1000) == 0 {
				log.Println("Imported", lc)
			}
		}

		log.Println("Finished importing", lc, "messages")
		return fs.d3.PutMetaTx(tx, "snapshot", entry.Time)
	})
	fs.d3.Timetable = oldTT
	if err != nil {
		return err
	}

	fs.snapshotTime = entry.Time

	return nil
}
