// Reference timetable
package darwintimetable

type Timetable struct {
  TimetableId     string              `xml:"timetableID,attr"`
  //Journeys      []*Journey            `xml:"Journey"`
  Journeys        map[string]*Journey `xml:"-"`
}

type Journey struct {
  RID             string        `xml:"rid,attr"`
  UID             string        `xml:"uid,attr"`
  TrainID         string        `xml:"trainId"`
  SSD             string        `xml:"ssd,attr"`
  Toc             string        `xml:"toc,attr"`
  TrainCat        string        `xml:"trainCat,attr"`
  Passenger       bool          `xml:"isPassengerSvc,attr"`
  // The schedule
  Schedule      []interface{}   `xml:,any`
  CancelReason    int           `xml:"cancelReason"`
  // Associations
  Associations  []*Association  `xml:"-"`
}

type OR struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // False destination to be used at this location
  FalseDest string      `xml:"fd,attr"`
}

type OPOR struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
}

type IP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
  // False destination to be used at this location
  FalseDest string      `xml:"fd,attr"`
}

type OPIP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type PP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wtp       string      `xml:"wtp,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type DT struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type OPDT struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type cancelReason struct {
  text string `xml:",chardata"`
}

type Association struct {
  Main      AssocService  `xml:"main"`
  Assoc     AssocService  `xml:"assoc"`
  Tiploc    string        `xml:"tiploc,attr"`
  Category  string        `xml:"category,attr"`
  Cancelled bool          `xml:"isCancelled,attr"`
  Deleted   bool          `xml:"isDeleted,attr"`
}

type AssocService struct {
  RID   string    `xml:"rid,attr"`
  Wta   string      `xml:"wta,attr"`
  Wtd   string      `xml:"wtd,attr"`
  Wtp   string      `xml:"wtp,attr"`
  Pta   string      `xml:"pta,attr"`
  Ptd   string      `xml:"ptd,attr"`
}
