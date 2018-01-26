package darwind3

import (
  "encoding/xml"
)

// Describes the identifier of a train in the train order
type trainOrderItem struct {
  // For trains in the train order where the train is the Darwin timetable,
  // it will be identified by its RID
  RID       string
  // One or more scheduled times to identify the instance of the location in
  // the train schedule for which the train order is set.
  Times     CircularTimes
  // Where a train in the train order is not in the Darwin timetable,
  // a Train ID (headcode) will be supplied
  TrainId   string
}

func (s *trainOrderItem) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  inRid := false
  inTrainId := false
  var rid []byte
  var tid []byte

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
      case xml.StartElement:
        switch tok.Name.Local {
          case "rid":
            inRid = true
            inTrainId = false
            s.Times.UnmarshalXMLAttributes( tok )

          case "trainID":
            inRid = false
            inTrainId = true

          default:
            if err := decoder.Skip(); err != nil {
              return err
            }
        }

      case xml.CharData:
        if inRid {
          rid = append( rid, tok... )
        } else if inTrainId {
          tid = append( rid, tok... )
        }

      case xml.EndElement:
        if inRid || inTrainId {
          inRid = false
          s.RID = string( rid )
        } else if inRid || inTrainId {
          inTrainId = false
          s.TrainId = string( tid )
        } else {
          return nil
        }
    }
  }
}
