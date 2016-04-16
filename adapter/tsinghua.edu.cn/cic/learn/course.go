package learn

import (
	"fmt"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
	"strconv"
)

func SchedulesURL(courseId string) string {
	return fmt.Sprintf("%s/b/course/info/timePlace/%s", BaseURL, courseId)
}

func AssistantsURL(courseId string) string {
	return fmt.Sprintf("%s/b/mycourse/AssistTeacher/list/%s", BaseURL, courseId)
}

func AttendedURL(semesterId string) string {
	return fmt.Sprintf("%s/b/myCourse/courseList/loadCourse4Student/%s", BaseURL, semesterId)
}

func (ada *Adapter) Schedules(courseId string) (schedules []*model.Schedule, status int, errMsg error) {
	schedules = make([]*model.Schedule, 0)
	status = http.StatusOK

	url := SchedulesURL(courseId)
	var v struct {
		ResultList []struct {
			Skzc string
			Skxq string
			Skjc string
			Skdd string
		}
	}

	if err := ada.GetJSON(url, &v); err != nil {
		return nil, http.StatusBadGateway, err
	}

	for _, result := range v.ResultList {
		day, err := strconv.Atoi(result.Skxq)
		if err != nil {
			status = http.StatusInternalServerError
			errMsg = fmt.Errorf("Failed to parse day of week from %s: %s", result.Skxq, err)
			return
		}
		slot, err := strconv.Atoi(result.Skjc)
		if err != nil {
			status = http.StatusInternalServerError
			errMsg = fmt.Errorf("Failed to parse period of day from %s: %s", result.Skjc, err)
			return
		}

		schedule := &model.Schedule{
			Weeks:    result.Skzc,
			Day:      day,
			Slot:     slot,
			Location: result.Skdd,
		}
		schedules = append(schedules, schedule)
	}

	return
}

func (ada *Adapter) Assistants(courseId string) (assistants []*model.User, status int, errMsg error) {
	assistants = make([]*model.User, 0)
	status = http.StatusOK

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

	if err := ada.GetJSON(url, &v); err != nil {
		return nil, http.StatusBadGateway, err
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

	return
}

func (ada *Adapter) Attended(semesterId string, english bool) (courses []*model.Course, status int, errMsg error) {
	courses = make([]*model.Course, 0)

	url := AttendedURL(semesterId)
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

	if err := ada.GetJSON(url, &v); err != nil {
		return nil, http.StatusBadGateway, err
	}

	sg := util.NewStatusGroup()

	for _, result := range v.ResultList {
		var name, description, department string
		if english {
			name = result.E_course_name
			description = result.Detail_e
			department = result.CodeDepartmentInfo.Dwywmc
		}
		// Fallback to Chinese version.
		if name == "" {
			name = result.Course_name
		}
		if description == "" {
			description = result.Detail_c
		}
		if department == "" {
			department = result.CodeDepartmentInfo.Dwmc
		}

		course := &model.Course{
			Id:          result.CourseId,
			SemesterId:  result.SemesterInfo.Id,
			Number:      result.Course_no,
			Sequence:    result.Course_seq,
			Name:        name,
			Credit:      result.Credit,
			Hour:        result.Course_time,
			Description: description,

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

		sg.Go(func(status *int, err *error) {
			course.Schedules, *status, *err = ada.Schedules(course.Id)
		})
		sg.Go(func(status *int, err *error) {
			course.Assistants, *status, *err = ada.Assistants(course.Id)
		})
		courses = append(courses, course)
	}

	status, errMsg = sg.Wait()
	return
}

func (ada *Adapter) NowAttended(english bool) (thisCourses []*model.Course, nextCourses []*model.Course, status int, errMsg error) {
	var thisSem, nextSem *model.Semester
	thisSem, nextSem, status, errMsg = ada.Semesters()
	if errMsg != nil {
		return
	}

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		thisCourses, *status, *err = ada.Attended(thisSem.Id, english)
	})
	sg.Go(func(status *int, err *error) {
		nextCourses, *status, *err = ada.Attended(nextSem.Id, english)
	})

	status, errMsg = sg.Wait()
	return
}

func (ada *Adapter) PastAttended(english bool) (courses []*model.Course, status int, errMsg error) {
	return ada.Attended("-1", english)
}

func (ada *Adapter) AllAttended(english bool) (courses []*model.Course, status int, errMsg error) {
	var pastCourses, thisCourses, nextCourses []*model.Course

	sg := util.NewStatusGroup()
	sg.Go(func(status *int, err *error) {
		thisCourses, nextCourses, *status, *err = ada.NowAttended(english)
	})
	sg.Go(func(status *int, err *error) {
		pastCourses, *status, *err = ada.PastAttended(english)
	})

	status, errMsg = sg.Wait()
	courses = append(nextCourses, thisCourses...)
	courses = append(courses, pastCourses...)
	return
}
