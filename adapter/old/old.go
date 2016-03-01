package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

const (
	BaseURL  = "https://learn.tsinghua.edu.cn"
	LoginURL = "https://learn.tsinghua.edu.cn/MultiLanguage/lesson/teacher/loginteacher.jsp"
)

// OldAdapter is the adapter for learn.tsinghua.edu.cn
type OldAdapter struct {
	client http.Client
}

func (adapter *OldAdapter) getOldResponse(path string, headers map[string]string) (doc *goquery.Document, err error) {
	url := BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("Failed to create the request: %s", err)
		return
	}

	// Set request headers
	for name, value := range headers {
		req.Header.Add(name, value)
	}
	// Do the request
	resp, err := adapter.client.Do(req)
	if err != nil {
		err = fmt.Errorf("Request Error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Construct goquery.Document
	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		err = fmt.Errorf("Failed to parse response: %s", err)
	}
	// Check if the login time limit exceed
	if doc.Find("div#err").Length() != 0 {
		err = fmt.Errorf("Login time limit exceed.")
	}
	return
}

func New(cookies []*http.Cookie, langCode string) *OldAdapter {
	adapter := &OldAdapter{}

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

func (adapter *OldAdapter) parseFileInfo(path string) (filename string, size int) {
	resp, err := adapter.client.Head(BaseURL + path)
	if err != nil {
		glog.Errorf("Failed to get header information of file %s: %s", path, err)
		return
	}

	// file size
	size, _ = strconv.Atoi(resp.Header.Get("Content-Length"))

	// file name
	disposition := resp.Header.Get("Content-Disposition")
	// decode from gbk
	disposition, _ = simplifiedchinese.GBK.NewDecoder().String(disposition)

	// parse disposition header
	disposition, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		glog.Errorf("Failed to parse header Content-Disposition of file %s", path)
		return
	}
	filename = params["filename"]
	return
}
