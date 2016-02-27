package cic

import (
	"github.com/golang/glog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type CicAdapter struct {
	client http.Client
}

const (
	BaseURL = "http://learn.cic.tsinghua.edu.cn"
)

var parsedBaseURL, _ = url.Parse(BaseURL)

func New(cookies []*http.Cookie) *CicAdapter {
	jar, err := cookiejar.New(nil)
	if err != nil {
		glog.Errorf("Unable to create cookie jar: %s", err)
		return nil
	}
	jar.SetCookies(parsedBaseURL, cookies)

	return &CicAdapter{
		client: http.Client{Jar: jar},
	}
}

// FetchInfo fetches info from url using HTTP GET/POST, and then parse it
// using the given parser. A HTTP status code is returned to indicate the
// result.
func (adapter *CicAdapter) FetchInfo(url string, method string, langCode string, p parser, info interface{}) (status int) {
	// Fetch data from url.
	glog.Infof("Fetching data from %s", url)

	var resp *http.Response
	var err error
	t_send := time.Now()

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

	t_receive := time.Now()
	glog.Infof("Fetched data from %s (%s)", url, t_receive.Sub(t_send))

	// Parse the data.
	if err := p.parse(resp.Body, info, langCode); err != nil {
		glog.Errorf("Unable to parse data received from %s: %s", url, err)
		return http.StatusBadGateway
	}

	glog.Infof("Parsed data from %s (%s)", url, time.Since(t_receive))

	// We are safe.
	return http.StatusOK
}
