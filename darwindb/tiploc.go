package darwindb

import "time"

type Tiploc struct {
	Id          int       `json:"id"`
	Tiploc      string    `json:"tiploc"`
	Crs         string    `json:"crs"`
	Stanox      int       `json:"stanox"`
	Name        string    `json:"name"`
	Nlc         int       `json:"nlc"`
	NlcCheck    string    `json:"nlccheck"`
	NlcDesc     string    `json:"nlcdesc"`
	Station     bool      `json:"station"`
	Deleted     bool      `json:"deleted"`
	DateExtract time.Time `json:"dateextract"`
}

func (s *DarwinDB) GetCrsTiploc(crs string) (Tiploc, error) {
	var t Tiploc
	err := s.db.QueryRow("select tiploc,crs,name from timetable.tiploc where crs = $1 order by stanox limit 1", crs).
		Scan(&t.Tiploc, &t.Crs, &t.Name)
	return t, err
}

func (s *DarwinDB) GetTiploc(tiploc string) (Tiploc, error) {
	var t Tiploc
	err := s.db.QueryRow("select tiploc,crs,name from timetable.tiploc where tiploc = $1 limit 1", tiploc).
		Scan(&t.Tiploc, &t.Crs, &t.Name)
	return t, err
}

func (s *DarwinDB) GetTiplocs() ([]Tiploc, error) {

	rows, err := s.db.Query("select tiploc,trim(crs) as crs,name from timetable.tiploc order by tiploc,name")
	if err != nil {
		return nil, err
	}

	var tiplocs []Tiploc
	for rows.Next() {
		t := Tiploc{}
		err = rows.Scan(&t.Tiploc, &t.Crs, &t.Name)
		if err != nil {
			return nil, err
		}

		tiplocs = append(tiplocs, t)
	}

	return tiplocs, nil
}
