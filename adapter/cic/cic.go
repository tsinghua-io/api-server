package cic

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	BaseURL         = "http://learn.cic.tsinghua.edu.cn"
	AuthURL         = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
	PersonalInfoURL = BaseURL + "/b/m/getStudentById"
)

// CicAdapter is the adapter for learn.cic.tsinghua.edu.cn
type CicAdapter struct {
	client http.Client
}

func Login(username string, password string) (cookies []*http.Cookie, err error) {
	location, err := getAuth(username, password)
	if err != nil {
		err = fmt.Errorf("Failed to get auth: %s", err)
		return
	}

	// Do not follow 302 redirect.
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		err = fmt.Errorf("Invalid location: %s", err)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		err = fmt.Errorf("Failed to login using auth: %s", err)
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
	req, err := http.NewRequest("POST", AuthURL, strings.NewReader(data))
	if err != nil {
		err = fmt.Errorf("Failed to create the request: %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		err = fmt.Errorf("Request error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Parse redirection.
	location = resp.Header.Get("Location")

	if strings.Contains(location, "status=SUCCESS") {
		err = nil
	} else if strings.Contains(location, "status=BAD_CREDENTIALS") {
		err = fmt.Errorf("Bad credentials.")
	} else {
		err = fmt.Errorf("Unknown new location: %s", location)
	}
	return
}

func New(cookies []*http.Cookie) *CicAdapter {
	adapter := &CicAdapter{}

	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		glog.Errorf("Unable to parse base URL: %s", BaseURL)
		return adapter
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		glog.Errorf("Unable to create cookie jar: %s", err)
		return adapter
	}

	jar.SetCookies(baseURL, cookies)
	adapter.client.Jar = jar

	return adapter
}

func (adapter *CicAdapter) PersonalInfo() (user *resource.User, err error) {
	resp, err := adapter.client.Post(PersonalInfoURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return
	}

	tempUser := resource.User{}

	tempUser.Id, err = j.GetPath("dataSingle", "id").String()
	if err != nil {
		return
	}
	tempUser.Name, err = j.GetPath("dataSingle", "name").String()
	if err != nil {
		return
	}
	tempUser.Department, err = j.GetPath("dataSingle", "majorName").String()
	if err != nil {
		return
	}
	tempUser.Class, err = j.GetPath("dataSingle", "classname").String()
	if err != nil {
		return
	}
	tempUser.Gender, err = j.GetPath("dataSingle", "gender").String()
	if err != nil {
		return
	}
	tempUser.Email, err = j.GetPath("dataSingle", "email").String()
	if err != nil {
		return
	}
	tempUser.Phone, err = j.GetPath("dataSingle", "phone").String()
	if err != nil {
		return
	}

	user = &tempUser
	return
}
