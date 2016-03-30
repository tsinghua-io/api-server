package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	BaseURL = "https://learn.tsinghua.edu.cn"
	AuthURL = "https://learn.tsinghua.edu.cn/MultiLanguage/lesson/teacher/loginteacher.jsp"
)

type Adapter struct{ util.Client }

func New(userId, password string) (ada *Adapter, status int, errMsg error) {
	ada = new(Adapter)
	ada.WithJar()
	status = http.StatusOK

	form := url.Values{}
	form.Add("userid", userId)
	form.Add("userpass", password)

	resp, err := ada.PostForm(AuthURL, form)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("Failed to post login form to %s: %s", AuthURL, err)
	}
	defer resp.Body.Close()

	location := "loginteacher_action.jsp"
	if body, _ := ioutil.ReadAll(resp.Body); !strings.Contains(string(body), location) {
		return nil, http.StatusUnauthorized, fmt.Errorf("Failed to login to %s: No \"%s\" found in response", AuthURL, location)
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
