package update

import (
  "compress/gzip"
  "encoding/xml"
  "github.com/peter-mount/nre-feeds/darwind3"
  "log"
  "os"
)

const (
  tempFile = "/tmp/inbound.xml.gz"
)

// timetableUpdateListener listens for real time updates for when new reference
// data is made available.
func (d *ReferenceUpdateService) updateTimetable( tid *darwind3.TimeTableId ) error {
  fname := tid.TTRefFile

  err := d.retrieveReference( fname )
  if err != nil {
    return err
  }

  err = d.importReference( tid.TimeTableId, fname )
  if err != nil {
    return err
  }

  return nil
}

func (d *ReferenceUpdateService) retrieveReference( fname string ) error {
  file, err := os.Create( tempFile )
  if err != nil {
    return err
  }
  defer file.Close()

  return d.config.S3.RetrieveFile( fname, file )
}

func (d *ReferenceUpdateService) importReference( id, fname string ) error {
  file, err := os.Open( tempFile )
  if err != nil {
    return err
  }
  defer file.Close()

  gr, err := gzip.NewReader( file )
  if err != nil {
    return err
  }

  log.Println( "Importing", id )
  err = xml.NewDecoder( gr ).Decode( d.ref.GetDarwinReference() )
  if err != nil {
    return err
  }

  return nil
}
