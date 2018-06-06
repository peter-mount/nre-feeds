package client

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "net/http"
)

// A remove client to the DarwinTimetable microservice
type DarwinRefClient struct {
  // The url prefix, e.g. "http://localhost:8080" of the remote service
  // Note no trailing "/" as the client will add a patch starting with "/"
  Url string
}

// Make a GET to a remote service
// path - path of rest endpoint
// v - instance of object to unmarshal
// returns (true, nil) if found and v contains data
// (false, nil) if not found or (false, error ) on error
func (c *DarwinRefClient) get( path string, v interface{} ) ( bool, error ) {
  if resp, err := http.Get( c.Url + path ); err != nil {
    return false, err
  } else {
    defer resp.Body.Close()

    if resp.StatusCode == 404 {
      return false, nil
    }

    if body, err := ioutil.ReadAll( resp.Body ); err != nil {
      return false, err
    } else {
      json.Unmarshal( body, v )
      return true, nil
    }
  }
}

func (c *DarwinRefClient) post( path string, b interface{}, v interface{} ) ( bool, error ) {
  if post, err := json.Marshal( b ); err != nil {
    return false, err
  } else if resp, err := http.Post( c.Url + path, "application/json", bytes.NewReader( post ) ); err != nil {
    return false, err
  } else {
    defer resp.Body.Close()

    if resp.StatusCode == 404 {
      return false, nil
    }

    if body, err := ioutil.ReadAll( resp.Body ); err != nil {
      return false, err
    } else {
      json.Unmarshal( body, v )
      return true, nil
    }
  }
}
