package resource

import (
	"github.com/tsinghua-io/api-server/adapter/tsinghua.edu.cn/cic/learn"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var Profile = Resource{
	"GET": util.AuthNeededHandler(GetProfile),
}

var GetProfile = learn.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, ada *learn.Adapter) {
	v, status, err := ada.Profile()
	util.JSON(rw, v, status, err)
})
