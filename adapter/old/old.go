package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	BaseURL  = "https://learn.tsinghua.edu.cn"
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

func (adapter *OldAdapter) getOldResponse(path string, headers map[string]string) (doc *goquery.Document, err error) {
	url := BaseURL + path
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
	resp, err := adapter.client.Do(req)
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

func (adapter *OldAdapter) PersonalInfo() (*resource.User, int) {
	url := "/MultiLanguage/vspace/vspace_userinfo1.jsp"
	doc, err := adapter.getOldResponse(url, make(map[string]string))

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
			return &resource.User{
				Id:     infos[0],
				Name:   infos[1],
				Type:   infos[14],
				Gender: infos[13],
				Email:  infos[6],
				Phone:  infos[7]}, http.StatusOK
		}
	}
}

func (adapter *OldAdapter) courseIds(typepage int) (courseIdList []string, err error) {
	path := "/MultiLanguage/lesson/student/MyCourse.jsp?typepage=" + strconv.Itoa(typepage)
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		err = fmt.Errorf("Failed to get response from learning web: %s", err)
	} else {
		// parsing the response body
		courseLinkList := doc.Find("#info_1 tr a")
		courseLinkList.Each(func(i int, s *goquery.Selection) {
			var href string
			var hrefUrl *url.URL
			var courseId string
			href, _ = s.Attr("href")
			if hrefUrl, err = url.Parse(href); err != nil {
				return
			}
			if courseId = hrefUrl.Query().Get("course_id"); courseId != "" {
				courseIdList = append(courseIdList, courseId)
			}
		})
	}
	return
}

func (adapter *OldAdapter) courseInfo(courseId string) (course *resource.Course, err error) {
	path := "/MultiLanguage/lesson/student/course_info.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	if err != nil {
		err = fmt.Errorf("Failed to get response from learning web: %s", err)
	} else {
		tds := doc.Find("table#table_box td")

		infos := tds.Map(func(i int, s *goquery.Selection) string {
			firstChild := s.Nodes[0].FirstChild
			if firstChild == nil {
				return ""
			}
			switch s.Nodes[0].FirstChild.Type {
			case html.TextNode:
				info, _ := s.Html()
				return info
			case html.ElementNode:
				info, _ := s.Children().Html()
				return info
			default:
				info, _ := s.Html()
				return info
			}
		})
		if len(infos) < 23 {
			err = fmt.Errorf("Course information parsing error: cannot parse all the informations from %s", infos)
			return
		}
		course = &resource.Course{
			Id:   courseId,
			Name: infos[5],
			Teacher: resource.User{
				Name:  infos[16],
				Email: infos[18],
				Phone: infos[20],
			},
			CourseNumber:   infos[1],
			CourseSequence: infos[3],
			Description:    infos[22],
		}

		if credit, err := strconv.Atoi(strings.TrimSpace(infos[7])); err == nil {
			course.Credit = credit
		}
		if hour, err := strconv.Atoi(strings.TrimSpace(infos[9])); err == nil {
			course.Hour = hour
		}
	}
	return
}

func (adapter *OldAdapter) Attending() (courses []*resource.Course, status int) {
	courseIdList, err := adapter.courseIds(1)
	if err != nil {
		glog.Errorf("Failed to get attending course list: %s", err)
		status = http.StatusBadGateway
		return
	}

	for _, courseId := range courseIdList {
		course, err := adapter.courseInfo(courseId)
		if err != nil {
			glog.Errorf("Failed to get course info of course %s: %s\n", courseId, err)
			status = http.StatusBadGateway
			return
		}
		courses = append(courses, course)
	}
	status = http.StatusOK
	return
}

func (adapter *OldAdapter) Attended() (courses []*resource.Course, status int) {
	courseIdList, err := adapter.courseIds(2)
	if err != nil {
		glog.Errorf("Failed to get attended course list: %s", err)
		status = http.StatusBadGateway
		return
	}

	for _, courseId := range courseIdList {
		course, err := adapter.courseInfo(courseId)
		if err != nil {
			glog.Errorf("Failed to get course info of course %s: %s\n", courseId, err)
			status = http.StatusBadGateway
			return
		}
		courses = append(courses, course)
	}
	status = http.StatusOK
	return
}

func (adapter *OldAdapter) Announcements(courseId string) (courses []*resource.Announcement, status int) {
	return
}

func (adapter *OldAdapter) Files(courseId string) (courses []*resource.File, status int) {
	return
}

func (adapter *OldAdapter) Homeworks(course_id string) (courses []*resource.Homework, status int) {
	return
}
