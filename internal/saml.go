package internal

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
)

func ParseSAMLResponse(samlBase64 string) (*StudentDetails, error) {
	studentDetails := new(StudentDetails)

	decoded, err := base64.StdEncoding.DecodeString(samlBase64)
	if err != nil {
		return studentDetails, fmt.Errorf("failed to decode SAML response: %v", err)
	}

	var response SAMLResponse
	if err := xml.Unmarshal(decoded, &response); err != nil {
		return studentDetails, fmt.Errorf("failed to parse SAML XML: %v", err)
	}

	for _, attr := range response.Assertion.AttributeStatement.Attributes {
		switch attr.Name {
		case "nim":
			studentDetails.NIM = attr.Value
		case "email":
			studentDetails.Email = attr.Value
		case "fullName":
			studentDetails.FullName = PascalCase(attr.Value)
		case "fakultas":
			studentDetails.Faculty = fmt.Sprintf("Fakultas %s", attr.Value)
		case "prodi":
			studentDetails.StudyProgram = attr.Value
		}
	}

	if len(studentDetails.NIM) >= 2 {
		studentDetails.SIAKADPhotoURL = fmt.Sprintf("https://siakad.ub.ac.id/dirfoto/foto/foto_20%s/%s.jpg", studentDetails.NIM[0:2], studentDetails.NIM)
		studentDetails.FileFILKOMPhotoUrl = fmt.Sprintf("https://file-filkom.ub.ac.id/fileupload/assets/foto/20%s/%s.png", studentDetails.NIM[0:2], studentDetails.NIM)
	}

	return studentDetails, nil
}
