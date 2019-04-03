package darwinref

type CISSource struct {
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

// The CIS Source
func (r *DarwinReference) getCISSource(s string) string {
	if val, ok := r.cisSource[s]; ok {
		return val
	}
	return ""
}
