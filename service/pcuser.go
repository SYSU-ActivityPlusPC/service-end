package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sysu-activitypluspc/service-end/dao"
	"github.com/sysu-activitypluspc/service-end/types"
)

type PCUser struct {
	dao.PCUser
}

type PCUserSlice []PCUser

// GetUserInformation returns user information with id
func (u *PCUser) GetUserInformation() (int, error) {
	session := GetSession()
	defer DeleteSession(session, true)
	has, err := u.GetByID(session)
	if err != nil {
		return 500, err
	}
	if !has {
		return 204, nil
	}
	return 200, nil
}

// AduitUser aduit user with id and verified status
func (u *PCUser) AduitUser(message string) (int, error) {
	if u.Verified != 1 && u.Verified != 2 {
		return 400, errors.New("Invalid verified status")
	}
	if u.ID <= 0 {
		return 400, errors.New("Invalid user message")
	}
	// Get user message
	session := GetSession()
	defer DeleteSession(session, true)
	verified := u.Verified
	has, err := u.GetByID(session)
	if err != nil {
		return 500, err
	}
	if !has {
		return 204, nil
	}
	if u.Verified == verified {
		return http.StatusNotModified, nil
	}
	u.Verified = verified

	var password string
	if u.Verified == 1 {
		now := time.Now().Add(time.Hour * 8)
		password = GeneratePassword(12)
		u.Password = getPassword(strconv.Itoa(u.ID), password)
		u.RegisterTime = &now
	} else {
		u.Password = ""
		u.Account = ""
	}
	// Update db
	affected, err := u.UpdateVerifiedStatus(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}

	// Send email
	type RejectMsg struct {
		RefuseInfo string
	}
	subject := "中大活动: 恭喜，您的账号注册请求被已通过"
	if u.Verified == 2 {
		subject = "中大活动: 很抱歉，您的账号注册请求未被通过"
	}
	content := fmt.Sprintf("您的账户由于%s<br />, 未通过审核", message)
	if u.Verified == 1 {
		content = fmt.Sprintf("您的登录账户信息为: %s<br />您的登录密码为: %s<br />感谢您使用中大活动", u.Email, password)
	}
	msg := types.EmailContent{"admin@sysuactivity.com", u.Email, subject, content}
	go SendMail(msg.From, msg.To, msg.Content, msg.Subject)
	return 200, nil
}

// ListUsers list users with given user type
func (us PCUserSlice) ListUsers(userType int) (int, error) {
	if userType != 1 && userType != 0 {
		return 400, errors.New("Invalid user type")
	}
	session := GetSession()
	u := new(dao.PCUser)
	defer DeleteSession(session, true)
	if userType == 0 {
		// Get only the verified user
		u.Verified = 1
		daoUsers, err := u.ListByVerifiedStatus(session)
		if err != nil {
			return 500, err
		}
		if len(daoUsers) == 0 {
			return 204, nil
		}
		for _, v := range daoUsers {
			tmp := PCUser{v}
			us = append(us, tmp)
		}
		return 200, nil
	}
	// Get all of the useru.Verified = 1
	u.Verified = 0
	daoUsers, err := u.ListByVerifiedStatus(session)
	if err != nil {
		return 500, err
	}
	for _, v := range daoUsers {
		tmp := PCUser{v}
		us = append(us, tmp)
	}
	u.Verified = 2
	daoUsers, err = u.ListByVerifiedStatus(session)
	if err != nil {
		return 500, err
	}
	for _, v := range daoUsers {
		tmp := PCUser{v}
		us = append(us, tmp)
	}
	if len(us) == 0 {
		return 204, nil
	}
	return 200, nil
}

// SignUp sign up a user
func (u PCUser) SignUp() (int, error) {
	// TODO: verify user data
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := u.Insert(session)
	if err != nil {
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, nil
}
