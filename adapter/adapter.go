// Adapters

package adapter

import (
	"github.com/tsinghua-io/api-server/middleware"
)

// Websites urls.
const (
	OldUrl = "https://learn.tsinghua.edu.cn"
	CicUrl = "https://learn.cic.tsinghua.edu.cn"
)

type Adapter interface {
	GetUserInfo()
}

type AdapterNewer func(middleware.UserSession, chan CommunicateUnit) Adapter

var AdapterNewerList = []AdapterNewer{NewOldAdapter, NewCicAdapter}

type CommunicateUnit struct {
	Resource interface{}
	Status   int
}
