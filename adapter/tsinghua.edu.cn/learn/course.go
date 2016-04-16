package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"net/url"
	"path"
	"regexp"
)

func AttendedURL(page int) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/MyCourse.jsp?typepage=%d&language=cn", BaseURL, page)
}

func ParseCourseURL(rawurl string) (id string) {
	if parsed, err := url.Parse(rawurl); err == nil {
		if parsed.Host == "" {
			id = parsed.Query().Get("course_id")
		} else if parsed.Host == "learn.cic.tsinghua.edu.cn" {
			id = path.Base(parsed.Path)
		}
	}
	return
}

func ParseCourseName(rawName string) (name, seq, sem string) {
	regex, _ := regexp.Compile("\\s*(.*)\\((\\S+)\\)\\((\\d{4}-\\d{4})(秋|春|夏)季学期\\)")
	if match := regex.FindStringSubmatch(rawName); len(match) > 4 {
		name = match[1]
		seq = match[2]

		switch match[4] {
		case "秋":
			sem = match[3] + "-1"
		case "春":
			sem = match[3] + "-2"
		case "夏":
			sem = match[3] + "-3"
		}
	}
	return
}

func (ada *Adapter) AttendedList(page int) (courses []*model.Course, status int, errMsg error) {
	status = http.StatusOK

	url := AttendedURL(page)
	doc, err := ada.GetDocument(url)
	if err != nil {
		status = http.StatusBadGateway
		errMsg = err
		return
	}

	links := doc.Find("#info_1 a")
	courses = make([]*model.Course, links.Size())

	links.EachWithBreak(func(i int, s *goquery.Selection) bool {
		if id := ParseCourseURL(s.AttrOr("href", "")); id == "" {
			errMsg = fmt.Errorf("Failed to find course id from href (%s).", s.AttrOr("href", ""))
		} else if name, seq, sem := ParseCourseName(s.Text()); name == "" {
			errMsg = fmt.Errorf("Failed to parse course name, sequence number and semester id from link text (%s)", s.Text())
		} else {
			course := &model.Course{
				Id:         id,
				SemesterId: sem,
				Sequence:   seq,
				Name:       name,
			}

			courses[i] = course
		}

		return errMsg == nil
	})

	if errMsg != nil {
		status = http.StatusInternalServerError
		errMsg = fmt.Errorf("Failed to parse %s: %s", url, errMsg)
		return
	}

	return
}

func (ada *Adapter) AllAttendedList() (courses []*model.Course, status int, errMsg error) {
	var thisCourses, pastCourses []*model.Course

	sg := util.NewStatusGroup()

	sg.Go(func(status *int, err *error) {
		thisCourses, *status, *err = ada.AttendedList(1)
	})
	sg.Go(func(status *int, err *error) {
		pastCourses, *status, *err = ada.AttendedList(2)
	})

	status, errMsg = sg.Wait()
	courses = append(thisCourses, pastCourses...)
	return
}
