package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	TimePlaceURL  = BaseURL + "/b/course/info/timePlace/{course_id}"
	AssistantsURL = BaseURL + "/b/mycourse/AssistTeacher/list/{course_id}"
	AttendedURL   = BaseURL + "/b/myCourse/courseList/loadCourse4Student/-1"
)

type timeLocationParser struct {
	ResultList []struct {
		Skzc string
		Skxq string
		Skjc string
		Skdd string
	}
}

func (p *timeLocationParser) parse(r io.Reader, info interface{}, _ string) error {
	timeLocations, ok := info.(*[]*resource.TimeLocation)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	for _, result := range p.ResultList {
		dayOfWeek, err := strconv.Atoi(result.Skxq)
		if err != nil {
			return fmt.Errorf("Failed to parse DayOfWeek to int: %s", err)
		}
		periodOfDay, err := strconv.Atoi(result.Skjc)
		if err != nil {
			return fmt.Errorf("Failed to parse PeriodOfDay to int: %s", err)
		}

		timeLocation := &resource.TimeLocation{
			Weeks:       result.Skzc,
			DayOfWeek:   dayOfWeek,
			PeriodOfDay: periodOfDay,
			Location:    result.Skdd,
		}
		*timeLocations = append(*timeLocations, timeLocation)
	}

	return nil
}

type assistantsParser struct {
	ResultList []struct {
		id     string
		dwmc   string
		phone  string
		email  string
		name   string
		gender string
	}
}

func (p *assistantsParser) parse(r io.Reader, info interface{}, _ string) error {
	users, ok := info.(*[]*resource.User)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	for _, result := range p.ResultList {
		user := &resource.User{
			Id:         result.id,
			Name:       result.name,
			Department: result.dwmc,
			Gender:     result.gender,
			Email:      result.email,
			Phone:      result.phone,
		}
		*users = append(*users, user)
	}

	return nil
}

type coursesParser struct {
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
			SemesterName  string
			SemesterEname string
		}
		Detail_c    string
		Detail_e    string
		Credit      int
		Course_time int
	}
}

func (p *coursesParser) parse(r io.Reader, info interface{}, langCode string) error {
	courses, ok := info.(*[]*resource.Course)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	// TODO: Here we loop through a struct array. Will Go copy every struct?
	// Try some benchmarks.
	for _, result := range p.ResultList {
		// Language specific fields.
		// TODO: Move out of loop?
		var semester, name, description, department string
		switch langCode {
		case "zh-CN":
			semester = result.SemesterInfo.SemesterName
			name = result.Course_name
			description = result.Detail_c
			department = result.CodeDepartmentInfo.Dwmc
		case "en":
			semester = result.SemesterInfo.SemesterEname
			name = result.E_course_name
			description = result.Detail_e
			department = result.CodeDepartmentInfo.Dwywmc
		}

		course := &resource.Course{
			Id:             result.CourseId,
			Semester:       semester,
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
		*courses = append(*courses, course)
	}

	return nil
}

// func (adapter *CicAdapter) Attending() (courses []*resource.Course, status int) {
//  status = adapter.FetchInfo(PersonalInfoURL, "POST", &attendingParser{}, &courses)
//  return
// }

func (adapter *CicAdapter) Attended() (courses []*resource.Course, status int) {
	if status = adapter.FetchInfo(AttendedURL, "GET", &coursesParser{}, &courses); status != http.StatusOK {
		return nil, status
	}
	chan_size := len(courses) * 2
	statuses := make(chan int, chan_size)

	for _, course := range courses {
		course := course // Avoid variable reusing.

		// Time & Place.
		go func() {
			url := strings.Replace(TimePlaceURL, "{course_id}", course.Id, -1)
			statuses <- adapter.FetchInfo(url, "GET", &timeLocationParser{}, &course.TimeLocations)
		}()

		// Assistants.
		go func() {
			url := strings.Replace(AssistantsURL, "{course_id}", course.Id, -1)
			statuses <- adapter.FetchInfo(url, "GET", &assistantsParser{}, &course.Assistants)
		}()
	}

	// Drain the channel.
	for i := 0; i < chan_size; i++ {
		if status = <-statuses; status != http.StatusOK {
			return nil, status
		}
	}

	return
}
