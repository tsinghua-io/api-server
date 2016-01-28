package adapter

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/middleware"
	"github.com/tsinghua-io/api-server/resources"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
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
	return resp.Header.Get("Set-Cookie"), err
}

func (ada OldAdapter) getOldResponse(url string, userSession middleware.UserSession,
	headers map[string]string) (*http.Response, string, error) {
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
		if err != nil {
			return &http.Response{}, "", err
		}
		req.Header.Add("Cookie", cookie)
	}
	// Set request headers
	for name, value := range headers {
		req.Header.Add(name, value)
	}

	resp, err := client.Do(req)
	// TODO: if session expired, do another login.
	// cookie = cookie + resp.Header.Get("Set-Cookie")
	//fmt.Printf("cookie: %s\n", cookie)
	return resp, cookie, err
}

func (ada OldAdapter) GetUserInfo() (CommunicateUnit, error) {
	url := OldUrl + "/MultiLanguage/vspace/vspace_userinfo1.jsp"
	resp, newCookie, err := ada.getOldResponse(url, ada.userSession, make(map[string]string))

	if err != nil {
		return CommunicateUnit{}, err
	} else {
		// parsing the response body
		//body, _ := ioutil.ReadAll(resp.Body)
		//respBody := string(body)
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			return CommunicateUnit{}, err
		}

		docTable := doc.Find("form")
		infos := docTable.Find(".tr_l,.tr_l2").Map(func(i int, valueTR *goquery.Selection) string {
			//fmt.Println(valueTR.Html())
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
