package util

import (
	"testing"
)

func TestGenerateScheduleKey(t *testing.T) {
	wta := NewWorkingTime("12:00:00")
	wtd := NewWorkingTime("12:01:30")
	wtp := NewWorkingTime("12:03:00")

	if GenerateScheduleKey("MSTONEE", nil, wtd, nil) != "MSTONEEmtybgemty" {
		t.Error("Failed origin")
	}

	if GenerateScheduleKey("MSTONEE", wta, wtd, nil) != "MSTONEEbeMbgemty" {
		t.Error("Failed stop")
	}

	if GenerateScheduleKey("MSTONEE", nil, nil, wtp) != "MSTONEEmtymtybhG" {
		t.Error("Failed pass")
	}

	if GenerateScheduleKey("MSTONEE", wta, nil, nil) != "MSTONEEbeMmtymty" {
		t.Error("Failed termination")
	}

}
