package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"net/url"
)

const (
	BaseURL = "http://learn.cic.tsinghua.edu.cn"
	AuthURL = "https://id.tsinghua.edu.cn/do/off/ui/auth/login/post/fa8077873a7a80b1cd6b185d5a796617/0?/j_spring_security_thauth_roaming_entry"
)

type Adapter struct{ util.Client }

func New(userId, password string) (ada *Adapter, status int, errMsg error) {
	ada = new(Adapter)
	ada.WithJar()
	status = http.StatusOK

	form := url.Values{}
	form.Add("i_user", userId)
	form.Add("i_pass", password)

	resp, err := ada.PostForm(AuthURL, form)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("Failed to post login form to %s: %s", AuthURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg = fmt.Errorf("Server replied %s.", http.StatusText(resp.StatusCode))
	} else if url := resp.Request.URL; url.Host != "learn.cic.tsinghua.edu.cn" {
		errMsg = fmt.Errorf("Stopped at incorrect URL: %s.", url.String())
	}
	if errMsg != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("Failed to login to %s: %s", AuthURL, errMsg)
	}

	return
}

func HandlerFunc(f func(http.ResponseWriter, *http.Request, *Adapter)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		userId, password, _ := req.BasicAuth()
		if ada, status, err := New(userId, password); err != nil {
			util.Error(rw, err.Error(), status)
		} else {
			f(rw, req, ada)
		}
	})
}
