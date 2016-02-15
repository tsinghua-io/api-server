package cic

import (
	"github.com/bitly/go-simplejson"
	"github.com/tsinghua-io/api-server/resource"
)

type userParser struct {
	Id         string
	Name       string
	Type       string
	Department string
	Class      string
	Gender     string
	Email      string
	Phone      string
}

type courseParser struct {
	Id             string
	Name           string
	Teacher        userParser
	Coteachers     userParser
	SchoolYear     string
	Semester       string
	CourseNumber   string
	CourseSequence string
	Credit         string
	Hour           string
	Description    string
	StudentCount   string
}

type announcementParser struct {
	Id        string
	CourseId  string
	Title     string
	Owner     userParser
	CreatedAt string
	Important string
	Body      string
}

type fileParser struct {
	Id          string
	CourseId    string
	Category    string
	Title       string
	Description string
	Filename    string
	Size        string
	DownloadUrl string
	Created_at  string
	Owner       userParser
}

type attachmentParser struct {
	Filename    string
	Size        string
	DownloadUrl string
}

type HomeworkParser struct {
	Id              string
	CourseId        string
	Title           string
	CreatedAt       string
	BeginAt         string
	DueAt           string
	SubmissionCount string
	MarkCount       string
	Body            string
	Attachment      attachmentParser
}

type submissionParser struct {
	CourseId          string
	HomeworkId        string
	Student           userParser
	CreatedAt         string
	MarkedAt          string
	Score             string
	Body              string
	Attachment        attachmentParser
	Comment           string
	CommentAttachment attachmentParser
}

// parseUser reads a User from a json, using the given paths.
func (parser *userParser) parse(j *simplejson.Json) (user *resource.User, err error) {
	tempUser := &resource.User{}
	if parser.Id != "" {
		if tempUser.Id, err = j.GetPath(parser.Id).String(); err != nil {
			return
		}
	}
	if parser.Name != "" {
		if tempUser.Name, err = j.GetPath(parser.Name).String(); err != nil {
			return
		}
	}
	if parser.Department != "" {
		if tempUser.Department, err = j.GetPath(parser.Department).String(); err != nil {
			return
		}
	}
	if parser.Class != "" {
		if tempUser.Class, err = j.GetPath(parser.Class).String(); err != nil {
			return
		}
	}
	if parser.Gender != "" {
		if tempUser.Gender, err = j.GetPath(parser.Gender).String(); err != nil {
			return
		}
	}
	if parser.Email != "" {
		if tempUser.Email, err = j.GetPath(parser.Email).String(); err != nil {
			return
		}
	}
	if parser.Phone != "" {
		if tempUser.Phone, err = j.GetPath(parser.Phone).String(); err != nil {
			return
		}
	}

	// Safe, we are done.
	user = tempUser
	return
}
