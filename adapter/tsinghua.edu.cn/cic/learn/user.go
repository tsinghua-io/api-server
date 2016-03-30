package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"net/http"
	"strings"
)

const (
	ProfileURL = BaseURL + "/b/m/getStudentById"
)

func (ada *Adapter) Profile() (profile *model.User, status int, errMsg error) {
	status = http.StatusOK

	var v struct {
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

	if err := ada.PostFormJSON(ProfileURL, nil, &v); err != nil {
		return nil, http.StatusBadGateway, err
	}

	data := v.DataSingle
	profile = &model.User{
		Id:         data.Id,
		Name:       data.Name,
		Type:       data.Title,
		Department: data.MajorName,
		Class:      strings.TrimSpace(data.Classname),
		Gender:     data.Gender,
		Email:      data.Email,
		Phone:      data.Phone,
	}

	return
}
