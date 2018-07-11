package darwinkb

import (
  "encoding/json"
  "errors"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/exec"
  "strings"
  "time"
)

type KBToken struct {
  Username      string              `json:"username"`
  Roles         map[string]string   `json:"roles"`
  Token         string              `json:"token"`
  tokenDate     time.Time
}

// Set the authentication token, requesting a new one as required
func (k *DarwinKB) setToken( mainReq *http.Request ) error {

  // Tokens last for an hour so if it's older than 45 minutes we'll request a new one
  now := time.Now().UTC()
  refresh := k.token.tokenDate.Add( 45 * time.Minute )
  if now.After( refresh ) {
    authString := "username=" + k.config.KB.Username + "&password=" + k.config.KB.Password
    payload := strings.NewReader( authString )

    req, err := http.NewRequest( "POST", "https://datafeeds.nationalrail.co.uk/authenticate", payload )
    if err != nil {
      log.Println( "DarwinKB: ", err )
      return err
    }

    req.Header.Add( "Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add( "Accept", "application/json, text/plain, */*")

    log.Println("DarwinKB: Requesting new token" )
    resp, err := http.DefaultClient.Do( req )
    if err != nil {
      log.Println( "DarwinKB: ", err )
      return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll( resp.Body )
    if err != nil {
      log.Println( "DarwinKB: ", err )
      return err
    }

    json.Unmarshal( body, &k.token )
    log.Println( "DarwinKB: Token issued" )

    k.token.tokenDate = now
  } else {
    log.Println( "DarwinKB: Reusing existing token" )
  }

  // Set the token on the outer request
  mainReq.Header.Add( "X-Auth-Token", k.token.Token )
  return nil
}

func (k *DarwinKB) refreshFile( filename, url string, maxAge time.Duration ) (bool,error) {
  fname := k.config.KB.DataDir + "static/" + filename

  finfo, err := os.Stat( fname )
  if err != nil {
    if os.IsNotExist( err ) {
      log.Println( "DarwinKB: File not exist", fname )
      return true, k.retrieveFile( fname, url )
    }
    log.Println( "DarwinKB:", err )
    return false, err
  }

  now := time.Now()
  if now.Sub( finfo.ModTime() ) >= maxAge {
    return true, k.retrieveFile( fname, url )
  }

  log.Println( "DarwinKB: Keeping file", fname )
  return false, nil
}

func (k *DarwinKB) retrieveFile( fname, url string ) error {

  // temp file name
  tempname := fname + ".temp"
  err := k.retrieveFileImpl( tempname, url )
  if err != nil {
    return err
  }

  // Although the xml files say they are utf-8 they really are iso-8859-1 so
  // run though iconv
  log.Println( "DarwinKB: Fixing utf-8 encoding on", fname )
  cmd := exec.Command( "iconv", "-f", "iso-8859-1", "-t", "utf-8", tempname )

  file, err := os.OpenFile( fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644 )
  if err != nil {
    log.Println( "DarwinKB:", err )
    return err
  }
  defer file.Close()
  cmd.Stdout = file

  err = cmd.Run()
  if err != nil {
    log.Println( "DarwinKB: iconv failed,", err )
    return err
  }

  log.Println( "DarwinKB: Retrieved", fname )
  return nil
}

func (k *DarwinKB) retrieveFileImpl( fname, url string ) error {

  req, err := http.NewRequest( "GET", url, nil )
  if err != nil {
    log.Println( "DarwinKB:", err )
    return err
  }

  err = k.setToken( req )
  if err != nil {
    return err
  }

  log.Println("DarwinKB: Retrieving ", url )
  resp, err := http.DefaultClient.Do( req )
  if err != nil {
    log.Println( "DarwinKB:", err )
    return err
  }
  defer resp.Body.Close()

  log.Println("DarwinKB: Response ", resp.Status, " length ", resp.ContentLength, "Uncompressed ", resp.Uncompressed )
  if resp.StatusCode > 399 {
    return errors.New( resp.Status )
  }

  file, err := os.OpenFile( fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644 )
  if err != nil {
    log.Println( "DarwinKB:", err )
    return err
  }
  defer file.Close()

  w, err := io.Copy( file, resp.Body )
  if err != nil {
    log.Println( "DarwinKB:", err )
    return err
  }

  log.Println( "DarwinKB: Written", w, "to", fname )
  return nil
}
