package models

import (
	"encoding/xml"
)

type PersonpostResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName               xml.Name `xml:"Body"`
		PersonpostXMLResponse struct {
			XMLName              xml.Name             `xml:"PersonpostXMLResponse"`
			Folkbokforingsposter Folkbokforingsposter `xml:"Folkbokforingsposter"`
		}
	}
}
