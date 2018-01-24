package darwinupdate

import (
  "compress/gzip"
  "encoding/xml"
  "github.com/jlaffaye/ftp"
  "github.com/peter-mount/golib/rest"
  "log"
  "regexp"
  "sort"
  "strings"
)

func (u *DarwinUpdate) ReferenceHandler( r *rest.Rest ) error {
  if err := u.ftp( u.ReferenceUpdate ); err != nil {
    return err
  }

  r.Status( 200 ).Value( "ok" )
  return nil
}

func (u *DarwinUpdate) ReferenceUpdate( con *ftp.ServerConn ) error {
  log.Println( "Looking for reference files" )

  entries, err := con.List( "." )
  if err != nil {
    return err
  }

  re := regexp.MustCompile( ".*_ref_v3.xml.gz" )

  var files []*ftp.Entry

  for _, e := range entries {
    if re.MatchString( e.Name ) {
      files = append( files, e )
    }
  }

  // Sort as in ISO format
  sort.SliceStable( files, func( i, j int ) bool {
    return strings.Compare( files[i].Name, files[i].Name ) < 0
  })

  if len( files ) < 1 {
    log.Println( "No reference files found" )
    return nil
  }

  file := files[ len(files)-1 ]
  tid := u.Ref.TimetableId()
  if tid != "" && strings.Compare( file.Name[:len(tid)], tid ) <= 0 {
    log.Println( "Ignoring", file.Name, "timetableId", tid )
    return nil
  }

  log.Println( "Retrieving", file.Name, file.Size, file.Time )

  resp, err := con.Retr( file.Name )
  if err != nil {
    return err
  }

  gr, err := gzip.NewReader( resp )
  if err != nil {
    log.Println( "Failed to gunzip")
    resp.Close()
    return err
  }

  if err := xml.NewDecoder( gr ).Decode( u.Ref ); err != nil {
    log.Println( err )
    resp.Close()
    return err
  }

  return resp.Close()
}