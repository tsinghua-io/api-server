package old

import (
	"github.com/golang/glog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"mime"
	"net/http"
	"net/url"
	"strconv"
)

const (
	BaseURL = "https://learn.tsinghua.edu.cn"
)

var parsedBaseURL, _ = url.Parse(BaseURL)

type OldAdapter struct {
	client http.Client
}

// TODO: move to adapter.go
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
