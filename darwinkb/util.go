package darwinkb

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/peter-mount/go-kernel/bolt"
)

func (r *DarwinKB) View(n string, f func(*bolt.Bucket) error) error {
	db, exists := r.db[n]
	if !exists {
		return errors.New("No bucket " + n)
	}
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(n)
		if bucket == nil {
			return errors.New("Bucket " + n + " not found")
		}
		return f(bucket)
	})
}

func (r *DarwinKB) Update(n string, f func(*bolt.Bucket) error) error {
	db, exists := r.db[n]
	if !exists {
		return errors.New("No bucket " + n)
	}
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(n)
		if bucket == nil {
			return errors.New("Bucket " + n + " not found")
		}
		return f(bucket)
	})
}

// Tests to see if a bucket is empty
func (r *DarwinKB) bucketEmpty(name string) (bool, error) {
	empty := false
	err := r.View(name, func(bucket *bolt.Bucket) error {
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

func bucketRemoveAll(bucket *bolt.Bucket) error {
	return bucket.ForEach(func(k string, v []byte) error {
		return bucket.Delete(k)
	})
}

// unmarshalBytes unmarshals json returning an object
func unmarshalBytes(b *bytes.Buffer) (map[string]interface{}, error) {
	a := make(map[string]interface{})
	err := json.Unmarshal(b.Bytes(), &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Finds and returns a named value.
// This will return the object containing the value as well as the value
// or nill if the entry does not exist
func findJsonValue(r map[string]interface{}, n ...string) (map[string]interface{}, interface{}, bool) {
	var o map[string]interface{}
	var v interface{}
	for _, k := range n {
		if o == nil {
			o = r
		} else if a, ok := v.(map[string]interface{}); ok {
			o = a
		} else {
			// Not an object so give up
			return nil, nil, false
		}

		var e bool
		v, e = o[k]

		if !e || v == nil {
			return nil, nil, false
		}
	}
	return o, v, true
}

func GetJsonObjectValue(r map[string]interface{}, n ...string) (interface{}, bool) {
	_, v, exists := findJsonValue(r, n...)
	return v, exists
}

func GetJsonObject(o map[string]interface{}, n ...string) (map[string]interface{}, bool) {
	v, e := GetJsonObjectValue(o, n...)
	if e {
		if a, ok := v.(map[string]interface{}); ok {
			return a, true
		}
	}
	return nil, false
}

func GetJsonArray(o map[string]interface{}, n ...string) ([]interface{}, bool) {
	v, e := GetJsonObjectValue(o, n...)
	if e {
		if a, ok := v.([]interface{}); ok {
			return a, ok
		}
		var a []interface{}
		a = append(a, v)
		return a, true
	}
	return nil, false
}

// Forces an entry to be a json array. If an entry is an object or value then
// it will be wrapped within a singular array
func ForceJsonArray(r map[string]interface{}, n ...string) {
	forceJsonArray(r, n, 0, len(n)-1)
}

func forceJsonArray(r map[string]interface{}, n []string, i, j int) {
	v, e := r[n[i]]
	if e {
		if i == j {
			if _, ok := v.([]interface{}); !ok {
				var a []interface{}
				a = append(a, v)
				r[n[len(n)-1]] = a
			}
			return
		}

		if a, ok := v.([]interface{}); ok {
			for _, e := range a {
				if o, ok := e.(map[string]interface{}); ok {
					forceJsonArray(o, n, i+1, j)
				}
			}
		} else if o, ok := v.(map[string]interface{}); ok {
			forceJsonArray(o, n, i+1, j)
		}
	}
}

// If an entry is "" then replace it with {} - seen in stations feed for InformationSystems
func ForceJsonObject(r map[string]interface{}, n ...string) {
	forceJsonObject(r, n, 0, len(n)-1)
}

func forceJsonObject(r map[string]interface{}, n []string, i, j int) {
	v, e := r[n[i]]
	if e {
		if i == j {
			if s, ok := v.(string); ok {
				if s == "" {
					r[n[i]] = make(map[string]interface{})
				}
			}
			return
		}

		if a, ok := v.([]interface{}); ok {
			for _, e := range a {
				if o, ok := e.(map[string]interface{}); ok {
					forceJsonObject(o, n, i+1, j)
				}
			}
		} else if o, ok := v.(map[string]interface{}); ok {
			forceJsonObject(o, n, i+1, j)
		}
	}
}
