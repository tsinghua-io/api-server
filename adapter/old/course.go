package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

		infos := tds.Map(func(i int, s *goquery.Selection) (info string) {
			firstChild := s.Nodes[0].FirstChild
			if firstChild != nil {
				switch s.Nodes[0].FirstChild.Type {
				case html.TextNode:
					info, _ = s.Html()
				case html.ElementNode:
					info, _ = s.Children().Html()
				default:
					info, _ = s.Html()
				}
			}
			return
		})
		if len(infos) < 23 {
			err = fmt.Errorf("Course information parsing error: cannot parse all the informations from %s", infos)
			return
		}
		course = &resource.Course{
			Id:   courseId,
			Name: infos[5],
			Teachers: []*resource.User{&resource.User{
				Name:  infos[16],
				Email: infos[18],
				Phone: infos[20],
			}},
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
