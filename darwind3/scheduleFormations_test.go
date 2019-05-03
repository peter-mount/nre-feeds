package darwind3

import (
	"encoding/xml"
	"log"
	"testing"
)

const (
	scheduleFormationXml = "<scheduleFormations rid=\"201905037116859\"><ns4:formation fid=\"201905037116859-001\"><ns4:coaches><ns4:coach coachNumber=\"0\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">None</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"1\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">Accessible</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"2\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">Standard</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"3\" coachClass=\"Mixed\"><ns4:toilet status=\"Unknown\">None</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"4\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">None</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"5\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">Accessible</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"6\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">Standard</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"7\" coachClass=\"Mixed\"><ns4:toilet status=\"Unknown\">None</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"8\" coachClass=\"Mixed\"><ns4:toilet status=\"Unknown\">None</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"9\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">Standard</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"10\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">Accessible</ns4:toilet></ns4:coach><ns4:coach coachNumber=\"11\" coachClass=\"Standard\"><ns4:toilet status=\"Unknown\">None</ns4:toilet></ns4:coach></ns4:coaches></ns4:formation></scheduleFormations>"
)

// Test unmarshalling a formation element in a uR
func Test_ScheduleFormation_XML(t *testing.T) {
	var formation ScheduleFormation

	err := xml.Unmarshal([]byte(scheduleFormationXml), &formation)
	if err != nil {
		t.Fatalf("Unmarshall: %v", err)
	}

	loadingTestEquals(t, "RID", "201905037116859", formation.RID)
	loadingTestEquals(t, "Fid", "201905037116859-001", formation.Formation.Fid)

	loadingTestEquals(t, "Coach length", 12, len(formation.Formation.Coaches))

	for i, c := range formation.Formation.Coaches {
		log.Println(i, c)
	}
}
