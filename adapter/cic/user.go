package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
)

const (
	PersonalInfoURL = BaseURL + "/b/m/getStudentById"
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

func (adapter *CicAdapter) Profile(username string) (user *resource.User, status int) {
	if username == "" {
		// Self Profile.
		user = &resource.User{}
		status = adapter.FetchInfo(PersonalInfoURL, "POST", &personalInfoParser{}, user)
		return user, status
	} else {
		// User Profile, not implemented yet.
		return nil, http.StatusBadRequest
	}
}
