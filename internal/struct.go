package internal

import (
	"encoding/xml"
	"net/http"
)

type Session struct {
	Client              *http.Client
	AuthSessionID       string
	AuthSessionIDLegacy string
	KCRestart           string
	SessionCode         string
	Execution           string
	TabID               string
}

type StudentDetails struct {
	NIM                string
	FullName           string
	Email              string
	Faculty            string
	StudyProgram       string
	SIAKADPhotoURL     string
	FileFILKOMPhotoUrl string
}

type SAMLResponse struct {
	XMLName   xml.Name  `xml:"Response"`
	Assertion Assertion `xml:"Assertion"`
}

type Assertion struct {
	XMLName            xml.Name           `xml:"Assertion"`
	AttributeStatement AttributeStatement `xml:"AttributeStatement"`
}

type AttributeStatement struct {
	XMLName    xml.Name    `xml:"AttributeStatement"`
	Attributes []Attribute `xml:"Attribute"`
}

type Attribute struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"AttributeValue"`
}
