package adapter

import (
	"github.com/tsinghua-io/api-server/middleware"
	//"github.com/tsinghua-io/api-server/resources"
	//"net/http"
)

// CicAdapter is the adapter for learn.cic.tsinghua.edu.cn
type CicAdapter struct {
	userSession middleware.UserSession
}

func NewCicAdapter(userSession middleware.UserSession) Adapter {
	return CicAdapter{userSession}
}

func (ada CicAdapter) GetUserInfo() (CommunicateUnit, error) {
	return CommunicateUnit{}, nil
	// ada.queue <- CommunicateUnit{resources.SUser{
	//	Id:         "2012011067",
	//	Name:       "宁雪妃ci",
	//	Department: "電子工程",
	//	User_type:  "undergraduate"},
	//	http.StatusUnauthorized,
	//	""}
}
