package controller

import (
	"net/http"

	"github.com/sysu-activitypluspc/service-end/model"
)

// ValidUserMiddleWare check token and decide the permission
// 0->Basic permission
// 1->Activity uploader
// 2->Super user
// Timeout authorization will be set to 0
// Suppose that only manager can access the api
func ValidUserMiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Read authorization from header
	auth := r.Header.Get("Authorization")
	if len(auth) <= 0 {
		r.Header.Set("X-Role", "0")
		next(rw, r)
		return
	}
	ok, name := CheckToken(auth)
	if ok != 2 {
		r.Header.Set("X-Role", "0")
		next(rw, r)
		return
	}
	// Check name
	user := model.GetUserInfo(name)
	if user.ID <= 0 {
		r.Header.Set("X-Role", "0")
		next(rw, r)
		return
	}
	isAdmin := CheckIsAdmin(name)
	if isAdmin {
		r.Header.Set("X-Role", "2")
		next(rw, r)
		return
	}
	next(rw, r)
}
