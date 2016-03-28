package resource

import (
	"github.com/gorilla/mux"
	"net/http"
)

var User = Resource{
	"GET": GetUser,
}

func GetUser(r *http.Request) (interface{}, int) {
	vars := mux.Vars(r)
	return vars, http.StatusOK
}
