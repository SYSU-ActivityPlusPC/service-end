package service

import (
	"strconv"

	"github.com/sysu-activitypluspc/service-end/dao"
)

type Session struct {
	Account  string
	Password string
}

// Login user login
func (s *Session) Login() (int, string, error) {
	session := GetSession()
	defer DeleteSession(session, true)
	u := new(dao.PCUser)
	u.Account = s.Account
	has, err := u.GetByAccount(session)
	if err != nil {
		return 500, "", err
	}
	if !has {
		return 400, "", nil
	}
	password := getPassword(strconv.Itoa(u.ID), s.Password)
	if password != u.Password {
		return 400, "", nil
	}
	token, _ := GenerateJWT(s.Account)
	return 200, token, nil
}
