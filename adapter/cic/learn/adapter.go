package learn

import (
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL     = "http://learn.cic.tsinghua.edu.cn"
	downloadURL = BaseURL + "/b/resource/downloadFileStream/{file_id}"
	AuthURL     = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
)

var parsedBaseURL, _ = url.Parse(BaseURL)

type Adapter struct {
	adapter.Adapter
}

func New(userId, password string) (ada *Adapter, status int) {
	ada = new(Adapter)
	ada.AddJar()

	form := url.Values{}
	form.Add("i_user", userId)
	form.Add("i_pass", password)

	resp, err := ada.PostForm(AuthURL, form)
	if err != nil {
		glog.Errorf("Failed to post login form: %s", err)
		return nil, http.StatusBadGateway
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, http.StatusUnauthorized
	}

	return ada, http.StatusOK
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

func DownloadURL(fileId string) string {
	return strings.Replace(downloadURL, "{file_id}", fileId, -1)
}
