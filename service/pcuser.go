package service

import (
	"github.com/sysu-activitypluspc/service-end/dao"
)

type PCUser struct {
	dao.PCUser
}

type PCUserSlice []PCUser

// GetUserInformation returns user information with id
func (u *PCUser) GetUserInformation() {}

// AduitUser aduit user with id and verified status
func ( *PCUser) AduitUser() {}

// ListUsers list users with given user type
func (us PCUserSlice) ListUsers(userType int) {}

// SignUp sign up a user
func (u PCUser) SignUp() {

}