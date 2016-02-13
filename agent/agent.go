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
	oldSession, ok := context.GetOk(r, "oldSession")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	oldAda := old.New(oldSession.([]*http.Cookie))
	res, status := oldAda.PersonalInfo()

	if status != http.StatusOK {
		// clear the cookie
		clearSession, ok := context.GetOk(r, "clearSession")
		clearSessionFunc := clearSession.(func(bool) bool)
		if !ok {
			glog.Warningln("No clearSession func in the request context.")
		}
		clearSessionFunc(false) // clear session of old learning web
	}

	w.WriteHeader(status)
	if res != nil {
		j, _ := json.Marshal(res)
		w.Write(j)
	}
}
