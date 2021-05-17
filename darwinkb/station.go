package darwinkb

import (
	"github.com/peter-mount/go-kernel/bolt"
	"log"
)

const (
	stationXml  = "station.xml"
	stationJson = "station.json"
)

func (r *DarwinKB) GetStation(crs string) ([]byte, error) {
	var data []byte
	err := r.View(stationsBucket, func(bucket *bolt.Bucket) error {
		data = bucket.Get(crs)
		return nil
	})
	return data, err
}

func (r *DarwinKB) refreshStations() {
	err := r.refreshStationsImpl()
	if err != nil {
		log.Println("refreshStations:", err)
	}
}

func (r *DarwinKB) refreshStationsImpl() error {

	updateRequired, err := r.refreshFile(stationXml, stationsUrl, stationsMaxAge)
	if err != nil {
		return err
	}

	// If no update check to see if the bucket is empty forcing an update
	if !updateRequired {
		updateRequired, err = r.bucketEmpty(stationsBucket)
		if err != nil {
			return err
		}
	}

	// Give up if no update is required
	if !updateRequired {
		return nil
	}

	b, err := r.xml2json(stationXml, stationJson)
	if err != nil {
		return err
	}

	log.Println("Parsing JSON")

	root, err := unmarshalBytes(b)
	if err != nil {
		return err
	}

	stations, _ := GetJsonArray(root, "StationList", "Station")
	log.Printf("Found %d stations", len(stations))

	err = r.Update(stationsBucket, func(bucket *bolt.Bucket) error {
		err := bucketRemoveAll(bucket)
		if err != nil {
			return err
		}

		// Insert entries one per entry using CrsCode as the key
		for _, c := range stations {

			d := c.(map[string]interface{})
			crs := d["CrsCode"].(string)

			// Fix certain fields into correct objects,
			// e.g. single values into an array if the element is "unbounded"
			// or "" into {}

			// This one isn't in the feed but as it's "unbounded" then we need to fix it
			ForceJsonArray(d, "AlternativeIdentifiers", "Tiplocs", "Tiploc")

			ForceJsonArray(d, "Fares", "PenaltyFares", "TrainOperator")
			ForceJsonArray(d, "Accessibility", "StaffHelpAvailable", "Open", "DayAndTimeAvailability")
			ForceJsonArray(d, "Accessibility", "StaffHelpAvailable", "Open", "DayAndTimeAvailability", "OpeningHours", "OpenPeriod")

			ForceJsonArray(d, "Fares", "PenaltyFares", "TrainOperator")
			ForceJsonArray(d, "Fares", "TicketOffice", "Open", "DayAndTimeAvailability")
			ForceJsonArray(d, "Fares", "TicketOffice", "Open", "DayAndTimeAvailability", "OpeningHours", "OpenPeriod")
			ForceJsonArray(d, "Fares", "Travelcard", "TravelcardZone")

			ForceJsonObject(d, "InformationSystems")
			ForceJsonArray(d, "InformationSystems", "CIS")
			ForceJsonArray(d, "InformationSystems", "DayAndTimeAvailability")
			ForceJsonArray(d, "InformationSystems", "DayAndTimeAvailability", "OpeningHours", "OpenPeriod")
			ForceJsonArray(d, "InformationSystems", "InformationServicesOpen", "DayAndTimeAvailability")
			ForceJsonArray(d, "InformationSystems", "InformationServicesOpen", "DayAndTimeAvailability", "OpeningHours", "OpenPeriod")
			ForceJsonArray(d, "InformationSystems", "InformationAvailableFromStaff")

			ForceJsonArray(d, "Interchange", "CarPark")
			ForceJsonArray(d, "Interchange", "CarPark", "Open", "DayAndTimeAvailability")
			ForceJsonArray(d, "Interchange", "CarPark", "Open", "DayAndTimeAvailability", "OpeningHours", "OpenPeriod")
			ForceJsonArray(d, "Interchange", "CycleStorage", "Type")
			ForceJsonArray(d, "Interchange", "RailReplacementServices", "RailReplacementMap")

			ForceJsonArray(d, "TrainOperatingCompanies", "TocRef")

			err = bucket.PutJSON(crs, d)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("Updated %d stations", len(stations))
	return nil
}
