package darwind3

import (
	"encoding/xml"
	"testing"
)

func TestTracingIdXML(t *testing.T) {
	src := "<trackingID><ns10:berth area=\"EH\">0734</ns10:berth><ns10:incorrectTrainID>1J04</ns10:incorrectTrainID><ns10:correctTrainID>1J23</ns10:correctTrainID></trackingID>"

	var id TrackingID

	if err := xml.Unmarshal([]byte(src), &id); err != nil {
		t.Errorf("Failed to unmarshal pport xml: %v", err)
	}

	if id.Berth.Area != "EH" {
		t.Errorf("area incorrect: %s", id.Berth.Area)
	}

	if id.Berth.Berth != "0734" {
		t.Errorf("berth incorrect: %s", id.Berth.Berth)
	}

	if id.IncorrectTrainID != "1J04" {
		t.Errorf("incorrectTrainID incorrect: %s", id.IncorrectTrainID)
	}

	if id.CorrectTrainID != "1J23" {
		t.Errorf("correctTrainID incorrect: %s", id.CorrectTrainID)
	}
}
