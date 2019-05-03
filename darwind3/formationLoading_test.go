package darwind3

import (
	"encoding/xml"
	"fmt"
	"github.com/peter-mount/nre-feeds/util"
	"strconv"
	"testing"
)

const (
	formationLoadingURXml = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><uR ts=\"2019-05-03T11:03:17.4837588+01:00\" version=\"16.0\" updateOrigin=\"CIS\" requestSource=\"CACI\" requestID=\"lv05031003179\"><formationLoading fid=\"201905037118527-001\" rid=\"201905037118527\" tpl=\"VICTRIC\" wta=\"00:00\" wtd=\"11:00\" pta=\"00:00\" ptd=\"11:00\"><ns6:loading coachNumber=\"0\">8</ns6:loading><ns6:loading coachNumber=\"1\">12</ns6:loading><ns6:loading coachNumber=\"2\">15</ns6:loading><ns6:loading coachNumber=\"3\">40</ns6:loading><ns6:loading coachNumber=\"4\">26</ns6:loading><ns6:loading coachNumber=\"5\">66</ns6:loading><ns6:loading coachNumber=\"6\">53</ns6:loading><ns6:loading coachNumber=\"7\">88</ns6:loading></formationLoading></uR>"
	formationLoadingXml   = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><formationLoading fid=\"201905037118527-001\" rid=\"201905037118527\" tpl=\"VICTRIC\" wta=\"00:00\" wtd=\"11:00\" pta=\"00:00\" ptd=\"11:00\"><ns6:loading coachNumber=\"0\">8</ns6:loading><ns6:loading coachNumber=\"1\">12</ns6:loading><ns6:loading coachNumber=\"2\">15</ns6:loading><ns6:loading coachNumber=\"3\">40</ns6:loading><ns6:loading coachNumber=\"4\">26</ns6:loading><ns6:loading coachNumber=\"5\">66</ns6:loading><ns6:loading coachNumber=\"6\">53</ns6:loading><ns6:loading coachNumber=\"7\">88</ns6:loading></formationLoading>"
)

// Test unmarshalling a formation element in a uR
func Test_FormationLoading_UR(t *testing.T) {
	var ur UR

	err := xml.Unmarshal([]byte(formationLoadingURXml), &ur)
	if err != nil {
		t.Fatalf("Unmarshall: %v", err)
	}

}
func loadingTestEquals(t *testing.T, n string, e, v interface{}) {
	if e != v {
		t.Errorf("%s expected \"%v\" got \"%v\"", n, e, v)
	}
}

// Test unmarshalling a formation element in a uR
func Test_FormationLoading_XML(t *testing.T) {
	var loading Loading

	err := xml.Unmarshal([]byte(formationLoadingXml), &loading)
	if err != nil {
		t.Fatalf("Unmarshall: %v", err)
	}

	loadingTestEquals(t, "Fid", "201905037118527-001", loading.Fid)
	loadingTestEquals(t, "RID", "201905037118527", loading.RID)
	loadingTestEquals(t, "Tiploc", "VICTRIC", loading.Tiploc)
	loadingTestEquals(t, "Coaches", 8, len(loading.Loading))

	for i, l := range []int{8, 12, 15, 40, 26, 66, 53, 88} {
		loadingTestEquals(t, fmt.Sprintf("Coach %d: number", i), strconv.FormatUint(uint64(i), 10), loading.Loading[i].CoachNumber)
		loadingTestEquals(t, fmt.Sprintf("Coach %d: loading", i), l, loading.Loading[i].Loading)
	}

}

func Test_FormationLoading_AppendFormationLoading(t *testing.T) {
	var loading Loading

	err := xml.Unmarshal([]byte(formationLoadingXml), &loading)
	if err != nil {
		t.Fatalf("Unmarshall: %v", err)
	}

	loc := &Location{
		Tiploc: loading.Tiploc,
		Times: util.CircularTimes{
			Wtd: util.NewWorkingTime("11:00"),
			Ptd: util.NewPublicTime("11:00"),
		},
	}

	sched := Schedule{
		RID: loading.RID,
		Locations: []*Location{
			loc,
		},
	}

	sched.appendFormationLoading(nil, &loading)

	if loc.Loading == nil {
		t.Error("Loading not appened to location")
	}
}
