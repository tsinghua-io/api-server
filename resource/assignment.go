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
	assignments, status, err := ada.Assignments(mux.Vars(req)["id"])
	util.JSON(rw, assignments, status, err)
})
