package learn

import (
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	BaseURL = "https://learn.tsinghua.edu.cn"
	AuthURL = "https://learn.tsinghua.edu.cn/MultiLanguage/lesson/teacher/loginteacher.jsp"
)

type Adapter struct {
	adapter.Adapter
}

func New(userId, password string) (ada *Adapter, status int) {
	ada = new(Adapter)
	ada.AddJar()

	form := url.Values{}
	form.Add("userid", userId)
	form.Add("userpass", password)

	resp, err := ada.PostForm(AuthURL, form)
	if err != nil {
		glog.Errorf("Failed to post login form to %s: %s", AuthURL, err)
		return nil, http.StatusBadGateway
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if strings.Contains(bodyStr, "用户名或密码错误，登录失败") ||
		strings.Contains(bodyStr, "您没有登陆网络学堂的权限") {
		return nil, http.StatusUnauthorized
	}

	return ada, http.StatusOK
}
