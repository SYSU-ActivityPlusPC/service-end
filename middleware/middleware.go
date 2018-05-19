package middleware

import (
	"net/http"

	"github.com/sysu-activitypluspc/service-end/controller"
	"github.com/sysu-activitypluspc/service-end/model"
)

// ValidUserMiddleWare check token and decide the permission
// 0->Basic permission
// 1->Activity uploader
// 2->Super user
// Timeout authorization will be set to 0
// Suppose that only manager can access the api
// User level and account name will be set
type ValidUserMiddleWare struct {
}

func (v ValidUserMiddleWare) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Read authorization from header
	r.Header.Del("X-Role")
	r.Header.Del("X-Account")
	auth := r.Header.Get("Authorization")
	if len(auth) <= 0 {
		r.Header.Set("X-Role", "0")
		next(rw, r)
		return
	}
	ok, name := controller.CheckToken(auth)
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
	isAdmin := controller.CheckIsAdmin(name)
	if isAdmin {
		r.Header.Set("X-Role", "2")
		next(rw, r)
		return
	}
	r.Header.Set("X-Role", "1")
	r.Header.Set("X-Account", user.Account)
	next(rw, r)
}
