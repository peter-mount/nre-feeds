package darwind3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	bolt "github.com/etcd-io/bbolt"
	"github.com/jlaffaye/ftp"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"os"
	"sort"
	"time"
)

type logEntry struct {
	path string
	time time.Time
	meta string
}

func (fs *FeedStatus) loadSnapshot(ts time.Time) error {
	log.Println("Suspending realtime message processing for NRDP synchronization")

	// List of retrieved files & delete them once we are done
	defer fs.cleanup()

	// declare err here & don't use := inside the ftpClient call else the new entries slice won't be exposed to us!
	var err error

	err = fs.d3.ftpClient(func(con *ftp.ServerConn) error {
		// latest full snapshot first
		fs.entries, err = fs.resolveFiles("snapshot", con, fs.entries)
		if err != nil {
			return err
		}

		// The pushport log files next
		fs.entries, err = fs.resolveFiles("pushport", con, fs.entries)
		if err != nil {
			return err
		}

		return nil
	})

	log.Println("Found", len(fs.entries), "files for import")

	for _, entry := range fs.entries {
		err = fs.importLogEntry(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FeedStatus) cleanup() {
	// Remove any entries from local disk
	for _, entry := range fs.entries {
		log.Println("Deleting", entry.path)
		err := os.Remove(entry.path)
		if err != nil {
			log.Println(err)
		}
	}

	// Empty the slice
	fs.entries = nil

	log.Println("Resuming realtime message processing")
}

func (fs *FeedStatus) resolveFiles(dirname string, con *ftp.ServerConn, origFiles []logEntry) ([]logEntry, error) {
	// The latest time we imported a file for this directory
	var latestTime time.Time
	err := fs.d3.GetMeta(dirname, &latestTime)
	if err != nil {
		return origFiles, err
	}

	if latestTime.IsZero() {
		log.Println("Looking for", dirname)
	} else {
		log.Println("Looking for", dirname, "after", latestTime.Format(util.HumanDateTime))
	}

	entries, err := con.List(dirname)
	if err != nil {
		return origFiles, err
	}

	var files []logEntry

	for _, entry := range entries {
		if entry.Time.After(latestTime) {
			srcName := dirname + "/" + entry.Name

			destName := "/tmp/" + entry.Name
			err = FtpCp(con, srcName, destName)
			if err != nil {
				return sortLogEntry(origFiles, files), err
			}

			files = append(files, logEntry{destName, entry.Time, dirname})
		}
	}

	// Sort this list of files by their timestamps
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].time.Before(files[j].time)
	})

	// return
	return sortLogEntry(origFiles, files), nil
}

// Sorts the new slice by their timestamps then appends them to orig
func sortLogEntry(orig, new []logEntry) []logEntry {
	sort.SliceStable(new, func(i, j int) bool {
		return new[i].time.Before(new[j].time)
	})

	return append(orig, new...)
}

func (fs *FeedStatus) importLogEntry(entry logEntry) error {
	err := fs.d3.BulkUpdate(func(tx *bolt.Tx) error {
		log.Println("Importing", entry.path)

		f, err := os.Open(entry.path)
		if err != nil {
			return err
		}
		defer f.Close()

		gr, err := gzip.NewReader(f)
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

			if len(ln) > 0 {
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
		}

		log.Println("Finished importing", lc, "messages")
		return fs.d3.PutMetaTx(tx, entry.meta, entry.time)
	})
	return err
}
