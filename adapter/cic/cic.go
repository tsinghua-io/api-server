package cic

import (
	"errors"
	"strings"
	//"github.com/tsinghua-io/api-server/resources"
	// "github.com/golang/glog"
	// "io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

const (
	BaseURL  = "https://learn.cic.tsinghua.edu.cn"
	LoginURL = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
)

// CicAdapter is the adapter for learn.cic.tsinghua.edu.cn
type CicAdapter struct {
	cookies []*http.Cookie
}

func Login(username string, password string) (cookies []*http.Cookie, err error) {
	location, err := getAuth(username, password)
	if err != nil {
		err = errors.New("Failed to get auth: " + err.Error())
		return
	}

	// Do not follow 302 redirect.
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		err = errors.New("Invalid location: " + err.Error())
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		err = errors.New("Failed to login using auth: " + err.Error())
		return
	}
	defer resp.Body.Close()

	cookies = resp.Cookies()
	return
}

func getAuth(username string, password string) (location string, err error) {
	form := url.Values{}
	form.Add("i_user", username)
	form.Add("i_pass", password)
	data := form.Encode()

	// Do not follow 302 redirect.
	req, err := http.NewRequest("POST", LoginURL, strings.NewReader(data))
	if err != nil {
		err = errors.New("Failed to create the request: " +
			err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		err = errors.New("Request error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// Parse redirection.
	location = resp.Header.Get("Location")

	if strings.Contains(location, "status=SUCCESS") {
		err = nil
	} else if strings.Contains(location, "status=BAD_CREDENTIALS") {
		err = errors.New("Bad credentials.")
	} else {
		err = errors.New("Unknown new location: " + location)
	}
	return
}

func New(jar cookiejar.Jar) CicAdapter {
	baseURL, err := url.Parse(BaseURL)
	if err != nil {

	}
	return CicAdapter{jar.Cookies(baseURL)}
}
