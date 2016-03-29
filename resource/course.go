package resource

import (
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/cic/learn"
	"github.com/tsinghua-io/api-server/util"
	"golang.org/x/text/language"
	"net/http"
)

var Attended = Resource{
	"GET": util.AuthNeededHandler(GetAttended),
}

var GetAttended = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	sem := req.URL.Query().Get("semester")
	english := (util.Language(req) == language.English)

	switch sem {
	case "", "latest":
		if this, next, status, err := ada.NowAttended(english); len(next) == 0 {
			util.JSON(rw, this, status, err)
		} else {
			util.JSON(rw, next, status, err)
		}
	case "this":
		this, _, status, err := ada.NowAttended(english)
		util.JSON(rw, this, status, err)
	case "next":
		_, next, status, err := ada.NowAttended(english)
		util.JSON(rw, next, status, err)
	case "past":
		past, status, err := ada.PastAttended(english)
		util.JSON(rw, past, status, err)
	case "all":
		all, status, err := ada.AllAttended(english)
		util.JSON(rw, all, status, err)
	default:
		courses, status, err := ada.Attended(sem, english)
		util.JSON(rw, courses, status, err)
	}
})
