package resource

import (
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/x/learn"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var CourseAnnouncements = Resource{
	"GET": util.AuthNeededHandler(GetCourseAnnouncements),
}

var GetCourseAnnouncements = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	announcements, status, err := ada.Announcements(mux.Vars(req)["id"])
	util.JSON(rw, announcements, status, err)
})
