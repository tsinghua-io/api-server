package cic

import (
	"encoding/json"
	"fmt"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/resource"
	"io"
	"net/http"
)

const (
	profileURL = BaseURL + "/b/m/getStudentById"
)

type profileParser struct {
	params map[string]string
	data   struct {
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
}

func (p *profileParser) Parse(r io.Reader, info interface{}) error {
	user, ok := info.(*resource.User)
	if !ok {
		return fmt.Errorf("The parser and the destination type do not match.")
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p.data); err != nil {
		return err
	}

	user.Id = p.data.DataSingle.Id
	user.Name = p.data.DataSingle.Name
	user.Type = p.data.DataSingle.Title
	user.Department = p.data.DataSingle.MajorName
	user.Class = p.data.DataSingle.Classname
	user.Gender = p.data.DataSingle.Gender
	user.Email = p.data.DataSingle.Email
	user.Phone = p.data.DataSingle.Phone

	return nil
}

func (ada *CicAdapter) Profile(username string, params map[string]string) (user *resource.User, status int) {
	if username != "" {
		// Not supported yet.
		return nil, http.StatusBadRequest
	}

	// Self Profile.
	parser := &profileParser{params: params}
	user = &resource.User{}

	status = adapter.FetchInfo(&ada.client, profileURL, "POST", parser, user)
	return user, status
}
