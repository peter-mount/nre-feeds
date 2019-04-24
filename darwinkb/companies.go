package darwinkb

import (
	"github.com/peter-mount/golib/kernel/bolt"
	"github.com/peter-mount/sortfold"
	"log"
	"sort"
)

const (
	companiesXml  = "companies.xml"
	companiesJson = "companies.json"
)

// An entry in the incident index
type CompanyEntry struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (r *DarwinKB) GetCompanies() ([]byte, error) {
	// Works as we have the index as a single key
	b, err := r.GetCompany("index")
	return b, err
}

func (r *DarwinKB) GetCompany(id string) ([]byte, error) {
	var data []byte
	err := r.View(tocsBucket, func(bucket *bolt.Bucket) error {
		data = bucket.Get(id)
		return nil
	})
	return data, err
}

func (r *DarwinKB) refreshCompanies() {
	err := r.refreshCompaniesImpl()
	if err != nil {
		log.Println("refreshIncidents:", err)
	}
}

func (r *DarwinKB) refreshCompaniesImpl() error {

	updateRequired, err := r.refreshFile(companiesXml, tocsUrl, tocsMaxAge)
	if err != nil {
		return err
	}

	// If no update check to see if the bucket is empty forcing an update
	if !updateRequired {
		updateRequired, err = r.bucketEmpty(tocsBucket)
		if err != nil {
			return err
		}
	}

	// Give up if no update is required
	if !updateRequired {
		return nil
	}

	b, err := r.xml2json(companiesXml, companiesJson)
	if err != nil {
		return err
	}

	log.Println("Parsing JSON")

	root, err := unmarshalBytes(b)
	if err != nil {
		return err
	}

	var index []*CompanyEntry

	companies, _ := GetJsonArray(root, "TrainOperatingCompanyList", "TrainOperatingCompany")
	log.Println("Found", len(companies), tocsBucket)

	err = r.Update(tocsBucket, func(bucket *bolt.Bucket) error {
		err := bucketRemoveAll(bucket)
		if err != nil {
			return err
		}

		for _, company := range companies {
			o := company.(map[string]interface{})

			atocCode, _ := GetJsonObjectValue(o, "AtocCode")
			name, _ := GetJsonObjectValue(o, "Name")

			// The individual entry
			err = bucket.PutJSON(atocCode.(string), company)
			if err != nil {
				return err
			}

			index = append(index, &CompanyEntry{Id: atocCode.(string), Name: name.(string)})
		}

		sort.SliceStable(index, func(i, j int) bool { return sortfold.CompareFold(index[i].Name, index[j].Name) < 0 })

		err = bucket.PutJSON("index", index)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("Updated %d companies", len(companies))
	return nil
}
