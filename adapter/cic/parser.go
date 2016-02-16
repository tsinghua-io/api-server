package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
)

type parser interface {
	parse(reader io.Reader, info interface{}) error
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
			Dwjc string
		}
		Detail_c    string
		Credit      int
		Course_time int
	}
}

type announcementsParser struct {
}

type filesParser struct {
}

type homeworksParser struct {
}

func (p *personalInfoParser) parse(r io.Reader, info interface{}) error {
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

func (p *courseListParser) parse(r io.Reader, info interface{}) error {
	courses, ok := info.(*[]*resource.Course)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	// s, _ := ioutil.ReadAll(r)
	// fmt.Println(string(s))

	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		return err
	}

	for _, result := range p.ResultList {
		course := &resource.Course{
			Id:   result.CourseId,
			Name: result.Course_name,
			Teacher: resource.User{
				Id:         result.TeacherInfo.Id,
				Name:       result.TeacherInfo.Name,
				Type:       result.TeacherInfo.Title,
				Department: result.CodeDepartmentInfo.Dwjc,
				Gender:     result.TeacherInfo.Gender,
				Email:      result.TeacherInfo.Email,
				Phone:      result.TeacherInfo.Phone,
			},
			CourseNumber:   result.Course_no,
			CourseSequence: result.Course_seq,
			Credit:         result.Credit,
			Hour:           result.Course_time,
			Description:    result.Detail_c,
		}
		*courses = append(*courses, course)
	}

	return nil
}
