package cic

import (
	"fmt"
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
	// AttendedURL     = BaseURL + "/b/myCourse/courseList/loadCourse4Student/-1"
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
	} else if location == "" {
		err = fmt.Errorf("No new location provided.")
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

func (adapter *CicAdapter) FetchInfo(url string, method string, p parser, info interface{}) (status int) {
	// Fetch data from url.
	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = adapter.client.Get(url)
	case "POST":
		resp, err = adapter.client.Post(url, "application/x-www-form-urlencoded", nil)
	default:
		glog.Errorf("Unknown method to fetch info: %s", method)
		return http.StatusInternalServerError
	}
	if err != nil {
		glog.Errorf("Unable to fetch info from %s: %s", url, err)
		return http.StatusBadGateway
	}
	defer resp.Body.Close()

	if err := p.parse(resp.Body, info); err != nil {
		glog.Errorf("Unable to parse data received from %s: %s", url, err)
		return http.StatusBadGateway
	}

	// We are safe.
	return http.StatusOK
}

func (adapter *CicAdapter) PersonalInfo() (user *resource.User, status int) {
	// Fill data into User.
	user = &resource.User{}
	status = adapter.FetchInfo(PersonalInfoURL, "POST", &personalInfoParser{}, user)
	return
}

func (adapter *CicAdapter) Attending() (courses []*resource.Course, status int) {
	return
}

func (adapter *CicAdapter) Attended() (courses []*resource.Course, status int) {
	return
}

func (adapter *CicAdapter) Announcements(course_id string) (courses []*resource.Announcement, status int) {
	return
}

func (adapter *CicAdapter) Files(course_id string) (courses []*resource.File, status int) {
	return
}

func (adapter *CicAdapter) Homeworks(course_id string) (courses []*resource.Homework, status int) {
	return
}
