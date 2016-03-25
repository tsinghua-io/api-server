package cic

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
	"strconv"
)

func TimeLocationsURL(courseId string) string {
	return fmt.Sprintf("%s/b/course/info/timePlace/%s", BaseURL, courseId)
}

func AssistantsURL(courseId string) string {
	return fmt.Sprintf("%s/b/mycourse/AssistTeacher/list/%s", BaseURL, courseId)
}

func AttendedURL(semesterID string) string {
	return fmt.Sprintf("%s/b/myCourse/courseList/loadCourse4Student/%s", BaseURL, semesterID)
}

func (ada *Adapter) TimeLocations(courseId string, _ map[string]string, timeLocations *[]*resource.TimeLocation) (status int) {
	if timeLocations == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := TimeLocationsURL(courseId)
	var v struct {
		ResultList []struct {
			Skzc string
			Skxq string
			Skjc string
			Skdd string
		}
	}

	if err := ada.GetJSON("GET", url, &v); err != nil {
		return http.StatusBadGateway
	}

	for _, result := range v.ResultList {
		dayOfWeek, err := strconv.Atoi(result.Skxq)
		if err != nil {
			return http.StatusBadGateway
		}
		periodOfDay, err := strconv.Atoi(result.Skjc)
		if err != nil {
			return http.StatusBadGateway
		}

		timeLocation := &resource.TimeLocation{
			Weeks:       result.Skzc,
			DayOfWeek:   dayOfWeek,
			PeriodOfDay: periodOfDay,
			Location:    result.Skdd,
		}
		*timeLocations = append(*timeLocations, timeLocation)
	}

	return http.StatusOK
}

func (ada *Adapter) Assistants(courseId string, _ map[string]string, assistants *[]*resource.User) (status int) {
	if assistants == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := AssistantsURL(courseId)
	var v struct {
		ResultList []struct {
			Id     string
			Dwmc   string
			Phone  string
			Email  string
			Name   string
			Gender string
		}
	}

	if err := ada.GetJSON("GET", url, &v); err != nil {
		return http.StatusBadGateway
	}

	for _, result := range v.ResultList {
		assistant := &resource.User{
			Id:         result.Id,
			Name:       result.Name,
			Department: result.Dwmc,
			Gender:     result.Gender,
			Email:      result.Email,
			Phone:      result.Phone,
		}
		*assistants = append(*assistants, assistant)
	}

	return http.StatusOK
}

func (ada *Adapter) Attended(semesterID string, params map[string]string, courses *[]*resource.Course) (status int) {
	if courses == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := AttendedURL(semesterID)
	var v struct {
		ResultList []struct {
			CourseId      string
			Course_no     string
			Course_seq    string
			Course_name   string
			E_course_name string
			TeacherInfo   struct {
				Id     string
				Name   string
				Email  string
				Phone  string
				Gender string
				Title  string
			}
			CodeDepartmentInfo struct {
				Dwmc   string
				Dwywmc string
			}
			SemesterInfo struct {
				Id string
			}
			Detail_c    string
			Detail_e    string
			Credit      int
			Course_time int
		}
	}

	if err := ada.GetJSON("GET", url, &v); err != nil {
		return http.StatusBadGateway
	}

	statuses := make(chan int, 1)
	count := 0

	// TODO: Here we loop through a struct array. Will Go copy every struct?
	// Try some benchmarks.
	for _, result := range v.ResultList {
		// Language specific fields.
		// TODO: Move out of loop?
		var name, description, department string
		switch params["lang"] {
		case "zh-CN", "":
			name = result.Course_name
			description = result.Detail_c
			department = result.CodeDepartmentInfo.Dwmc
		case "en":
			name = result.E_course_name
			description = result.Detail_e
			department = result.CodeDepartmentInfo.Dwywmc
		}

		course := &resource.Course{
			Id:             result.CourseId,
			Semester:       result.SemesterInfo.Id,
			CourseNumber:   result.Course_no,
			CourseSequence: result.Course_seq,
			Name:           name,
			Credit:         result.Credit,
			Hour:           result.Course_time,
			Description:    description,

			Teachers: []*resource.User{
				&resource.User{
					Id:         result.TeacherInfo.Id,
					Name:       result.TeacherInfo.Name,
					Type:       result.TeacherInfo.Title,
					Department: department,
					Gender:     result.TeacherInfo.Gender,
					Email:      result.TeacherInfo.Email,
					Phone:      result.TeacherInfo.Phone,
				},
			},
		}

		count += 2
		go func() {
			statuses <- ada.TimeLocations(course.Id, params, &course.TimeLocations)
		}()
		go func() {
			statuses <- ada.Assistants(course.Id, params, &course.Assistants)
		}()
		*courses = append(*courses, course)
	}

	status = http.StatusOK
	for i := 0; i < count; i++ {
		if s := <-statuses; s != http.StatusOK {
			status = s
		}
	}

	return status
}
