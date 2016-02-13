package old

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"fmt"
)

const (
	BaseURL = "https://learn.tsinghua.edu.cn"
	LoginURL = "https://learn.tsinghua.edu.cn/MultiLanguage/lesson/teacher/loginteacher.jsp"
)

// OldAdapter is the adapter for learn.tsinghua.edu.cn
type OldAdapter struct {
	client http.Client
}


func Login(name string, pass string) (cookies []*http.Cookie, err error) {
	form := url.Values{}
	form.Add("userid", name)
	form.Add("userpass", pass)
	resp, err := http.PostForm(LoginURL, form)
	if err != nil {
		err = fmt.Errorf("Failed to create the request: %s", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(body), "用户名或密码错误，登录失败") ||
		strings.Contains(string(body), "您没有登陆网络学堂的权限") {
		cookies = []*http.Cookie{}
		err = fmt.Errorf("Bad credentials.")
	} else {
		// Login success
		cookies = resp.Cookies()
		err = nil
	}
	return
}

func (ada *OldAdapter) getOldResponse(url string, headers map[string]string) (doc *goquery.Document, err error) {
	url = BaseURL + url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("Failed to create the request: %s", err)
		return
	}

	// Set request headers
	for name, value := range headers {
		req.Header.Add(name, value)
	}
	// Do the request
	resp, err := ada.client.Do(req)
	if err != nil {
		err = fmt.Errorf("Request Error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Construct goquery.Document
	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		err = fmt.Errorf("Failed to parse response: %s", err)
	}
	return
}

func (ada *OldAdapter) PersonalInfo() (interface{}, int) {
	url := "/MultiLanguage/vspace/vspace_userinfo1.jsp"
	doc, err := ada.getOldResponse(url, make(map[string]string))

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
		return nil, http.StatusBadGateway
	} else {
		// parsing the response body
		docTable := doc.Find("form")
		infos := docTable.Find(".tr_l,.tr_l2").Map(func(i int, valueTR *goquery.Selection) string {
			switch valueTR.Nodes[0].FirstChild.Type {
			case html.TextNode:
				info, _ := valueTR.Html()
				return info
			case html.ElementNode:
				info, _ := valueTR.Children().Attr("value")
				return info
			default:
				info, _ := valueTR.Html()
				return info
			}

		})
		if len(infos) < 15 {
			glog.Errorf("User information parsing error: cannot parse all the informations from %s", infos)
			return nil, http.StatusBadGateway
		} else {
			return resource.User{
				Id:        infos[0],
				Name:      infos[1],
				Type: infos[14],
				Gender: infos[13],
				Email:  infos[6],
				Phone:  infos[7]}, http.StatusOK
		}
	}
}

func New(cookies []*http.Cookie) *OldAdapter {
	adapter := &OldAdapter{}

	baseURL, err := url.Parse(BaseURL)
	if err != nil {
		glog.Errorf("Unable to parse base URL: %s", BaseURL)
		return adapter
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		glog.Errorf("Unable to create cookie jar: %s", err)
		return adapter
	}

	jar.SetCookies(baseURL, cookies)
	adapter.client.Jar = jar
	return adapter
}
