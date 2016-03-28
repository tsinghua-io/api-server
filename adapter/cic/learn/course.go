package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
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

func (ada *Adapter) TimeLocations(courseId string) (timeLocations []*model.TimeLocation, status int) {
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
		status = http.StatusBadGateway
		return
	}

	status = http.StatusInternalServerError

	for _, result := range v.ResultList {
		dayOfWeek, err := strconv.Atoi(result.Skxq)
		if err != nil {
			return
		}
		periodOfDay, err := strconv.Atoi(result.Skjc)
		if err != nil {
			return
		}

		timeLocation := &model.TimeLocation{
			Weeks:       result.Skzc,
			DayOfWeek:   dayOfWeek,
			PeriodOfDay: periodOfDay,
			Location:    result.Skdd,
		}
		timeLocations = append(timeLocations, timeLocation)
	}

	status = http.StatusOK
	return
}

func (ada *Adapter) Assistants(courseId string) (assistants []*model.User, status int) {
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
		status = http.StatusBadGateway
		return
	}

	for _, result := range v.ResultList {
		assistant := &model.User{
			Id:         result.Id,
			Name:       result.Name,
			Department: result.Dwmc,
			Gender:     result.Gender,
			Email:      result.Email,
			Phone:      result.Phone,
		}
		assistants = append(assistants, assistant)
	}

	status = http.StatusOK
	return
}

func (ada *Adapter) Attended(semesterID string, english bool) (courses []*model.Course, status int) {
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
		status = http.StatusBadGateway
		return
	}

	sg := util.NewStatusGroup()

	for _, result := range v.ResultList {
		var name, description, department string
		if english {
			name = result.E_course_name
			description = result.Detail_e
			department = result.CodeDepartmentInfo.Dwywmc
		} else {
			name = result.Course_name
			description = result.Detail_c
			department = result.CodeDepartmentInfo.Dwmc
		}

		course := &model.Course{
			Id:             result.CourseId,
			Semester:       result.SemesterInfo.Id,
			CourseNumber:   result.Course_no,
			CourseSequence: result.Course_seq,
			Name:           name,
			Credit:         result.Credit,
			Hour:           result.Course_time,
			Description:    description,

			Teachers: []*model.User{
				&model.User{
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

		sg.Add(2)
		go func() {
			var status int
			defer sg.Done(status)
			course.TimeLocations, status = ada.TimeLocations(course.Id)
		}()
		go func() {
			var status int
			defer sg.Done(status)
			course.Assistants, status = ada.Assistants(course.Id)
		}()
		courses = append(courses, course)
	}

	status = sg.Wait()
	return
}

func (ada *Adapter) NowAttended(english bool) (thisCourses []*model.Course, nextCourses []*model.Course, status int) {
	var thisSem, nextSem string
	thisSem, nextSem, status = ada.Semesters()
	if status != http.StatusOK {
		return
	}

	sg := util.NewStatusGroup()
	sg.Add(2)
	go func() {
		var status int
		defer sg.Done(status)
		thisCourses, status = ada.Attended(thisSem, english)
	}()
	go func() {
		var status int
		defer sg.Done(status)
		nextCourses, status = ada.Attended(nextSem, english)
	}()

	status = sg.Wait()
	return
}

func (ada *Adapter) AllAttended(english bool) (courses []*model.Course, status int) {
	var pastCourses, thisCourses, nextCourses []*model.Course

	sg := util.NewStatusGroup()
	sg.Add(2)
	go func() {
		var status int
		defer sg.Done(status)
		thisCourses, nextCourses, status = ada.NowAttended(english)
	}()
	go func() {
		var status int
		defer sg.Done(status)
		pastCourses, status = ada.Attended("-1", english)
	}()

	status = sg.Wait()
	courses = append(nextCourses, thisCourses...)
	courses = append(courses, pastCourses...)
	return
}
