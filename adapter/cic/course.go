package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	timePlaceURL        = BaseURL + "/b/course/info/timePlace/{course_id}"
	courseAssistantsURL = BaseURL + "/b/mycourse/AssistTeacher/list/{course_id}"
	attendedURL         = BaseURL + "/b/myCourse/courseList/loadCourse4Student/-1"
)

type timeLocationParser struct {
	params map[string]string
	data   struct {
		ResultList []struct {
			Skzc string
			Skxq string
			Skjc string
			Skdd string
		}
	}
}

func (p *timeLocationParser) Parse(r io.Reader, info interface{}) error {
	timeLocations, ok := info.(*[]*resource.TimeLocation)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p.data); err != nil {
		return err
	}

	for _, result := range p.data.ResultList {
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
	params map[string]string
	data   struct {
		ResultList []struct {
			Id     string
			Dwmc   string
			Phone  string
			Email  string
			Name   string
			Gender string
		}
	}
}

func (p *assistantsParser) Parse(r io.Reader, info interface{}) error {
	users, ok := info.(*[]*resource.User)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p.data); err != nil {
		return err
	}

	for _, result := range p.data.ResultList {
		user := &resource.User{
			Id:         result.Id,
			Name:       result.Name,
			Department: result.Dwmc,
			Gender:     result.Gender,
			Email:      result.Email,
			Phone:      result.Phone,
		}
		*users = append(*users, user)
	}

	return nil
}

type coursesParser struct {
	params map[string]string
	data   struct {
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
}

func (p *coursesParser) Parse(r io.Reader, info interface{}) error {
	courses, ok := info.(*[]*resource.Course)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p.data); err != nil {
		return err
	}

	// TODO: Here we loop through a struct array. Will Go copy every struct?
	// Try some benchmarks.
	for _, result := range p.data.ResultList {
		// Language specific fields.
		// TODO: Move out of loop?
		var semester, name, description, department string
		switch p.params["lang"] {
		case "zh-CN", "":
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

func (ada *CicAdapter) Attended(username string, params map[string]string) (courses []*resource.Course, status int) {
	if username != "" {
		// Not suppported yet.
		return nil, http.StatusBadRequest
	}

	parser := &coursesParser{params: params}
	if status = adapter.FetchInfo(&ada.client, attendedURL, "GET", parser, &courses); status != http.StatusOK {
		return nil, status
	}
	chan_size := len(courses) * 2
	statuses := make(chan int, chan_size)

	for _, course := range courses {
		course := course // Avoid variable reusing.

		// Time & Place.
		go func() {
			URL := strings.Replace(timePlaceURL, "{course_id}", course.Id, -1)
			parser := &timeLocationParser{params: params}
			statuses <- adapter.FetchInfo(&ada.client, URL, "GET", parser, &course.TimeLocations)
		}()

		// Assistants.
		go func() {
			URL := strings.Replace(courseAssistantsURL, "{course_id}", course.Id, -1)
			parser := &assistantsParser{params: params}
			statuses <- adapter.FetchInfo(&ada.client, URL, "GET", parser, &course.Assistants)
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
