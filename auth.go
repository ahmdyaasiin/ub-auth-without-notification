package ub_auth

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type StudentDetails struct {
	NIM          string `json:"nim"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Fakultas     string `json:"fakultas"`
	ProgramStudi string `json:"program_studi"`
}

type ResponseDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *ResponseDetails) Error() string {
	return r.Message
}

func AuthUB(username, password string) (*StudentDetails, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://brone.ub.ac.id/my/", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://brone.ub.ac.id/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := fmt.Sprintf("%s", respBody)
	h := fmt.Sprintf("%s", resp.Header["Set-Cookie"])

	authSessionID := getBetween("AUTH_SESSION_ID=", ";", h)
	authSessionIDLegacy := getBetween("AUTH_SESSION_ID_LEGACY=", ";", h)
	kcRestart := getBetween("KC_RESTART=", ";", h)

	fullURL := getBetween(`action="`, `" `, r)
	sessionCode := getBetween("session_code=", "&amp", fullURL)
	execution := getBetween("execution=", "&amp", fullURL)
	tabID := strings.Split(fullURL, "tab_id=")[1]

	data := strings.NewReader(fmt.Sprintf("username=%s&password=%s&credentialId=%s", username, password, ""))
	req, err = http.NewRequest("POST", fmt.Sprintf("https://iam.ub.ac.id/auth/realms/ub/login-actions/authenticate?session_code=%s&execution=%s&client_id=brone.ub.ac.id&tab_id=%s", sessionCode, execution, tabID), data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "null")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("cookie", fmt.Sprintf("AUTH_SESSION_ID=%s; AUTH_SESSION_ID_LEGACY=%s; KC_RESTART=%s", authSessionID, authSessionIDLegacy, kcRestart))

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	r = fmt.Sprintf("%s", respBody)
	if !strings.Contains(r, "SAMLResponse") {
		if strings.Contains(r, "Invalid username or password.") {

			return &StudentDetails{}, &ResponseDetails{401, "Invalid username or password"}
		} else {

			return &StudentDetails{}, &ResponseDetails{500, "Unexpected error"}
		}
	}

	samlResponse := getBetween(`name="SAMLResponse" value="`, `"/>`, r)

	decodedSamlResponseByte, err := base64.StdEncoding.DecodeString(samlResponse)
	if err != nil {
		log.Fatal(err)
	}

	decodedSamlResponse := string(decodedSamlResponseByte)

	// parse student details
	nim := getBetween(`<saml:Attribute FriendlyName="nim" Name="nim" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">`, "</saml", decodedSamlResponse)
	fullName := getBetween(`<saml:Attribute FriendlyName="fullName" Name="fullName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">`, "</saml", decodedSamlResponse)
	email := getBetween(`<saml:Attribute FriendlyName="email" Name="email" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">`, "</saml", decodedSamlResponse)
	fakultas := fmt.Sprintf("Fakultas %s", getBetween(`<saml:Attribute FriendlyName="fakultas" Name="fakultas" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">`, "</saml", decodedSamlResponse))
	programStudi := getBetween(`<saml:Attribute FriendlyName="prodi" Name="prodi" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:basic"><saml:AttributeValue xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="xs:string">`, "</saml", decodedSamlResponse)

	return &StudentDetails{nim, fullName, email, fakultas, programStudi}, nil
}

func getBetween(start, end, text string) string {
	r := strings.Split(text, start)[1]

	return strings.Split(r, end)[0]
}
