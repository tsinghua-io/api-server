package cic

import (
	"github.com/golang/glog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	authURL = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
)

func getAuth(username string, password string) (location string, status int) {
	form := url.Values{}
	form.Add("i_user", username)
	form.Add("i_pass", password)
	data := form.Encode()

	// Do not follow 302 redirect.
	req, err := http.NewRequest("POST", authURL, strings.NewReader(data))
	if err != nil {
		glog.Errorf("Failed to create request to %s:, %s", authURL, err)
		return "", http.StatusInternalServerError
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		glog.Errorf("Failed to get response from %s: %s", authURL, err)
		return "", http.StatusBadGateway
	}
	defer resp.Body.Close()

	// Parse redirection.
	location = resp.Header.Get("Location")

	if strings.Contains(location, "status=SUCCESS") {
		return location, http.StatusOK
	} else if strings.Contains(location, "status=BAD_CREDENTIALS") {
		return location, http.StatusUnauthorized
	} else if location == "" {
		glog.Error("Empty redirection got from %s.", authURL)
		return location, http.StatusBadGateway
	} else {
		glog.Error("Unknown redirection got from %s: %s", authURL, location)
		return location, http.StatusBadGateway
	}
}

func Login(username string, password string) (ada *CicAdapter, status int) {
	location, status := getAuth(username, password)
	if status != http.StatusOK {
		return nil, status
	}

	// Do not follow 302 redirect.
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		glog.Errorf("Failed to create request to %s:, %s", location, err)
		return nil, http.StatusInternalServerError
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		glog.Errorf("Failed to login using auth (%s): %s", location, err)
		return nil, http.StatusBadGateway
	}
	defer resp.Body.Close()

	// Construct adapter.
	jar, err := cookiejar.New(nil)
	if err != nil {
		glog.Errorf("Unable to create cookie jar: %s", err)
		return nil, http.StatusInternalServerError
	}
	jar.SetCookies(parsedBaseURL, resp.Cookies())

	return &CicAdapter{client: http.Client{Jar: jar}}, http.StatusOK
}
