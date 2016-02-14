// Resources Agent

package agent

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/tsinghua-io/api-server/adapter/old"
	//"github.com/tsinghua-io/api-server/adapter/cic"
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
	session, ok := context.GetOk(r, "session")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ada := old.New(session.([]*http.Cookie))

	res, status := ada.PersonalInfo()

	if status != http.StatusOK {
		// clear the cookie
		clearSession, ok := context.GetOk(r, "clearSession")
		clearSessionFunc := clearSession.(func() bool)
		if !ok {
			glog.Warningln("No clearSession func in the request context.")
		}
		clearSessionFunc() // clear session of learning web
	}

	w.WriteHeader(status)
	if res != nil {
		j, _ := json.Marshal(res)
		w.Write(j)
	}
}
