package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"golang.org/x/text/encoding"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"strconv"
)

type Adapter struct {
	*http.Client
}

func New() *Adapter {
	client := new(http.Client)
	client.Jar, _ = cookiejar.New(nil)

	return &Adapter{client}
}

func (ada *Adapter) GetDocument(url string) (doc *goquery.Document, err error) {
	glog.Infof("Getting document from %s", url)

	resp, err := ada.Get(url)
	if err != nil {
		glog.Errorf("Unable to get document from %s: %s", url, err)
		return nil, err
	}

	return goquery.NewDocumentFromResponse(resp)
}

func (ada *Adapter) GetJSON(method, url string, v interface{}) error {
	glog.Infof("Doing %s on JSON from %s", method, url)

	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = ada.Get(url)
	case "POST":
		resp, err = ada.PostForm(url, nil)
	default:
		err = fmt.Errorf("Unknown method: %s", method)
	}

	if err != nil {
		glog.Errorf("Unable to get JSON from %s: %s", url, err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(v); err != nil {
		glog.Errorf("Unable to decode JSON: %s", err)
		return err
	}

	return nil
}

func (ada *Adapter) FileInfo(url string, encoding encoding.Encoding) (filename string, size int, status int) {
	status = http.StatusBadGateway

	resp, err := ada.Head(url)
	if err != nil {
		glog.Errorf("Failed to HEAD file at %s: %s", url, err)
		return
	}

	status = http.StatusInternalServerError

	// Filename.
	disposition := resp.Header.Get("Content-Disposition")
	disposition, err = encoding.NewDecoder().String(disposition)
	if err != nil {
		glog.Errorf("Failed to decode Content-Disposition: %s", err)
		return
	}

	// Parse disposition header.
	_, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		glog.Errorf("Failed to parse header Content-Disposition of file at %s: %s", url, err)
		return
	}
	filename = params["filename"]

	// File size.
	if size != nil {
		sizeStr := resp.Header.Get("Content-Length")
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			glog.Errorf("Failed to convert Content-Length (%s) to int: %s", sizeStr, err)
			return
		}
	}

	status = http.StatusOK
	return
}

func MergeStatus(statuses ...int) (status int) {
	for _, s := range statuses {
		if s != 0 && s != http.StatusOK {
			return s
		}
	}
	return http.StatusOK
}
