package ub_auth

import (
	"fmt"
	"github.com/ahmdyaasiin/ub-auth-without-notification/v3/internal"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Auth(username, password string) (*internal.StudentDetails, error) {
	studentDetails := new(internal.StudentDetails)

	session, err := internal.GetSession()
	if err != nil {
		return studentDetails, err
	}

	data := strings.NewReader(fmt.Sprintf("username=%s&password=%s&credentialId=%s", url.QueryEscape(username), url.QueryEscape(password), ""))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://iam.ub.ac.id/auth/realms/ub/login-actions/authenticate?session_code=%s&execution=%s&client_id=brone.ub.ac.id&tab_id=%s", session.SessionCode, session.Execution, session.TabID), data)
	if err != nil {
		return studentDetails, fmt.Errorf("failed to create POST request: %v", err)
	}

	for k, v := range internal.GetHeaders() {
		req.Header.Set(k, v)
	}

	req.Header.Set("origin", "null")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", fmt.Sprintf("AUTH_SESSION_ID=%s; AUTH_SESSION_ID_LEGACY=%s; KC_RESTART=%s", session.AuthSessionID, session.AuthSessionIDLegacy, session.KCRestart))

	resp, err := session.Client.Do(req)
	if err != nil {
		return studentDetails, fmt.Errorf("failed to perform POST request: %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return studentDetails, fmt.Errorf("failed to read response body: %v", err)
	}

	r := fmt.Sprintf("%s", respBody)
	if !strings.Contains(r, "SAMLResponse") {
		if strings.Contains(r, "Invalid username or password.") {

			return studentDetails, fmt.Errorf("invalid username or password")
		} else {

			return studentDetails, fmt.Errorf("unexpected error")
		}
	}

	samlResponse, err := internal.GetSubstringBetween(`name="SAMLResponse" value="`, `"/>`, r)
	if err != nil {
		return studentDetails, fmt.Errorf("failed to get SAML response")
	}

	studentDetails, err = internal.ParseSAMLResponse(samlResponse)
	if err != nil {
		return studentDetails, err
	}

	return studentDetails, nil
}
