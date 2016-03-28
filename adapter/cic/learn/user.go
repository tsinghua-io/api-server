package learn

import (
	"github.com/tsinghua-io/api-server/model"
	"net/http"
)

const (
	ProfileURL = BaseURL + "/b/m/getStudentById"
)

func (ada *Adapter) Profile() (profile *model.User, status int) {
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

	if err := ada.GetJSON("POST", ProfileURL, &v); err != nil {
		status = http.StatusBadGateway
		return
	}

	data := v.DataSingle
	profile = &model.User{
		Id:         data.Id,
		Name:       data.Name,
		Type:       data.Title,
		Department: data.MajorName,
		Class:      data.Classname,
		Gender:     data.Gender,
		Email:      data.Email,
		Phone:      data.Phone,
	}

	status = http.StatusOK
	return
}
