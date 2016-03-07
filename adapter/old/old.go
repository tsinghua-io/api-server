package old

import (
	"fmt"
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
func parseFileInfo(client http.Client, path string) (filename string, size int, err error) {
	resp, err := client.Head(path)
	if err != nil {
		return "", 0, fmt.Errorf("Failed to get header information of file %s: %s", path, err)
	}

	// file size
	sizeStr := resp.Header.Get("Content-Length")
	size, err = strconv.Atoi(sizeStr)
	if err != nil {
		return "", 0, err
	}

	// file name
	disposition := resp.Header.Get("Content-Disposition")
	// decode from gbk
	disposition, err = simplifiedchinese.GBK.NewDecoder().String(disposition)
	if err != nil {
		return "", 0, err
	}

	// parse disposition header
	disposition, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		return "", 0, fmt.Errorf("Failed to parse header Content-Disposition of file %s", path)
	}
	filename = params["filename"]

	return filename, size, nil
}
