package middleware

import (
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/sysu-activitypluspc/service-end/service"
	"github.com/sysu-activitypluspc/service-end/dao"
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
		ok, name := service.CheckToken(auth)
		if ok != 2 {
			r.Header.Set("X-Role", "0")
		} else {
			// Check user account
			user := dao.GetUserInfo(name)
			if user.ID <= 0 {
				r.Header.Set("X-Role", "0")
				userId = -1
			} else {
				isAdmin := service.CheckIsAdmin(name)
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
		// Refuse all the anonymous user
		if role == 0 {
			rw.WriteHeader(401)
			return
		}
		// Refuse some of the party user request
		if role == 1 {
			if path == "/act" && r.Method == "GET" {
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
			// for "/act/{clubId}/list" or "/act/{clubId}/status"
			if strings.HasPrefix(path, "/act/") && (strings.HasSuffix(path, "/list") || strings.HasSuffix(path, "/status")) {
				// judge whether clubId and token match
				pathArr := strings.Split(path, "/")
				clubId := pathArr[2]
				intClubId, err := strconv.Atoi(clubId)
				if err != nil {
					fmt.Fprint(os.Stderr, err)
					rw.WriteHeader(400)
					return
				}

				if userId == intClubId {
					next(rw, r)
					return
				} else {
					rw.WriteHeader(401)
					return
				}
			}
			
			if strings.HasPrefix(path, "/act/") && r.Method == "POST" {
				pathArr := strings.Split(path, "/")
				actId := pathArr[2]
				intActId, err := strconv.Atoi(actId)
				if err != nil {
					fmt.Fprint(os.Stderr, err)
					rw.WriteHeader(400)
					return
				}

				yes, _ := dao.IsPublishedByClub(userId, intActId)
				if yes {
					next(rw, r)
					return
				} else {
					rw.WriteHeader(401)
					return
				}
			}

			if strings.HasPrefix(path, "/act/") && r.Method == "PUT" {
				actid := mux.Vars(r)["actid"]
				intActID, err := strconv.Atoi(actid)
				if intActID <= 0 || err != nil {
					rw.WriteHeader(400)
					return
				}
				if ok, err := dao.IsPublishedByClub(userId, intActID); err != nil || !ok {
					rw.WriteHeader(401)
					return
				}
			}

			if path == "actApply" {
				r.ParseForm()
				actID := r.FormValue("act")
				intActID, err := strconv.Atoi(actID)
				if intActID <= 0 || err != nil {
					rw.WriteHeader(400)
					return
				}
				if ok, err := dao.IsPublishedByClub(userId, intActID); err != nil || !ok {
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
	}
	next(rw, r)
}
