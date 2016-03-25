package old

import (
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	loginURL = "https://learn.tsinghua.edu.cn/MultiLanguage/lesson/teacher/loginteacher.jsp"
)

func Login(name string, pass string) (ada *OldAdapter, status int) {
	ada = &OldAdapter{}
	ada.client.Jar, _ = cookiejar.New(nil)

	form := url.Values{}
	form.Set("userid", name)
	form.Set("userpass", pass)

	resp, err := ada.client.PostForm(loginURL, form)
	if err != nil {
		glog.Errorf("Failed to create request to %s: %s", loginURL, err)
		return nil, http.StatusBadGateway
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if strings.Contains(bodyStr, "用户名或密码错误，登录失败") ||
		strings.Contains(bodyStr, "您没有登陆网络学堂的权限") {
		return nil, http.StatusUnauthorized
	}

	// Login success
	return ada, http.StatusOK
}
