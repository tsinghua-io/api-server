package learn

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
)

func UserURL(_ string) string {
	return fmt.Sprintf("%s/b/m/getStudentById", BaseURL)
}

func (ada *Adapter) User(userId string, _ map[string]string, user *resource.User) (status int) {
	if user == nil {
		glog.Errorf("nil received")
		return http.StatusInternalServerError
	}

	url := UserURL(userId)
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

	if err := ada.GetJSON("POST", url, &v); err != nil {
		return http.StatusBadGateway
	}

	data := v.DataSingle
	*user = resource.User{
		Id:         data.Id,
		Name:       data.Name,
		Type:       data.Title,
		Department: data.MajorName,
		Class:      data.Classname,
		Gender:     data.Gender,
		Email:      data.Email,
		Phone:      data.Phone,
	}

	return http.StatusOK
}
