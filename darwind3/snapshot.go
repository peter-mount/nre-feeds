package darwind3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"github.com/jlaffaye/ftp"
	"log"
	"sync"
	"time"
)

func (fs *FeedStatus) loadSnapshot(ts time.Time, m *sync.Mutex) error {
	m.Lock()
	defer m.Unlock()

	return fs.d3.ftpClient(func(con *ftp.ServerConn) error {
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
		log.Println("Retrieving", n)
		r, err := con.Retr(n)
		if err != nil {
			return err
		}
		defer r.Close()

		log.Println("Importing snapshot")

		gr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}

		lc := 0
		scanner := bufio.NewScanner(gr)
		for scanner.Scan() {
			lc++
			if (lc % 1000) == 0 {
				log.Println("Imported", lc)
			}
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
		}

		log.Println("Finished importing", lc, "messages")

		return nil
	})
}
