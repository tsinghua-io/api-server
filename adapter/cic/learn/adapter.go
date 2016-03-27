package learn

import (
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"net/http"
	"net/url"
)

const (
	BaseURL = "http://learn.cic.tsinghua.edu.cn"
	AuthURL = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
)

type Adapter struct {
	*adapter.Adapter
}

func New(userId, password string) (ada *Adapter, status int) {
	ada = &Adapter{adapter.New()}

	form := url.Values{}
	form.Add("i_user", userId)
	form.Add("i_pass", password)

	resp, err := ada.PostForm(AuthURL, form)
	if err != nil {
		glog.Errorf("Failed to post login form to %s: %s", AuthURL, err)
		return nil, http.StatusBadGateway
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, http.StatusUnauthorized
	}

	return ada, http.StatusOK
}
