// Reference timetable
package darwintimetable

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
