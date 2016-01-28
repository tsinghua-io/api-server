// Resources Agent

package agent

import (
	"encoding/json"
	"github.com/golang/glog"
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

	setSession, ok := context.GetOk(r, "setSession")
	setSessionFunc := setSession.(func(string, bool) bool)
	if !ok {
		glog.Warningln("No setSession func in the request context.")
	}

	oldAda := adapter.NewOldAdapter(userSession.(middleware.UserSession))
	ans, err := oldAda.GetUserInfo()

	if err != nil {
		glog.Warningln("Failed getting user info using old learning web. ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update session
	if ans.Session != userSession.(middleware.UserSession).Session {
		setSessionFunc(ans.Session, false)
	}
	j, _ := json.Marshal(ans.Resource)
	w.WriteHeader(ans.Status)
	w.Write(j)
}
