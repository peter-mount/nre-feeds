package darwind3

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"sort"
	"testing"
)

// create new location
// i 0 use pta, 1 ptd, 2 wta & 3 wtd, 4 becomes 0 as %4 used
//
func location_new(i int, s string) *Location {
	var v *Location = &Location{
		Type:   "TS",
		Tiploc: "MSTONEE",
	}
	switch i % 4 {
	case 0:
		v.Times.Pta = util.NewPublicTime(s)
		v.Forecast.Arrival.ET = &util.WorkingTime{}
		v.Forecast.Arrival.ET.Set(v.Times.Pta.Get() * 60)
	case 1:
		v.Times.Ptd = util.NewPublicTime(s)
		v.Forecast.Departure.AT = &util.WorkingTime{}
		v.Forecast.Departure.AT.Set(v.Times.Ptd.Get() * 60)
	case 2:
		v.Times.Wta = util.NewWorkingTime(s)
		v.Forecast.Arrival.ET = v.Times.Wta
	case 3:
		v.Times.Wtd = util.NewWorkingTime(s)
		v.Forecast.Departure.AT = v.Times.Wtd
	}
	v.UpdateTime()

	return v
}

func location_testBool(t *testing.T, m string, e bool, f func() bool) {
	v := f()
	if v != e {
		t.Errorf("%s: got %v expected %v", m, v, e)
	}
}

// Test Location.Compare() works correctly
func TestLocation_Compare(t *testing.T) {

	a := location_new(0, "01:02")
	b := location_new(1, "02:03")

	// a < b
	location_testBool(t, "a.Compare(b)", true, func() bool {
		return a.Compare(b)
	})

	// b > b
	location_testBool(t, "b.Compare(a)", false, func() bool {
		return b.Compare(a)
	})

	// a = a & b = b
	location_testBool(t, "a.Compare(a)", false, func() bool {
		return a.Compare(a)
	})
	location_testBool(t, "b.Compare(b)", false, func() bool {
		return b.Compare(b)
	})

}

// Test that slices sort correctly
func location_timesSlice() []*Location {
	var ary []*Location
	var location_times = [...]string{
		"09:55",
		"09:50",
		"09:14",
		"09:18",
		"09:39",
		"09:33",
		"09:37",
		"10:14",
		"09:32",
		"09:40",
		"09:47",
		"09:52",
		"09:25",
		"19:50",
		"19:14",
		"19:18",
		"19:55",
		"19:39",
		"19:33",
		"19:37",
		"20:14",
		"19:32",
		"19:40",
		"19:47",
		"19:52",
		"19:25",
	}

	for i, s := range location_times {
		ary = append(ary, location_new(i, s))
	}
	return ary
}

func pary(l string, a []*Location) {
	var ary []string
	for _, av := range a {
		ary = append(ary, av.Time.String())
	}
	log.Println(l, ary)
}

func TestLocation_SliceStable(t *testing.T) {

	ary := location_timesSlice()

	pary("before sort", ary)

	sort.SliceStable(ary, func(i, j int) bool {
		return ary[i].Compare(ary[j])
	})

	pary("after sort", ary)

	var l *util.WorkingTime
	for i, v := range ary {
		if i > 0 && v.Time.Before(l) {
			t.Errorf("Element %d not in correct place. Last %v Got %v", i, l.String(), v.Time.String())
		}
		l = &v.Time
	}

}

const (
	testlocationXmlparseLengthXml      = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><ns5:Location tpl=\"ELMW\" wta=\"21:00:30\" wtd=\"21:01\" pta=\"21:01\" ptd=\"21:01\"><ns5:arr et=\"21:01\" wet=\"21:00\" src=\"Darwin\"/><ns5:dep et=\"21:01\" src=\"Darwin\"/><ns5:plat>4</ns5:plat><ns5:length>8</ns5:length></ns5:Location>"
	testlocationXmlparseLengthExpected = 8
)

// Test that the train length is parsed correctly
func TestLocation_XMLParse_Length(t *testing.T) {
	var loc Location

	err := xml.Unmarshal([]byte(testlocationXmlparseLengthXml), &loc)
	if err != nil {
		t.Fatalf("Length unmarshall: %v", err)
	}

	log.Println("Train length", loc.Length)

	if loc.Length != testlocationXmlparseLengthExpected {
		t.Errorf("XML Length, expected %d got %d", testlocationXmlparseLengthExpected, loc.Length)
	}

}
