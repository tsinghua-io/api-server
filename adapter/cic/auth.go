package cic

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	AuthURL = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
)

func getAuth(username string, password string) (location string, err error) {
	form := url.Values{}
	form.Add("i_user", username)
	form.Add("i_pass", password)
	data := form.Encode()

	// Do not follow 302 redirect.
	req, err := http.NewRequest("POST", AuthURL, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("Failed to create the request: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", fmt.Errorf("Request error: %s", err)
	}
	defer resp.Body.Close()

	// Parse redirection.
	location = resp.Header.Get("Location")

	if strings.Contains(location, "status=SUCCESS") {
		err = nil
	} else if strings.Contains(location, "status=BAD_CREDENTIALS") {
		err = fmt.Errorf("Bad credentials.")
	} else if location == "" {
		err = fmt.Errorf("No new location provided.")
	} else {
		err = fmt.Errorf("Unknown new location: %s", location)
	}
	return
}

func Login(username string, password string) (cookies []*http.Cookie, err error) {
	location, err := getAuth(username, password)
	if err != nil {
		return nil, fmt.Errorf("Failed to get auth: %s", err)
	}

	// Do not follow 302 redirect.
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return nil, fmt.Errorf("Invalid location: %s", err)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to login using auth: %s", err)
	}
	defer resp.Body.Close()

	cookies = resp.Cookies()
	return
}
