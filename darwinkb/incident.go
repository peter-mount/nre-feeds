package darwinkb

import (
  "time"
)

// Type for individual Incident.
type KBIncident struct {
  // Time of creation of incident
  Created             time.Time     `json:"created" xml:"CreationTime"`
  // Who changed the data most recently.
  //ChangeHistory       // ChangeHistoryStructure
  // Unique identifier of system issuing entry identifier.
  // If absent, taken from Context element.
  // May be different from that in Context, indicating that the incident is
  // forwarded from another system  - without being allocated a new
  // identifier by the inbtermediate system.
  // Note that the ExternalCode may be used to retain the external System's
  // identifier to allow round trip processing.
  ParticipantRef    []string        `json:"participantRef" xml:"ParticipantRef"`
  // Identifier of entry.
  // Must be unique within Participant's current data horizon.
  // Monotonically increasing,  seqience with time of issue.
  // Normally also unique within Participant (ie also outside of the current
  // horizon) so that a uniform namespace can also be used for
  // archived messages as well.
  IncidentNumber      string        `json:"incidentNumber" xml:"IncidentNumber"`
  // Version number if entry is update to a previous version.
  // Unique within IncidentNumber. Monotonically increasing within IncidentNumber.
  // Any values for  classification, description, affects, effects that are
  // present in an update replace any values on previous incidents and updates
  // with the same identifier.  Values that are not updated remain in effect.
  Version             int           `json:"version" xml:"Version"`
  // Twitter hash tag for the source.
  TwitterHashtag      string        `json:"twitterHashtag" xml:"Source>TwitterHashtag"`
  //OuterValidityPeriod // HalfOpenTimestampRangeStructure
  // Overall inclusive Period of applicability of incident
  //ValidityPeriod      // HalfOpenTimestampRangeStructure
  // Whether the incident was planned (eg engineering works) or
  // unplanned (eg service alteration).
  // Default is false, i.e. unplanned.
  Planned             bool          `json:"planned" xml:"Planned"`
  // >Summary of incident. If absent should be generated from structure
  // elements / and or by condensing Description.
  Summary             string        `json:"summary" xml:"Summary"`
  // Description of incident. Should not repeat any strap line incldued in Summary.
  Description         string        `json:"description" xml:"Description"`
  // Hyperlinks to other resources associated with incident.
  InfoLink          []string        `json:"links" xml:"InfoLinks>InfoLink>Uri"`
  // Structured model identifiying parts of transport network affected by incident.
  // Operator and Network values will be defaulted to values in general Context unless explicitly overridden.
  Operators         []KBAffectsOperator   `json:"operators" xml:"Affects>Operators>AffectedOperator"`
  Routes            []string              `json:"routes" xml:"Affects>RoutesAffected"`
  ClearedIncident     bool                `json:"cleared" xml:"ClearedIncident"`
  IncidentPriority    int                 `json:"priority" xml:"IncidentPriority"`
  P0Summary           string              `json:"P0Summary" xml:"P0Summary"`
}

type KBAffects struct {
}

type KBAffectsOperator struct {
  Ref       string    `json:"toc" xml:"OperatorRef"`
  Name      string    `json:"name" xml:"OperatorName"`
}
