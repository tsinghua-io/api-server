package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"strconv"
)

type parser interface {
	parse(reader io.Reader, info interface{}, langCode string) error
}

type personalInfoParser struct {
	DataSingle struct {
		Classname string
		Email     string
		Gender    string
		Id        string
		MajorName string
		Name      string
		Phone     string
		Title     string
	}
}

func (p *personalInfoParser) parse(r io.Reader, info interface{}, _ string) error {
	user, ok := info.(*resource.User)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	user.Id = p.DataSingle.Id
	user.Name = p.DataSingle.Name
	user.Type = p.DataSingle.Title
	user.Department = p.DataSingle.MajorName
	user.Class = p.DataSingle.Classname
	user.Gender = p.DataSingle.Gender
	user.Email = p.DataSingle.Email
	user.Phone = p.DataSingle.Phone

	return nil
}

type timeLocationParser struct {
	ResultList []struct {
		Skzc string
		Skxq string
		Skjc string
		Skdd string
	}
}

func (p *timeLocationParser) parse(r io.Reader, info interface{}, _ string) error {
	course, ok := info.(*resource.Course)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}
	if len(p.ResultList) == 0 {
		return nil // No time location info
	}

	var err error
	course.Weeks = p.ResultList[0].Skzc
	if course.DayOfWeek, err = strconv.Atoi(p.ResultList[0].Skxq); err != nil {
		return fmt.Errorf("Failed to parse DayOfWeek to int: %s", err)
	}
	if course.PeriodOfDay, err = strconv.Atoi(p.ResultList[0].Skjc); err != nil {
		return fmt.Errorf("Failed to parse PeriodOfDay to int: %s", err)
	}
	course.Location = p.ResultList[0].Skdd

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

type courseListParser struct {
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

func (p *courseListParser) parse(r io.Reader, info interface{}, langCode string) error {
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

// type announcementsParser struct {
// 	paginationList struct {
// 		recordList []struct {
// 			status       string
// 			courseNotice struct {
// 				id          string
// 				title       string
// 				owner       string
// 				regDate     string
// 				courseId    string
// 				msgPriority string
// 				detail      string
// 			}
// 		}
// 	}
// }

// type filesParser struct {
// }

// type homeworksParser struct {
// 	resultList []struct {
// 		courseHomeworkRecord struct {
// 		}
// 		courseHomeworkInfo struct {
// 		}
// 	}
// }
