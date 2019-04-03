package darwinkb

import (
	"github.com/peter-mount/golib/kernel/bolt"
	"log"
	"time"
)

const (
	serviceIndicatorsXml  = "serviceIndicators.xml"
	serviceIndicatorsJson = "serviceIndicators.json"
)

func (r *DarwinKB) GetServiceIndicators() ([]byte, error) {
	// Works as we have the index as a single key
	b, err := r.GetServiceIndicator("index")
	return b, err
}

func (r *DarwinKB) GetServiceIndicator(toc string) ([]byte, error) {
	var data []byte
	err := r.View("serviceIndicators", func(bucket *bolt.Bucket) error {
		data = bucket.Get(toc)
		return nil
	})
	return data, err
}

func (r *DarwinKB) refreshServiceIndicators() {
	err := r.refreshServiceIndicatorsImpl()
	if err != nil {
		log.Println("refreshServiceIndicators:", err)
	}
}

func (r *DarwinKB) refreshServiceIndicatorsImpl() error {

	updateRequired, err := r.refreshFile(serviceIndicatorsXml, "https://datafeeds.nationalrail.co.uk/api/staticfeeds/4.0/serviceIndicators", 9*time.Minute)
	if err != nil {
		return err
	}

	// If no update check to see if the bucket is empty forcing an update
	if !updateRequired {
		updateRequired, err = r.bucketEmpty("serviceIndicators")
		if err != nil {
			return err
		}
	}

	// Give up if no update is required
	if !updateRequired {
		return nil
	}

	b, err := r.xml2json(serviceIndicatorsXml, serviceIndicatorsJson)
	if err != nil {
		return err
	}

	log.Println("Parsing JSON")

	root, err := unmarshalBytes(b)
	if err != nil {
		return err
	}

	// Force all ServiceGroup's into arrays
	ForceJsonArray(root, "NSI", "TOC", "ServiceGroup")

	serviceIndicators, _ := GetJsonArray(root, "NSI", "TOC")
	log.Println("Found", len(serviceIndicators), "serviceIndicators")

	err = r.Update("serviceIndicators", func(bucket *bolt.Bucket) error {
		err := bucketRemoveAll(bucket)
		if err != nil {
			return err
		}

		err = bucket.PutJSON("index", serviceIndicators)
		if err != nil {
			return err
		}

		for _, status := range serviceIndicators {
			o := status.(map[string]interface{})

			tocCode, _ := GetJsonObjectValue(o, "TocCode")

			// The individual entry
			err = bucket.PutJSON(tocCode.(string), status)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("Updated %d serviceIndicators", len(serviceIndicators))
	return nil
}
