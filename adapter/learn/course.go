package old

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const (
	attendedURL   = BaseURL + "/MultiLanguage/lesson/student/MyCourse.jsp?typepage={page}&language=cn"
	courseInfoURL = BaseURL + "/MultiLanguage/lesson/student/course_info.jsp?course_id={course_id}"
)

func URL2CourseId(courseURL string) (courseId string, err error) {
	if courseURL == "" {
		return "", fmt.Errorf("Empty course URL.")
	}

	parsedURL, err := url.Parse(courseURL)
	if err != nil {
		return "", err
	}

	if courseId = parsedURL.Query().Get("course_id"); courseId != "" {
		return courseId, nil
	} else if courseId = path.Base(parsedURL.Path); courseId != "" {
		return courseId, nil
	} else {
		return "", fmt.Errorf("Unknown course URL format: %s", courseURL)
	}
}

func courseName2Semester(name string) (semester string, err error) {
	// Semester.
	begin := len(name) - 22
	semesterBegin := len(name) - 13
	semesterEnd := len(name) - 10
	end := len(name) - 1

	if len(name) < 23 || string(name[begin-1]) != "(" || string(name[end]) != ")" {
		return "", fmt.Errorf("Unknown course name format: %s", name)
	}

	yearStr := name[begin:semesterBegin]
	switch string(name[semesterBegin:semesterEnd]) {
	case "秋":
		return yearStr + "-1", nil
	case "春":
		return yearStr + "-2", nil
	case "夏":
		return yearStr + "-3", nil
	default:
		return "", fmt.Errorf("Unknown course name format: %s", name)
	}
}

func parseCourseLink(s *goquery.Selection, course *resource.Course) error {
	text := s.Text()

	// Course ID.
	href, _ := s.Attr("href")
	id, err := URL2CourseId(href)
	if err != nil {
		return err
	}

	// Semester.
	semester, err := courseName2Semester(text)
	if err != nil {
		return err
	}

	// We are safe.
	course.Id = id
	course.Semester = semester

	return nil
}

type CourseInfoParser struct {
	params map[string]string
}

func (p *CourseInfoParser) Parse(r io.Reader, info interface{}) error {
	course, ok := info.(*resource.Course)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	tds := doc.Find("table#table_box td")

	infos := tds.Map(func(i int, s *goquery.Selection) (info string) {
		return strings.TrimSpace(s.Text())
	})
	if len(infos) < 23 {
		return fmt.Errorf("Course information parsing error: cannot parse all the informations from %s", infos)
	}

	// Yes here's all we need for now.

	// credit, _ := strconv.Atoi(infos[7])
	// hour, _ := strconv.Atoi(infos[9])

	course.CourseNumber = infos[1]
	course.CourseSequence = infos[3]
	// course.Name = infos[5]
	// course.Credit = credit
	// course.Hour = hour
	// course.Description = infos[22]
	// course.Teachers = []*resource.User{
	// 	&resource.User{
	// 		Name:  infos[16],
	// 		Email: infos[18],
	// 		Phone: infos[20],
	// 	},
	// }

	return nil
}

type courseListParser struct {
	params map[string]string
}

func (p *courseListParser) Parse(r io.Reader, info interface{}) error {
	courses, ok := info.(*[]*resource.Course)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	links := doc.Find("#info_1 tr a")
	count := links.Size()
	*courses = make([]*resource.Course, count)

	links.EachWithBreak(func(i int, s *goquery.Selection) bool {
		course := &resource.Course{}
		if err = parseCourseLink(s, course); err != nil {
			return false
		}
		(*courses)[i] = course
		return true
	})
	return err
}

func (ada *OldAdapter) Attended(username string, params map[string]string) (courses []*resource.Course, status int) {
	if params["detail"] == "true" {
		glog.Errorf("Course detail is not supported by old adapter.")
		return nil, http.StatusInternalServerError
	}

	switch semester := params["semester"]; semester {
	case "", "all":
		var page1, page2 []*resource.Course
		status1 := make(chan int, 1)
		status2 := make(chan int, 1)

		go func() {
			var status int
			page1, status = ada.attendedPage(username, params, 1)
			status1 <- status
		}()
		go func() {
			var status int
			page2, status = ada.attendedPage(username, params, 2)
			status2 <- status
		}()

		if status := <-status1; status != http.StatusOK {
			return nil, status
		}
		if status := <-status2; status != http.StatusOK {
			return nil, status
		}

		return append(page1, page2...), http.StatusOK
	default:
		glog.Errorf("Semester (%s) is not supported by old adapter.	", semester)
		return nil, http.StatusNotImplemented
	}
}

func (ada *OldAdapter) attendedPage(_ string, params map[string]string, tab int) (courses []*resource.Course, status int) {
	URL := strings.Replace(attendedURL, "{page}", strconv.Itoa(tab), -1)
	parser := &courseListParser{params: params}

	if status = adapter.FetchInfo(&ada.client, URL, "GET", parser, &courses); status != http.StatusOK {
		return nil, status
	}

	count := len(courses)
	statuses := make(chan int, count)

	for _, course := range courses {
		course := course // Avoid variable reusing.

		go func() {
			// Skip new courses.
			if strings.Contains(course.Id, "-") {
				statuses <- http.StatusOK
				return
			}
			URL := strings.Replace(courseInfoURL, "{course_id}", course.Id, -1)
			parser := &CourseInfoParser{params: params}
			statuses <- adapter.FetchInfo(&ada.client, URL, "GET", parser, course)
		}()
	}

	// Drain the channel.
	for i := 0; i < count; i++ {
		if status = <-statuses; status != http.StatusOK {
			return nil, status
		}
	}

	return
}

func CourseIdMap(courses []*resource.Course) map[string]string {
	cic2old := make(map[string]string)

	for _, course := range courses {
		if strings.Contains(course.Id, "-") {
			continue // Already a cic id.
		}
		cicId := course.Semester + "-" + course.CourseNumber + "-" + course.CourseSequence
		cic2old[cicId] = course.Id
	}
	return cic2old
}
