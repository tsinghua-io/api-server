package adapter

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/middleware"
	"github.com/tsinghua-io/api-server/resources"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

// OldAdapter is the adapter for learn.tsinghua.edu.cn
type OldAdapter struct {
	userSession middleware.UserSession
}

func NewOldAdapter(userSession middleware.UserSession) Adapter {
	return OldAdapter{userSession}
}

func loginOld(name string, pass string) (string, error) {
	form := url.Values{}
	form.Add("userid", name)
	form.Add("userpass", pass)
	resp, err := http.PostForm("https://learn.tsinghua.edu.cn/MultiLanguage/lesson/teacher/loginteacher.jsp", form)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	match, _ := regexp.MatchString("用户名或密码错误，登录失败", string(body))
	if match {
		// Login failed.
		return "", err
	} else {
		// Login success.
		return resp.Header.Get("Set-Cookie"), err
	}
}

func (ada OldAdapter) getOldResponse(url string, userSession middleware.UserSession,
	headers map[string]string) (*goquery.Document, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	var cookie string
	if userSession.Session != "" {
		cookie = userSession.Session
		req.Header.Add("Cookie", cookie)
	} else {
		// login and get session
		var err error
		cookie, err = loginOld(userSession.LoginName, userSession.LoginPass)
		if err != nil || cookie == "" {
			return &goquery.Document{}, "", err
		}
		req.Header.Add("Cookie", cookie)
	}

	// Set request headers
	for name, value := range headers {
		req.Header.Add(name, value)
	}
	// Do the request
	resp, err := client.Do(req)
	defer resp.Body.Close()
	// Construct goquery.Document
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return &goquery.Document{}, cookie, err
	}
	content, _ := doc.Html()
	match, _ := regexp.MatchString("请您登陆后再使用", content)
	if match {
		// Maybe session expire. Try relogin once
		var err error
		cookie, err = loginOld(userSession.LoginName, userSession.LoginPass)
		if err != nil || cookie == "" {
			// Error or login failed.
			return &goquery.Document{}, "", err
		}
		req.Header.Set("Cookie", cookie)

		resp, err := client.Do(req)
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			return &goquery.Document{}, cookie, err
		}
		return doc, cookie, err
	} else {
		return doc, cookie, err
	}
}

func (ada OldAdapter) GetUserInfo() (CommunicateUnit, error) {
	url := OldUrl + "/MultiLanguage/vspace/vspace_userinfo1.jsp"
	doc, newCookie, err := ada.getOldResponse(url, ada.userSession, make(map[string]string))

	if err != nil {
		return CommunicateUnit{}, err
	} else if newCookie == "" {
		// Login failed
		return CommunicateUnit{nil, http.StatusUnauthorized, newCookie}, nil
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
		if len(infos) <= 0 {
			return CommunicateUnit{}, nil
		} else {
			return CommunicateUnit{
				resources.User{
					SUser: resources.SUser{
						Id:        infos[0],
						Name:      infos[1],
						User_type: infos[14]},
					Gender: infos[13],
					Email:  infos[6],
					Phone:  infos[7]},
				http.StatusOK,
				newCookie}, nil
		}
	}
}
