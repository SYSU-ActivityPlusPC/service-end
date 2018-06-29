package middleware

import (
	"net/http"
	"strconv"

	"github.com/sysu-activitypluspc/service-end/dao"
	"github.com/sysu-activitypluspc/service-end/service"
)

// ParseToken parse token and decide the role and permission
// 0->Basic permission
// 1->Activity uploader
// 2->Super user
// User Role, ID and Account name will be set
type ParseToken int

func (v ParseToken) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Read authorization from header
	r.Header.Del("X-Role")
	r.Header.Del("X-Account")
	r.Header.Del("X-ID")
	auth := r.Header.Get("Authorization")
	r.Header.Del("Authorization")

	// Check user identity
	if len(auth) <= 0 {
		r.Header.Set("X-Role", "0")
	} else {
		ok, name := service.CheckToken(auth)
		if ok != 2 {
			r.Header.Set("X-Role", "0")
		} else {
			// Check user account
			u := new(dao.PCUser)
			u.Account = name
			session := service.GetSession()
			defer service.DeleteSession(session, true)
			has, err := u.GetByAccount(session)
			if err != nil || !has {
				r.Header.Set("X-Role", "0")
			} else {
				isAdmin := service.CheckIsAdmin(name)
				if isAdmin {
					r.Header.Set("X-Role", "2")
				} else {
					r.Header.Set("X-Role", "1")
				}
				r.Header.Set("X-Account", u.Account)
				r.Header.Set("X-ID", strconv.Itoa(u.ID))
			}
		}
	}
	next(rw, r)
}
