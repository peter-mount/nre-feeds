package darwinkb

import (
	"bytes"
	xj "github.com/peter-mount/goxml2json"
	"log"
	"os"
)

func (r *DarwinKB) xml2json(xmlFile, jsonFile string) (*bytes.Buffer, error) {

	log.Println("Converting", xmlFile, "to json")

	dir := r.config.KB.DataDir + "static/"

	in, err := os.Open(dir + xmlFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer in.Close()

	json, err := xj.Convert(in)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if jsonFile != "" {
		log.Println("Writing", jsonFile)
		out, err := os.OpenFile(dir+jsonFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer out.Close()

		out.WriteString(json.String())
	}

	log.Println("Done")
	return json, nil
}
