package adapter

import (
	"github.com/tsinghua-io/api-server/middleware"
	"github.com/tsinghua-io/api-server/resources"
	"net/http"
)

// OldAdapter is the adapter for learn.tsinghua.edu.cn
type OldAdapter struct {
	userSession middleware.UserSession
	queue       chan CommunicateUnit
}

func NewOldAdapter(userSession middleware.UserSession,
	queue chan CommunicateUnit) Adapter {
	return OldAdapter{userSession, queue}
}

func (ada OldAdapter) GetUserInfo() {
	ada.queue <- CommunicateUnit{resources.SUser{
		Id:         "2012011067",
		Name:       "宁雪妃",
		Department: "電子工程",
		User_type:  "undergraduate"},
		http.StatusOK}
}
