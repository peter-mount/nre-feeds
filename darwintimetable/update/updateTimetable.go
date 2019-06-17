package update

import (
	"compress/gzip"
	"encoding/xml"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/golib/kernel/logger"
	"github.com/peter-mount/nre-feeds/darwind3"
	"os"
)

const (
	tempFile = "/tmp/inbound.xml.gz"
)

// timetableUpdateListener listens for real time updates for when new reference
// data is made available.
func (d *TimetableUpdateService) updateTimetable(tid *darwind3.TimeTableId) error {
	return d.logger.Report(
		"NRDP Timetable Update",
		func(log *logger.Logger) error {

			log.Printf("New timetable %s", tid.TimeTableId)

			fname := tid.TTFile

			err := d.retrieveTimetable(log, fname)
			if err != nil {
				return err
			}

			err = d.uploadFile(log, tid)
			if err != nil {
				return err
			}

			err = d.importTimetable(log, tid.TimeTableId, fname)
			if err != nil {
				return err
			}

			return nil
		})
}

func (d *TimetableUpdateService) retrieveTimetable(log *logger.Logger, fname string) error {
	file, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return d.config.S3.RetrieveFile(log, fname, file)
}

func (d *TimetableUpdateService) importTimetable(log *logger.Logger, id, fname string) error {
	file, err := os.Open(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	// Run a prune first
	_, err = d.timetable.GetTimetable().PruneSchedules()
	if err != nil {
		return err
	}

	log.Println("Importing", id)
	err = xml.NewDecoder(gr).Decode(d.timetable.GetTimetable())
	if err != nil {
		return err
	}

	// Run a prune afterwards
	_, err = d.timetable.GetTimetable().PruneSchedules()
	if err != nil {
		return err
	}

	// Report current database size
	return d.timetable.GetTimetable().View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("DarwinAssoc"))
		if b != nil {
			log.Printf("Associations: %d", b.Stats().KeyN)
		}

		b = tx.Bucket([]byte("DarwinJourney"))
		if b != nil {
			log.Printf("    Journeys: %d", b.Stats().KeyN)
		}

		return nil
	})
}

func (d *TimetableUpdateService) uploadFile(log *logger.Logger, tid *darwind3.TimeTableId) error {
	if d.config.Upload.Enabled {
		path, err := tid.GetPath()
		if err != nil {
			return err
		}

		file, err := os.Open(tempFile)
		if err != nil {
			return err
		}
		defer file.Close()

		err = d.config.Upload.UploadFile(log, file, path+tid.TTFile)
		if err != nil {
			return err
		}
	}
	return nil
}
