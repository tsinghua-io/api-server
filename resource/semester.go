package resource

import (
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/x/learn"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var Semester = Resource{
	"GET": util.AuthNeededHandler(GetSemester),
}

var GetSemester = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	semester, _, status, err := ada.Semesters()
	util.JSON(rw, semester, status, err)
})
