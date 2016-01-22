// Resources Agent

package agent

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter"
	"github.com/tsinghua-io/api-server/middleware"
	"github.com/tsinghua-io/api-server/webapp"
	"net/http"
)

type userAgent struct {
}

var UserAgent = userAgent{}

func (useragent userAgent) BindRoute(app *webapp.WebApp) {
	app.HandleFunc("/users/me", useragent.GetInfo)
}

// GetInfo of userAgent get and update the Info field of this user
func (useragent *userAgent) GetInfo(w http.ResponseWriter, r *http.Request) {
	userSession, ok := context.GetOk(r, "userSession")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	queue := make(chan adapter.CommunicateUnit,
		len(adapter.AdapterNewerList))
	for _, adaNewer := range adapter.AdapterNewerList {
		ada := adaNewer(userSession.(middleware.UserSession), queue)
		go ada.GetUserInfo()
	}
	var ans adapter.CommunicateUnit
	var answerStatus bool
	for {
		ans = <-queue
		// fixme: status should be choosen as the most
		// meaningful status. 500 > 401 > 400
		// fixme: 304 not modified
		if ans.Status == http.StatusOK ||
			ans.Status == http.StatusCreated {
			j, _ := json.Marshal(ans.Resource)
			w.Write(j)
			answerStatus = true
			break
		}
	}
	if !answerStatus {
		w.WriteHeader(ans.Status)
	}

}
