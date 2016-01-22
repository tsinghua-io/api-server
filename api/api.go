/*
Restful api for communicating with mobile app.
*/
package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (webapp *WebApp) BindRoute() {
	// 按名字bind/按照资源划分...  这个以后改
	webapp.Router.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	webapp.Router.HandleFunc("/users/{id}", postUserHandler).Methods("POST")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	userid := mux.Vars(r)["id"]

	// test data
	user := map[string]string{
		"id":         userid,
		"name":       "名字",
		"department": "院系",
		"type":       "undergraduate", // master, phd, teacher

		// Only in full type.
		"class":  "班级",    // 可能为null
		"gender": "male",  // female, unknown
		"email":  "email", // 可能为null
		"phone":  "phone number"}

	j, _ := json.Marshal(user)

	w.Write(j)
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
}
