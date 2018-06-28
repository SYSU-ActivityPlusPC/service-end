package service

import (
	"github.com/sysu-activitypluspc/service-end/dao"
)

type PCUser struct {
	dao.PCUser
}

type PCUserSlice []PCUser

// GetUserInformation returns user information with id
func (u *PCUser) GetUserInformation() {
	session := GetSession()
	defer DeleteSession(session, true)
	u.GetByID(session)
}

// AduitUser aduit user with id and verified status
func (u *PCUser) AduitUser() bool {
	if u.Verified != 1 && u.Verified != 2 {
		return false
	}
	session := GetSession()
	defer DeleteSession(session, true)
	_, err := u.UpdateVerifiedStatus(session)
	if err != nil {
		DeleteSession(session, false)
		return false
	}
	// TODO: Send email
	return true
}

// ListUsers list users with given user type
func (us PCUserSlice) ListUsers(userType int) {
	session := GetSession()
	defer DeleteSession(session, true)
	// TODO: make the requirment clear and finish
}

// SignUp sign up a user
func (u PCUser) SignUp() bool {
	// TODO: verify user data
	session := GetSession()
	defer DeleteSession(session, true)
	return u.Insert(session)
}
