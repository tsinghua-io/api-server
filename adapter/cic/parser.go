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

type courseParser struct {
}

type announcementParser struct {
}

type fileParser struct {
}

type attachmentParser struct {
}

type homeworkParser struct {
}

type submissionParser struct {
}

// parseUser reads a User from a json, using the given paths.
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

// // parseUser reads a User from a json, using the given paths.
// func (parser *courseParser) parse(reader *io.Reader) (course *resource.Course, err error) {
// }
