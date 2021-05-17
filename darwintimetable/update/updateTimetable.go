package update

import (
	"compress/gzip"
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"log"
	"os"
	"strings"
)

const (
	tempFile = "/tmp/inbound.xml.gz"
)

// timetableUpdateListener listens for real time updates for when new reference
// data is made available.
func (d *TimetableUpdateService) updateTimetable(tid *darwind3.TimeTableId) error {
	if !strings.Contains(tid.TTRefFile, "v8") {
		return nil
	}

	log.Printf("New timetable %s", tid.TimeTableId)

	fname := tid.TTFile

	err := d.retrieveTimetable(fname)
	if err != nil {
		return err
	}

	err = d.uploadFile(tid)
	if err != nil {
		return err
	}

	err = d.importTimetable(tid.TimeTableId, fname)
	if err != nil {
		return err
	}

	return nil
}

func (d *TimetableUpdateService) retrieveTimetable(fname string) error {
	file, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return d.config.S3.RetrieveFile(fname, file)
}

func (d *TimetableUpdateService) importTimetable(id, fname string) error {
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

func (d *TimetableUpdateService) uploadFile(tid *darwind3.TimeTableId) error {
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

		err = d.config.Upload.UploadFile(file, path+tid.TTFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *TimetableUpdateService) findUpdates() {
	log.Println("Looking for updates")

	config := d.config.S3

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(config.Region),
			Credentials: credentials.NewStaticCredentials(
				config.AccessKey,
				config.SecretKey,
				"",
			),
		})
	if err != nil {
		log.Println("Failed to create aws session", err)
	}

	client := s3.New(sess)

	input := &s3.ListObjectsInput{
		Bucket:  aws.String(config.Bucket),
		MaxKeys: aws.Int64(256),
	}

	objs, err := client.ListObjects(input)
	if err != nil {
		log.Println("Failed to find updates", err)
	}

	var file string
	for _, obj := range objs.Contents {
		if strings.Contains(*obj.Key, "v8.xml") && !strings.Contains(*obj.Key, "ref") {
			file = *obj.Key
		}
	}

	if file != "" {
		if strings.HasPrefix(file, config.Path) {
			file = file[len(config.Path):]
		}

		id := file
		if i := strings.Index(id, "_"); i > -1 {
			id = id[:i]
		}

		_ = d.updateTimetable(&darwind3.TimeTableId{
			TimeTableId: id,
			TTFile:      file,
		})
	}
}
