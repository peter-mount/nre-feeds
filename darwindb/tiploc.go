package darwindb

import (
	"github.com/peter-mount/nrod-cif/cif"
)

func (s *DarwinDB) GetCrsTiploc(crs string) (cif.Tiploc, error) {
	var t cif.Tiploc
	err := s.db.QueryRow("select tiploc,crs,name from timetable.tiploc where crs = $1 order by stanox limit 1", crs).
		Scan(&t.Tiploc, &t.CRS, &t.Desc)
	return t, err
}

func (s *DarwinDB) GetTiploc(tiploc string) (cif.Tiploc, error) {
	var t cif.Tiploc
	err := s.db.QueryRow("select tiploc,crs,name from timetable.tiploc where tiploc = $1 limit 1", tiploc).
		Scan(&t.Tiploc, &t.CRS, &t.Desc)
	return t, err
}

func (s *DarwinDB) GetTiplocs() ([]cif.Tiploc, error) {

	rows, err := s.db.Query("select tiploc,trim(crs) as crs,name from timetable.tiploc order by tiploc,name")
	if err != nil {
		return nil, err
	}

	var tiplocs []cif.Tiploc
	for rows.Next() {
		t := cif.Tiploc{}
		err = rows.Scan(&t.Tiploc, &t.CRS, &t.Desc)
		if err != nil {
			return nil, err
		}

		tiplocs = append(tiplocs, t)
	}

	return tiplocs, nil
}
