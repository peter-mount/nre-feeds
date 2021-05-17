package update

import (
	"compress/gzip"
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
func (d *ReferenceUpdateService) updateTimetable(tid *darwind3.TimeTableId) error {
	if strings.Contains(tid.TTRefFile, "ref_v3") {
		return nil
	}

	log.Printf("New reference %s", tid.TimeTableId)

	fname := tid.TTRefFile

	err := d.retrieveReference(fname)
	if err != nil {
		return err
	}

	err = d.uploadFile(tid)
	if err != nil {
		return err
	}

	err = d.importReference(tid.TimeTableId, fname)
	if err != nil {
		return err
	}

	return nil
}

func (d *ReferenceUpdateService) retrieveReference(fname string) error {
	file, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return d.config.S3.RetrieveFile(fname, file)
}

func (d *ReferenceUpdateService) importReference(id, fname string) error {
	file, err := os.Open(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	log.Println("Importing", id)
	err = xml.NewDecoder(gr).Decode(d.ref.GetDarwinReference())
	if err != nil {
		return err
	}

	return nil
}

func (d *ReferenceUpdateService) uploadFile(tid *darwind3.TimeTableId) error {
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

		err = d.config.Upload.UploadFile(file, path+tid.TTRefFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *ReferenceUpdateService) findUpdates() {
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
		if strings.Contains(*obj.Key, "ref") {
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
			TTRefFile:   file,
		})
	}
}
