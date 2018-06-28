package dao

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

func (user *PCUser) GetByAccount(session *xorm.Session) {
	account := user.Account
	_, err := session.Where("account = ?", account).Get(&user)
	if err != nil {
		fmt.Println(err)
		user.ID = -1
	}
}

func (user *PCUser) GetByEmail(session *xorm.Session) {
	email := user.Email
	_, err := session.Where("email=?", email).Get(user)
	if err != nil {
		fmt.Println(err)
		user = nil
	}
}

func (user *PCUser) GetByID(session *xorm.Session) {
	id := user.ID
	_, err := session.Where("id=?", id).Get(user)
	if err != nil {
		fmt.Println(err)
		user = nil
	}
}

func (user *PCUser) UpdateVerifiedStatus(session *xorm.Session) (int, error) {
	id := user.ID
	affected, err := session.Where("id=?", id).Cols("verified").Cols("account").Cols("password").Cols("register_time").Update(user)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), err
}

func (user *PCUser) ListByVerifiedStatus(session *xorm.Session) []PCUser {
	verify := user.Verified
	ret := make([]PCUser, 0)
	err := session.Where("verified = ?", verify).Incr("id").Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

func (user *PCUser) Insert(session *xorm.Session) bool {
	_, err := session.InsertOne(&user)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
