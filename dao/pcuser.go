package dao

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

func (user *PCUser) GetByAccount(session *xorm.Session) (bool, error) {
	account := user.Account
	has, err := session.Where("account = ?", account).Get(&user)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return has, nil
}

func (user *PCUser) GetByEmail(session *xorm.Session) (bool, error) {
	email := user.Email
	has, err := session.Where("email=?", email).Get(user)
	if err != nil {
		fmt.Println(err)
		user = nil
		return false, err
	}
	return has, nil
}

func (user *PCUser) GetByID(session *xorm.Session) (bool, error) {
	id := user.ID
	has, err := session.Where("id=?", id).Get(user)
	if err != nil {
		fmt.Println(err)
		user = nil
		return false, err
	}
	return has, nil
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

func (user *PCUser) ListByVerifiedStatus(session *xorm.Session) ([]PCUser, error) {
	verify := user.Verified
	ret := make([]PCUser, 0)
	err := session.Where("verified = ?", verify).Incr("id").Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func (user *PCUser) Insert(session *xorm.Session) (int, error) {
	affected, err := session.InsertOne(&user)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}
