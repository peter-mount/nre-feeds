package darwind3

import (
	"encoding/xml"
	"github.com/peter-mount/filecache"
	"log"
	"testing"
	"time"
)

const (
	testAlarmSetTsXml = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Pport xmlns=\"http://www.thalesgroup.com/rtti/PushPort/v16\" xmlns:ns2=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v3\" xmlns:ns3=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v2\" xmlns:ns4=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v2\" xmlns:ns5=\"http://www.thalesgroup.com/rtti/PushPort/Forecasts/v3\" xmlns:ns6=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v1\" xmlns:ns7=\"http://www.thalesgroup.com/rtti/PushPort/StationMessages/v1\" xmlns:ns8=\"http://www.thalesgroup.com/rtti/PushPort/TrainAlerts/v1\" xmlns:ns9=\"http://www.thalesgroup.com/rtti/PushPort/TrainOrder/v1\" xmlns:ns10=\"http://www.thalesgroup.com/rtti/PushPort/TDData/v1\" xmlns:ns11=\"http://www.thalesgroup.com/rtti/PushPort/Alarms/v1\" xmlns:ns12=\"http://thalesgroup.com/RTTI/PushPortStatus/root_1\" ts=\"2019-04-05T09:49:28.1704587+01:00\" version=\"16.0\"><uR updateOrigin=\"Darwin\"><alarm><ns11:set id=\"223686\"><ns11:tdAreaFail>SX</ns11:tdAreaFail></ns11:set></alarm></uR></Pport>"
)

func Test_Alarm_Unmarshal_Set_TS(t *testing.T) {
	//test_alarm_set_ts_xml := "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Pport xmlns=\"http://www.thalesgroup.com/rtti/PushPort/v16\" xmlns:ns2=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v3\" xmlns:ns3=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v2\" xmlns:ns4=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v2\" xmlns:ns5=\"http://www.thalesgroup.com/rtti/PushPort/Forecasts/v3\" xmlns:ns6=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v1\" xmlns:ns7=\"http://www.thalesgroup.com/rtti/PushPort/StationMessages/v1\" xmlns:ns8=\"http://www.thalesgroup.com/rtti/PushPort/TrainAlerts/v1\" xmlns:ns9=\"http://www.thalesgroup.com/rtti/PushPort/TrainOrder/v1\" xmlns:ns10=\"http://www.thalesgroup.com/rtti/PushPort/TDData/v1\" xmlns:ns11=\"http://www.thalesgroup.com/rtti/PushPort/Alarms/v1\" xmlns:ns12=\"http://thalesgroup.com/RTTI/PushPortStatus/root_1\" ts=\"2019-04-05T09:49:28.1704587+01:00\" version=\"16.0\"><uR updateOrigin=\"Darwin\"><alarm><ns11:set id=\"223686\"><ns11:tdAreaFail>SX</ns11:tdAreaFail></ns11:set></alarm></uR></Pport>"

	type pport struct {
		TS time.Time `xml:"ts,attr"`
		UR struct {
			Alarm RttiAlarm `xml:"alarm"`
		} `xml:"uR"`
	}
	wrapper := &pport{}

	err := xml.Unmarshal([]byte(testAlarmSetTsXml), wrapper)
	if err != nil {
		t.Errorf("Failed to parse xml: %v", err)
		return
	}

	obj := &wrapper.UR.Alarm

	if obj.Set.ID != "223686" {
		t.Errorf("ID is \"%s\" expecting \"223686\"", obj.Set.ID)
	}

	if obj.Set.TDArea != "SX" {
		t.Errorf("TDArea is \"%s\" expecting \"SX\"", obj.Set.TDArea)
	}

	if obj.Set.Tyrell != "" {
		t.Errorf("Tyrell is \"%s\" expecting \"\"", obj.Set.Tyrell)
	}

	if obj.Clear != "" {
		t.Errorf("Clear is \"%s\" expecting \"\"", obj.Clear)
	}

	// The Alarm object with date set
	alarm := &obj.Set
	alarm.Date = wrapper.TS

	// Marshal to Bytes
	b := filecache.ToJsonBytes(&alarm)
	if err != nil {
		t.Errorf("Bytes Failed: %v", err)
		return
	}
	if b == nil || len(b) == 0 {
		t.Error("No bytes returned")
	}
	log.Println(string(b))
}
