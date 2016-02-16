package cic

import (
	"encoding/json"
	"github.com/tsinghua-io/api-server/resource"
	"io"
)

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
func (parser *personalInfoParser) parse(reader io.Reader) (user *resource.User, err error) {
	dec := json.NewDecoder(reader)
	if err = dec.Decode(parser); err != nil {
		return nil, err
	}

	return &resource.User{
		Id:         parser.DataSingle.Id,
		Name:       parser.DataSingle.Name,
		Type:       parser.DataSingle.Title,
		Department: parser.DataSingle.MajorName,
		Class:      parser.DataSingle.Classname,
		Gender:     parser.DataSingle.Gender,
		Email:      parser.DataSingle.Email,
		Phone:      parser.DataSingle.Phone,
	}, nil
}
