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
func (d *TimetableUpdateService) updateTimetable( tid *darwind3.TimeTableId ) error {
  fname := tid.TTFile

  err := d.retrieveTimetable( fname )
  if err != nil {
    return err
  }

  err = d.uploadFile( tid )
  if err != nil {
    return err
  }

  err = d.importTimetable( tid.TimeTableId, fname )
  if err != nil {
    return err
  }

  return nil
}

func (d *TimetableUpdateService) retrieveTimetable( fname string ) error {
  file, err := os.Create( tempFile )
  if err != nil {
    return err
  }
  defer file.Close()

  return d.config.S3.RetrieveFile( fname, file )
}

func (d *TimetableUpdateService) importTimetable( id, fname string ) error {
  file, err := os.Open( tempFile )
  if err != nil {
    return err
  }
  defer file.Close()

  gr, err := gzip.NewReader( file )
  if err != nil {
    return err
  }

  // Run a prune first
  _, err = d.timetable.GetTimetable().PruneSchedules()
  if err != nil {
    return err
  }

  log.Println( "Importing", id )
  err = xml.NewDecoder( gr ).Decode( d.timetable.GetTimetable() )
  if err != nil {
    return err
  }

  // Run a prune afterwards
  _, err = d.timetable.GetTimetable().PruneSchedules()
  if err != nil {
    return err
  }

  return nil
}

func (d *TimetableUpdateService) uploadFile( tid *darwind3.TimeTableId ) error {
  if d.config.Upload.Enabled {
    path, err := tid.GetPath()
    if err != nil {
      return err
    }

    file, err := os.Open( tempFile )
    if err != nil {
      return err
    }
    defer file.Close()

    err = d.config.Upload.UploadFile( file, path + tid.TTFile )
    if err != nil {
      return err
    }
  }
  return nil
}
