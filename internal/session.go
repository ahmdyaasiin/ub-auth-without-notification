package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetSession() (*Session, error) {
	session := new(Session)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://brone.ub.ac.id/my/", nil)
	if err != nil {
		return session, fmt.Errorf("failed to create new request: %v", err)
	}

	for k, v := range GetHeaders() {
		req.Header.Set(k, v)
	}

	req.Header.Set("referer", "https://brone.ub.ac.id/")

	resp, err := client.Do(req)
	if err != nil {
		return session, fmt.Errorf("failed to perform GET request: %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return session, fmt.Errorf("failed to read response body: %v", err)
	}

	r := fmt.Sprintf("%s", respBody)
	h := fmt.Sprintf("%s", resp.Header["Set-Cookie"])

	authSessionID, err := GetSubstringBetween("AUTH_SESSION_ID=", ";", h)
	if err != nil {
		return session, fmt.Errorf("failed to get auth-session-id")
	}

	authSessionIDLegacy, err := GetSubstringBetween("AUTH_SESSION_ID_LEGACY=", ";", h)
	if err != nil {
		return session, fmt.Errorf("failed to get auth-session-id-legacy")
	}

	kcRestart, err := GetSubstringBetween("KC_RESTART=", ";", h)
	if err != nil {
		return session, fmt.Errorf("failed to get kc-start")
	}

	fullURL, err := GetSubstringBetween(`action="`, `" `, r)
	if err != nil {
		return session, fmt.Errorf("failed to get full-url")
	}

	sessionCode, err := GetSubstringBetween("session_code=", "&amp", fullURL)
	if err != nil {
		return session, fmt.Errorf("failed to get session-code")
	}

	execution, err := GetSubstringBetween("execution=", "&amp", fullURL)
	if err != nil {
		return session, fmt.Errorf("failed to get execution")
	}

	tabIDSlice := strings.Split(fullURL, "tab_id=")
	if len(tabIDSlice) < 2 {
		return session, fmt.Errorf("failed to get tab-id")
	}

	session.Client = client
	session.AuthSessionID = authSessionID
	session.AuthSessionIDLegacy = authSessionIDLegacy
	session.KCRestart = kcRestart
	session.SessionCode = sessionCode
	session.Execution = execution
	session.TabID = tabIDSlice[1]

	return session, nil
}
