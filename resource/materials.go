package resource

import (
	"github.com/gorilla/mux"
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/x/learn"
	"github.com/tsinghua-io/api-server/model"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var CourseMaterials = Resource{
	"GET": util.AuthNeededHandler(GetCourseMaterials),
}

var GetCourseMaterials = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	v, status, err := BatchResourceFunc(
		mux.Vars(req)["id"],
		func(id string) (interface{}, int, error) {
			materials := new(model.Materials)
			sg := util.NewStatusGroup()

			sg.Go(func(status *int, err *error) {
				materials.Announcements, *status, *err = ada.Announcements(id)
			})
			sg.Go(func(status *int, err *error) {
				materials.Files, *status, *err = ada.Files(id)
			})
			sg.Go(func(status *int, err *error) {
				materials.Assignments, *status, *err = ada.Assignments(id)
			})

			status, err := sg.Wait()
			return materials, status, err
		})
	util.JSON(rw, v, status, err)
})
