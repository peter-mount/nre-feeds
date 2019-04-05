package darwind3

import (
	"bytes"
	"encoding/xml"
	"testing"
	"time"
)

const (
	issue6Xml = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Pport xmlns=\"http://www.thalesgroup.com/rtti/PushPort/v16\" xmlns:ns2=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v3\" xmlns:ns3=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v2\" xmlns:ns4=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v2\" xmlns:ns5=\"http://www.thalesgroup.com/rtti/PushPort/Forecasts/v3\" xmlns:ns6=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v1\" xmlns:ns7=\"http://www.thalesgroup.com/rtti/PushPort/StationMessages/v1\" xmlns:ns8=\"http://www.thalesgroup.com/rtti/PushPort/TrainAlerts/v1\" xmlns:ns9=\"http://www.thalesgroup.com/rtti/PushPort/TrainOrder/v1\" xmlns:ns10=\"http://www.thalesgroup.com/rtti/PushPort/TDData/v1\" xmlns:ns11=\"http://www.thalesgroup.com/rtti/PushPort/Alarms/v1\" xmlns:ns12=\"http://thalesgroup.com/RTTI/PushPortStatus/root_1\" ts=\"2019-04-05T11:20:18.8416587+01:00\" version=\"16.0\"><uR updateOrigin=\"CIS\" requestSource=\"at01\" requestID=\"0000000000014065\"><trainOrder tiploc=\"BRGEND\" crs=\"BGN\" platform=\"1\"><ns9:set><ns9:first><ns9:trainID>6B09</ns9:trainID></ns9:first></ns9:set></trainOrder></uR></Pport>"
)

// Test_Issue6 looks to ensure that we can decode an xml message correctly.
// issue6Xml contains a live instance of this where the xml parser failed to process this line
func Test_Issue6_TrainOrder_XML_Parse(t *testing.T) {

	type pport struct {
		TS time.Time `xml:"ts,attr"`
		UR struct {
			Order trainOrderWrapper `xml:"trainOrder"`
		} `xml:"uR"`
	}
	p := &pport{}

	r := bytes.NewReader([]byte(issue6Xml))
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(p)

	if err != nil {
		t.Errorf("Failed to parse xml")
	}

	obj := p.UR.Order

	if obj.Tiploc != "BRGEND" {
		t.Errorf("Tiploc \"%s\" expected \"BRGEND\"", obj.Tiploc)
	}

	if obj.CRS != "BGN" {
		t.Errorf("Crs\"%s\" expected \"BGN\"", obj.CRS)
	}

	if obj.Platform != "1" {
		t.Errorf("Platform \"%s\" expected \"1\"", obj.Platform)
	}

	if obj.Set == nil {
		t.Error("Set is nil")
	} else if obj.Set.First == nil {
		t.Error("Set.First is nil")
	} else {
		f := obj.Set.First

		if f.RID != "" {
			t.Errorf("RID \"%s\" expected \"\"", f.RID)
		}

		// TODO test f.Times

		if f.TrainId != "6B09" {
			t.Errorf("TrainId \"%s\" expected \"6B09\"", f.TrainId)
		}
	}

	if obj.Set.Second != nil {
		t.Error("Set.Second is not nil")
	}

	if obj.Set.Third != nil {
		t.Error("Set.Third is not nil")
	}
}
