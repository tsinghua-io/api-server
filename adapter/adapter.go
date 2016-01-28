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
	GetUserInfo() (CommunicateUnit, error)
}

type AdapterNewer func(middleware.UserSession) Adapter

var AdapterNewerList = []AdapterNewer{NewOldAdapter, NewCicAdapter}

type CommunicateUnit struct {
	Resource interface{}
	Status   int    // response http status
	Session  string // the session used in this request
}
