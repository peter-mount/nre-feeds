package darwinkb

import (
  "bytes"
  "encoding/json"
  "github.com/peter-mount/golib/kernel/bolt"
  "errors"
)

func (r *DarwinKB) View( n string, f func( *bolt.Bucket ) error ) error {
  db, exists := r.db[n]
  if !exists {
    return errors.New( "No bucket " + n )
  }
  return db.View( func( tx *bolt.Tx ) error {
    bucket := tx.Bucket( n )
    if bucket == nil {
      return errors.New( "Bucket " + n + " not found" )
    }
    return f( bucket )
  } )
}

func (r *DarwinKB) Update( n string, f func( *bolt.Bucket ) error ) error {
  db, exists := r.db[n]
  if !exists {
    return errors.New( "No bucket " + n )
  }
  return db.Update( func( tx *bolt.Tx ) error {
    bucket := tx.Bucket( n )
    if bucket == nil {
      return errors.New( "Bucket " + n + " not found" )
    }
    return f( bucket )
  } )
}

// Tests to see if a bucket is empty
func (r *DarwinKB) bucketEmpty( name string ) (bool, error) {
  empty := false
  err := r.View( name, func( bucket *bolt.Bucket ) error {
    cursor := bucket.Cursor()
    k, _ := cursor.First()
    if k == "" {
      empty = true
    }
    return nil
  })
  if err != nil {
    return false, err
  }
  return empty, nil
}

func bucketRemoveAll( bucket *bolt.Bucket ) error {
  return bucket.ForEach( func(k string, v []byte) error {
    return bucket.Delete( k )
  })
}

// unmarshalBytes unmarshals json returning an object
func unmarshalBytes( b *bytes.Buffer ) ( map[string]interface{}, error ) {
  a := make( map[string]interface{} )
  err := json.Unmarshal( b.Bytes(), &a )
  if err != nil {
    return nil, err
  }
  return a, nil
}

func GetJsonObjectValue( r map[string]interface{}, n ...string ) ( interface{}, bool ) {
  var o map[string]interface{}
  var v interface{}
  for _, k := range n {
    if o == nil {
      o = r
    } else {
      o = v.(map[string]interface{})
    }

    var e bool
    v, e = o[ k ]

    if !e || v == nil {
      return nil, false
    }
  }
  return v, true
}

func GetJsonObject( o map[string]interface{}, n ...string ) ( map[string]interface{}, bool ) {
  v, e := GetJsonObjectValue( o, n... )
  if e {
    return v.(map[string]interface{}), true
  }
  return nil, false
}

func GetJsonArray( o map[string]interface{}, n ...string ) ( []interface{}, bool ) {
  v, e := GetJsonObjectValue( o, n... )
  if e {
    if a, ok := v.([]interface{}); ok {
      return a, ok
    }
    var a []interface{}
    a = append( a, v )
    return a, true
  }
  return nil, false
}
