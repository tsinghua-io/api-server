package cic

import (
	"github.com/golang/glog"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL     = "http://learn.cic.tsinghua.edu.cn"
	DownloadURL = BaseURL + "/b/resource/downloadFileStream/{file_id}"
)

type CicAdapter struct {
	client   http.Client
	LangCode string
}

var parsedBaseURL, _ = url.Parse(BaseURL)

func New(cookies []*http.Cookie, langCode string) *CicAdapter {
	jar, err := cookiejar.New(nil)
	if err != nil {
		glog.Errorf("Unable to create cookie jar: %s", err)
		return nil
	}
	jar.SetCookies(parsedBaseURL, cookies)

	return &CicAdapter{
		client:   http.Client{Jar: jar},
		LangCode: langCode,
	}
}

// FetchInfo fetches info from url using HTTP GET/POST, and then parse it
// using the given parser. A HTTP status code is returned to indicate the
// result.
func (adapter *CicAdapter) FetchInfo(url string, method string, p parser, info interface{}) (status int) {
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
	if err := p.parse(resp.Body, info, adapter.LangCode); err != nil {
		glog.Errorf("Unable to parse data received from %s: %s", url, err)
		return http.StatusInternalServerError
	}

	glog.Infof("Parsed data from %s (%s)", url, time.Since(t_receive))

	// We are safe.
	return http.StatusOK
}

// Helper functions.

type parser interface {
	parse(reader io.Reader, info interface{}, langCode string) error
}

func parseRegDate(regDate int64) string {
	// Return empty string for 0.
	// Will not work for 1970-01-01T08:00:00+0800.
	if regDate == 0 {
		return ""
	}
	return time.Unix(regDate/1000, 0).Format("2006-01-02T15:04:05+0800")
}

func fileID2DownloadUrl(fileID string) string {
	return strings.Replace(DownloadURL, "{file_id}", fileID, -1)
}
