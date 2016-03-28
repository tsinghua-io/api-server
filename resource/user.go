package resource

import (
	"github.com/tsinghua-io/api-server/adapter/cic/learn"
	"github.com/tsinghua-io/api-server/util"
	"net/http"
)

var Profile = Resource{
	"GET": util.AuthNeededHandler(GetProfile),
}

var GetProfile = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	userId, password, _ := req.BasicAuth()
	if ada, status, err := learn.New(userId, password); err != nil {
		util.Error(rw, err.Error(), status)
	} else if v, status, err := ada.Profile(); err != nil {
		util.Error(rw, err.Error(), status)
	} else {
		util.JSON(rw, v, status)
	}
})
