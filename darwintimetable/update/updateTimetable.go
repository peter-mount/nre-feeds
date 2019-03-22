package update

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/peter-mount/nre-feeds/darwind3"
  "os"
  "log"
)

const (
  tempFile = "/tmp/inbound.xml.gz"
)

// timetableUpdateListener listens for real time updates for when new reference
// data is made available.
func (d *TimetableUpdateService) updateTimetable( t *darwind3.TimeTableId ) error {
  err := d.retrieveFile( t.TTFile )
  if err != nil {
    return err
  }

  return nil
}

func (d *TimetableUpdateService) retrieveFile( fname string ) error {
  file, err := os.Create( fname )
   if err != nil {
     return err
   }

  defer file.Close()
  log.Printf("Retrieving %s", fname )
  log.Println( d.config.S3.Region )
  log.Println( d.config.S3.Bucket )

  sess, _ := session.NewSession(
    &aws.Config{
      Region: aws.String(d.config.S3.Region),
      Credentials: credentials.NewStaticCredentials(
        d.config.S3.AccessKey,
        d.config.S3.SecretKey,
        "",
      ),
    })

  downloader := s3manager.NewDownloader(sess)

  numBytes, err := downloader.Download(
    file,
    &s3.GetObjectInput{
      Bucket: aws.String(d.config.S3.Bucket),
      Key:    aws.String( d.config.S3.Path + "/" + fname ),
    })
  if err != nil {
    return err
  }

  log.Println("Downloaded", file.Name(), numBytes, "bytes")

  return nil
}
