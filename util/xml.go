package util

import "encoding/xml"

type xmlIntType struct {
	Value int `xml:",chardata"`
}

// DecodeXML_Int parses a simple XML element of type <element>integer</element> and returns the integer value.
// It's used within custom xml.UnmarshalXML() code.
func DecodeXML_Int(decoder *xml.Decoder, token *xml.StartElement) (int, error) {
	intType := xmlIntType{}
	if err := decoder.DecodeElement(&intType, token); err != nil {
		return 0, err
	}
	return intType.Value, nil
}

type xmlBoolType struct {
	Value bool `xml:",chardata"`
}

// DecodeXML_Bool parses a simple XML element of type <element>bool</element> and returns the bool value.
// It's used within custom xml.UnmarshalXML() code.
func DecodeXML_Bool(decoder *xml.Decoder, token *xml.StartElement) (bool, error) {
	boolType := xmlBoolType{}
	if err := decoder.DecodeElement(&boolType, token); err != nil {
		return false, err
	}
	return boolType.Value, nil
}
