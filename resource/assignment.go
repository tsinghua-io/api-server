package resource

import (
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/x/learn"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var CourseAssignments = Resource{
	"GET": util.AuthNeededHandler(GetCourseAssignments),
}

var GetCourseAssignments = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	v, status, err := BatchResourceFunc(
		mux.Vars(req)["id"],
		func(id string) (interface{}, int, error) {
			return ada.Assignments(id)
		})
	util.JSON(rw, v, status, err)
})
