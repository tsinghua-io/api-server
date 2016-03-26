package learn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func AttendedURL(page int) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/MyCourse.jsp?typepage=%d&language=cn", BaseURL, page)
}

func CourseInfoURL(courseId string) string {
	return fmt.Sprintf("%s/MultiLanguage/lesson/student/course_info.jsp?course_id=%s", BaseURL, courseId)
}

func isNewId(id string) bool {
	return strings.Contains(id, "-")
}

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

func (ada *Adapter) parseCourseLink(s *goquery.Selection, course *resource.Course) (status int) {
	if course == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	var tempCourse resource.Course
	text := s.Text()
	var err error

	// Course ID.
	href, _ := s.Attr("href")
	tempCourse.Id, err = URL2CourseId(href)
	if err != nil {
		glog.Errorf("Failed to parse course id from href (%s): %s", href, err)
		return http.StatusBadGateway
	}

	// Semester.
	tempCourse.Semester, err = courseName2Semester(text)
	if err != nil {
		glog.Errorf("Failed to parse semester from course name (%s): %s", text, err)
		return http.StatusBadGateway
	}

	// Info.
	if !isNewId(tempCourse.Id) {
		if status := ada.CourseInfo(tempCourse.Id, nil, &tempCourse); status != http.StatusOK {
			return status
		}
	}

	// We are safe.
	*course = tempCourse
	return http.StatusOK
}

func (ada *Adapter) CourseInfo(courseId string, _ map[string]string, course *resource.Course) (status int) {
	if course == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := CourseInfoURL(courseId)
	doc, err := ada.GetDocument(url)
	if err != nil {
		return http.StatusBadGateway
	}

	tds := doc.Find("table#table_box td")

	infos := tds.Map(func(i int, s *goquery.Selection) (info string) {
		return strings.TrimSpace(s.Text())
	})
	if len(infos) < 23 {
		glog.Errorf("Course information parsing error: cannot parse all the informations from %s", infos)
		return http.StatusBadGateway
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

	return http.StatusOK
}

func (ada *Adapter) attended(page int, _ map[string]string, courses *[]*resource.Course) (status int) {
	if courses == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := AttendedURL(page)
	doc, err := ada.GetDocument(url)
	if err != nil {
		return http.StatusBadGateway
	}

	links := doc.Find("#info_1 tr a")
	rows := links.Size()
	*courses = make([]*resource.Course, rows)

	// Parse each row.
	statuses := make(chan int, 1)

	links.Each(func(i int, s *goquery.Selection) {
		course := new(resource.Course)
		go func() {
			statuses <- ada.parseCourseLink(s, course)
		}()
		(*courses)[i] = course
	})

	for i := 0; i < rows; i++ {
		status = adapter.MergeStatus(status, <-statuses)
	}

	return status
}

func (ada *Adapter) Attended(_ string, params map[string]string, courses *[]*resource.Course) (status int) {
	switch semester := params["semester"]; semester {
	case "", "all":
		var page1, page2 []*resource.Course
		statuses := make(chan int, 1)

		go func() {
			statuses <- ada.attended(1, params, &page1)
		}()
		go func() {
			statuses <- ada.attended(2, params, &page2)
		}()

		status1 := <-statuses
		status2 := <-statuses
		if status1 != http.StatusOK {
			return status1
		}
		if status2 != http.StatusOK {
			return status2
		}

		*courses = append(page1, page2...)
		return http.StatusOK
	default:
		glog.Errorf("Semester (%s) is not supported by old adapter.	", semester)
		return http.StatusNotImplemented
	}
}

func (ada *Adapter) CourseIdMap(userId string, params map[string]string, cic2old map[string]string) (status int) {
	var courses []*resource.Course
	if status := ada.Attended(userId, params, &courses); status != http.StatusOK {
		return status
	}

	for _, course := range courses {
		if !isNewId(course.Id) {
			cicId := course.Semester + "-" + course.CourseNumber + "-" + course.CourseSequence
			cic2old[cicId] = course.Id
		}
	}

	return http.StatusOK
}
