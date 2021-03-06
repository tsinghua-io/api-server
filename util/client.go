package util

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"golang.org/x/text/encoding"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

type Client struct{ http.Client }

func (client *Client) WithJar() {
	client.Jar, _ = cookiejar.New(nil)
}

func (client *Client) GetDocument(url string) (doc *goquery.Document, errMsg error) {
	glog.Infof("Getting document from %s", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to get document from %s: %s", url, err)
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse document received from %s: %s", url, err)
	}

	return
}

func (client *Client) PostFormJSON(url string, data url.Values, v interface{}) error {
	glog.Infof("Posting %v to %s", data, url)

	resp, err := client.PostForm(url, data)
	if err != nil {
		return fmt.Errorf("Failed to post form to %s: %s", url, err)
	}
	return parseJSON(resp, v)
}

func (client *Client) GetJSON(url string, v interface{}) error {
	glog.Infof("Getting JSON from %s", url)

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to get from %s: %s", url, err)
	}
	return parseJSON(resp, v)
}

func parseJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("Failed to parse JSON received from %s: %s", resp.Request.URL.String(), err)
	}

	return nil
}

func (client *Client) FileInfo(url string, encoding encoding.Encoding) (filename string, size int, status int, errMsg error) {
	status = http.StatusOK

	resp, err := client.Head(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = fmt.Errorf("Failed to HEAD from %s: %s", url, err)
		return
	}

	// Filename.
	disposition := resp.Header.Get("Content-Disposition")
	disposition, err = encoding.NewDecoder().String(disposition)
	if err != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to decode Content-Disposition: %s", err)
		return
	}

	// Parse disposition header.
	_, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse header Content-Disposition of file at %s: %s", url, err)
		return
	}
	filename = params["filename"]

	// File size.
	sizeStr := resp.Header.Get("Content-Length")
	size, err = strconv.Atoi(sizeStr)
	if err != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to convert Content-Length (%s) to int: %s", sizeStr, err)
		return
	}

	return
}
