package old

import (
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Login(name string, pass string) (cookies []*http.Cookie, status int) {
	form := url.Values{}
	form.Set("userid", name)
	form.Set("userpass", pass)

	resp, err := http.PostForm(LoginURL, form)
	if err != nil {
		glog.Errorf("Failed to create the request: %s", err)
		status = http.StatusBadGateway
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "用户名或密码错误，登录失败") ||
		strings.Contains(string(body), "您没有登陆网络学堂的权限") {
		status = http.StatusUnauthorized
	} else {
		// Login success
		cookies = resp.Cookies()
		status = http.StatusOK
	}
	return
}
