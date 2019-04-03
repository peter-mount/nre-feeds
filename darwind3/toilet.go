package darwind3

// The availability of a toilet in coach formation data.
// If no availability is supplied, it should be assumed to have the value "Unknown".
// Defined in rttiPPTFormations_v2
type Toilet struct {
	// An indication of the availability of a toilet in a coach in a train formation.
	// E.g. "Unknown", "None" , "Standard" or "Accessible".
	// Note that other values may be supplied in the future without a schema change.
	Type string `json:"type,omitempty" xml:",chardata,omitempty"`
	// The service status of this toilet. E.g. "Unknown", "InService" or "NotInService".
	// Default if blank "InService".
	Status string `json:"status,omitempty" xml:"status,attr,omitempty"`
}
