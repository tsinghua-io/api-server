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

type UserParser struct {
	Id         string
	Name       string
	Type       string
	Department string
	Class      string
	Gender     string
	Email      string
	Phone      string
}

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

// readUser reads a User from a json, using the given paths.
func readUser(j *simplejson.Json, parser *UserParser) (user *resource.User, err error) {

	tempUser := &resource.User{}
	if parser.Id != "" {
		if tempUser.Id, err = j.GetPath(parser.Id).String(); err != nil {
			return
		}
	}
	if parser.Name != "" {
		if tempUser.Name, err = j.GetPath(parser.Name).String(); err != nil {
			return
		}
	}
	if parser.Department != "" {
		if tempUser.Department, err = j.GetPath(parser.Department).String(); err != nil {
			return
		}
	}
	if parser.Class != "" {
		if tempUser.Class, err = j.GetPath(parser.Class).String(); err != nil {
			return
		}
	}
	if parser.Gender != "" {
		if tempUser.Gender, err = j.GetPath(parser.Gender).String(); err != nil {
			return
		}
	}
	if parser.Email != "" {
		if tempUser.Email, err = j.GetPath(parser.Email).String(); err != nil {
			return
		}
	}
	if parser.Phone != "" {
		if tempUser.Phone, err = j.GetPath(parser.Phone).String(); err != nil {
			return
		}
	}

	// Safe, we are done.
	user = tempUser
	return
}

func (adapter *CicAdapter) PersonalInfo() (user *resource.User, status int) {
	resp, err := adapter.client.Post(PersonalInfoURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		glog.Errorf("Unable to fetch personal info: %s", err)
		status = http.StatusBadGateway
		return
	}
	defer resp.Body.Close()

	// Parse into JSON.
	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		glog.Errorf("Unable to parse the response to JSON: %s", err)
		status = http.StatusBadGateway
		return
	}

	// Fill data into User.
	parser := &UserParser{
		Id:         "id",
		Name:       "name",
		Type:       "",
		Department: "majorName",
		Class:      "classname",
		Gender:     "gender",
		Email:      "email",
		Phone:      "phone",
	}
	if user, err = readUser(j.Get("dataSingle"), parser); err != nil {
		// Failed.
		glog.Errorf("Unable to parse all the fields: %s", err)
		status = http.StatusBadGateway
		return
	}

	status = http.StatusOK
	return
}
