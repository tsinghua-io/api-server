package cic

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL     = "http://learn.cic.tsinghua.edu.cn"
	downloadURL = BaseURL + "/b/resource/downloadFileStream/{file_id}"
)

var parsedBaseURL, _ = url.Parse(BaseURL)

type CicAdapter struct {
	client http.Client
}

// Helper functions.

func parseRegDate(regDate int64) string {
	// Return empty string for 0.
	// Will not work for 1970-01-01T08:00:00+0800.
	if regDate == 0 {
		return ""
	}
	return time.Unix(regDate/1000, 0).Format("2006-01-02T15:04:05+0800")
}

func fileID2DownloadUrl(fileID string) string {
	return strings.Replace(downloadURL, "{file_id}", fileID, -1)
}
