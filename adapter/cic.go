package adapter

import (
	"github.com/tsinghua-io/api-server/middleware"
	"github.com/tsinghua-io/api-server/resources"
	"net/http"
)

// CicAdapter is the adapter for learn.cic.tsinghua.edu.cn
type CicAdapter struct {
	userSession middleware.UserSession
	queue       chan CommunicateUnit
}

func NewCicAdapter(userSession middleware.UserSession,
	queue chan CommunicateUnit) Adapter {
	return CicAdapter{userSession, queue}
}

func (ada CicAdapter) GetUserInfo() {
	ada.queue <- CommunicateUnit{resources.SUser{
		Id:         "2012011067",
		Name:       "宁雪妃ci",
		Department: "電子工程",
		User_type:  "undergraduate"},
		http.StatusUnauthorized}
}
