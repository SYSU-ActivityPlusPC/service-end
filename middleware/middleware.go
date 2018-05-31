package middleware

import (
	"net/http"
	"strconv"
	"strings"

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
	// Allow upload and signup request
	if r.URL.Path == "/images" || (r.URL.Path == "/pcusers" && r.Method == "POST") {
		next(rw, r)
		return
	}
	// Read authorization from header
	role := 0
	r.Header.Del("X-Role")
	r.Header.Del("X-Account")
	auth := r.Header.Get("Authorization")
	// Check user identity
	if len(auth) <= 0 {
		r.Header.Set("X-Role", "0")
	} else {
		ok, name := controller.CheckToken(auth)
		if ok != 2 {
			r.Header.Set("X-Role", "0")
		} else {
			// Check user account
			user := model.GetUserInfo(name)
			if user.ID <= 0 {
				r.Header.Set("X-Role", "0")
			} else {
				isAdmin := controller.CheckIsAdmin(name)
				if isAdmin {
					role = 2
					r.Header.Set("X-Role", "2")
				} else {
					role = 1
					r.Header.Set("X-Role", "1")
				}
				r.Header.Set("X-Account", user.Account)
				r.Header.Set("X-ID", strconv.Itoa(user.ID))
			}
		}
	}
	// Handle according to different url
	path := r.URL.Path
	// Pass /session path
	if path != "/session" {
		// Allow anonymous user add activity
		if path == "/act" && r.Method == "POST" {
			next(rw, r)
			return
		}
		// Refuse all the anyous user
		if role == 0 {
			rw.WriteHeader(401)
			return
		}
		// Refuse some of the party user request
		if role == 1 {
			if path == "/act" && (r.Method == "GET" || r.Method == "PUT") {
				rw.WriteHeader(401)
				return
			}
			if strings.HasPrefix(path, "/act/") && r.Method == "DELETE" {
				rw.WriteHeader(401)
				return
			}
			if strings.HasPrefix(path, "/message") {
				rw.WriteHeader(401)
				return
			} 
		}
	}
	next(rw, r)
}
