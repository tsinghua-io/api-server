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
		infos := docTable.Find(".tr_l,.tr_l2").Map(func(i int, valueTR *goquery.Selection) (info string) {
			switch valueTR.Nodes[0].FirstChild.Type {
			case html.TextNode:
				info, _ = valueTR.Html()
			case html.ElementNode:
				info, _ = valueTR.Children().Attr("value")
			default:
				info, _ = valueTR.Html()
			}
			return
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

func (adapter *OldAdapter) Announcements(courseId string) (announcements []*resource.Announcement, status int) {
	path := "/MultiLanguage/public/bbs/getnoteid_student.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	status = http.StatusBadGateway

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		trs := doc.Find("tr.tr1, tr.tr2")
		trs.Each(func(i int, s *goquery.Selection) {
			infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
				info, _ = tdSelection.Html()
				return
			})

			hrefSelection := s.Find("td a")

			var href string
			var hrefUrl *url.URL
			var err error
			href, _ = hrefSelection.Attr("href")

			if hrefUrl, err = url.Parse(href); err != nil {
				return
			}

			if announcementId := hrefUrl.Query().Get("id"); announcementId != "" {
				body := adapter.announcementBody(href)

				var important bool
				var title string
				switch hrefSelection.Nodes[0].FirstChild.Type {
				case html.TextNode:
					important = false
					title, _ = hrefSelection.Html()
				default:
					important = true
					title, _ = hrefSelection.Children().Html()
				}

				announcements = append(announcements, &resource.Announcement{
					Id:       announcementId,
					CourseId: courseId,
					Title:    title,
					Owner: resource.User{
						Name: infos[2],
					},
					CreatedAt: infos[3],
					Important: important,
					Body:      body,
				})
			}
		})
		if len(announcements) > 0 {
			status = http.StatusOK
		}
	}
	return
}

func (adapter *OldAdapter) Files(courseId string) (files []*resource.File, status int) {
	path := "/MultiLanguage/lesson/student/download.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	status = http.StatusBadGateway

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		// Find all categories
		categories := doc.Find("td.textTD").Map(func(i int, s *goquery.Selection) (info string) {
			info, _ = s.Html()
			return
		})

		categoryDivs := doc.Find("div.layerbox")
		categoryDivs.Each(func(i int, div *goquery.Selection) {
			category := categories[i]
			trs := div.Find("#table_box tr~tr")
			trs.Each(func(i int, s *goquery.Selection) {
				infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
					info, _ = tdSelection.Html()
					return
				})

				hrefSelection := s.Find("td a")

				var href string
				var hrefUrl *url.URL
				var err error
				if href, _ = hrefSelection.Attr("href"); href == "" {
					return
				}
				if hrefUrl, err = url.Parse(href); err != nil {
					return
				}

				if fileId := hrefUrl.Query().Get("file_id"); fileId != "" {
					title, _ := hrefSelection.Html()
					title = strings.TrimSpace(title)
					file := &resource.File{
						Id:          fileId,
						CourseId:    courseId,
						Category:    []string{category},
						Title:       title,
						Description: infos[2],
						DownloadUrl: href,
						CreatedAt:   infos[4],
					}

					file.Filename, file.Size = adapter.parseFileInfo(href)
					files = append(files, file)
				}
			})
		})
		if len(files) > 0 {
			status = http.StatusOK
		}
	}
	return
}

func (adapter *OldAdapter) Homeworks(courseId string) (homeworks []*resource.Homework, status int) {
	path := "/MultiLanguage/lesson/student/hom_wk_brw.jsp?course_id=" + courseId
	doc, err := adapter.getOldResponse(path, make(map[string]string))

	status = http.StatusBadGateway

	if err != nil {
		glog.Errorf("Failed to get response from learning web: %s", err)
	} else {
		trs := doc.Find("tr.tr1, tr.tr2")
		trs.Each(func(i int, s *goquery.Selection) {
			infos := s.Find("td").Map(func(i int, tdSelection *goquery.Selection) (info string) {
				info, _ = tdSelection.Html()
				return
			})

			hrefSelection := s.Find("td a")

			var href string
			var hrefUrl *url.URL
			var err error
			href, _ = hrefSelection.Attr("href")

			if hrefUrl, err = url.Parse(href); err != nil {
				return
			}

			if homeworkId := hrefUrl.Query().Get("id"); homeworkId != "" {
				title, _ := hrefSelection.Html()

				homework := &resource.Homework{
					Id:        homeworkId,
					CourseId:  courseId,
					Title:     title,
					CreatedAt: infos[1],
					BeginAt:   infos[1],
					DueAt:     infos[2],
				}
				homework.Body, homework.Attachment = adapter.parseHomeworkInfo(href)
				homeworks = append(homeworks, homework)
			}
		})
		if len(homeworks) > 0 {
			status = http.StatusOK
		}
	}
	return
}
