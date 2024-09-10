package models

import (
	"encoding/xml"

	"github.com/scheiblingco/gofn/errors"
)

type PersonpostRequestBestallning struct {
	OrgNr          string `xml:"OrgNr"`
	BestallningsId string `xml:"BestallningsId"`
}

type PersonpostRequest struct {
	XMLName     xml.Name                     `xml:"http://xmls.skatteverket.se/se/skatteverket/folkbokforing/na/epersondata/V1 PersonpostRequest"`
	Bestallning PersonpostRequestBestallning `xml:"Bestallning"`
	PersonId    string                       `xml:"PersonId"`
}

type PersonpostRequestBody struct {
	XMLName           xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	PersonpostRequest PersonpostRequest
}

type PersonpostRequestV4 struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    PersonpostRequestBody
}

func NewPersonpostV4Request(orgNr, bestId, personId string) (*PersonpostRequestV4, error) {
	if orgNr == "" || bestId == "" || personId == "" {
		return nil, errors.MissingValueError("orgNr, bestId or personId")
	}

	if len(orgNr) != 12 {
		return nil, errors.InvalidFieldError("orgNr must be 12 characters")
	}

	if len(personId) != 12 {
		return nil, errors.InvalidFieldError("personId must be 12 characters")
	}

	return &PersonpostRequestV4{
		Body: PersonpostRequestBody{
			PersonpostRequest: PersonpostRequest{
				Bestallning: PersonpostRequestBestallning{
					OrgNr:          orgNr,
					BestallningsId: bestId,
				},
				PersonId: personId,
			},
		},
	}, nil
}
