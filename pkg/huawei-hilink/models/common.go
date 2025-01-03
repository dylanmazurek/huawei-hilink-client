package models

import "encoding/xml"

type HilinkResponse struct {
	XMLName     xml.Name
	SessionInfo string `xml:"SesInfo"`
	TokenInfo   string `xml:"TokInfo"`
}

func (t *HilinkResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias HilinkResponse
	aux := &struct {
		*Alias
		Code    string `xml:"code"`
		Message string `xml:"message"`
	}{
		Alias: (*Alias)(t),
	}

	err := d.DecodeElement(aux, &start)
	if err != nil {
		return err
	}

	return nil
}