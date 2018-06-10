package middleware

import (
	"github.com/gorilla/mux"
	"os"
	"fmt"
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
	userId := 0
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
				userId = -1
			} else {
				isAdmin := controller.CheckIsAdmin(name)
				userId = user.ID
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
			if strings.HasPrefix(path, "/act/{clubId}/list") {
				// judge whether clubId and token match
				clubId := mux.Vars(r)["clubId"]
				intClubId, err := strconv.Atoi(clubId)
				if err != nil {
					fmt.Fprint(os.Stderr, err)
					rw.WriteHeader(400)
					return
				}

				if userId == intClubId {
					next(rw, r)
				} else {
					rw.WriteHeader(401)
					return
				}
			}
		}
		// Refuse all the pcusers api
		if path == "/pcusers" || strings.HasPrefix(path, "/pcusers/") {
			if role != 2 {
				rw.WriteHeader(401)
				return
			}
		}
		// Refuse super manager to access this url
		// if role == 2 {
		// 	if strings.HasPrefix(path, "/act/club") {
		// 		rw.WriteHeader(401)
		// 		return
		// 	}
		// }
	}
	next(rw, r)
}
