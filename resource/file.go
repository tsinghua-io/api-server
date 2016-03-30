package resource

import (
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/x/learn"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var CourseFiles = Resource{
	"GET": util.AuthNeededHandler(GetCourseFiles),
}

var GetCourseFiles = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	files, status, err := ada.Files(mux.Vars(req)["id"])
	util.JSON(rw, files, status, err)
})
