package model

import (
	"fmt"
	"time"

	"github.com/sysu-activitypluspc/service-end/types"
)

// CheckPCUser check if the account exists
func CheckPCUser(username string) bool {
	has, _ := Engine.Where("account = ?", username).Exist(&PCUser{})
	return has
}

// GetUserInfo returns user information with given username
func GetUserInfo(account string) PCUser {
	var user PCUser
	_, err := Engine.Where("account = ?", account).Get(&user)
	if err != nil {
		fmt.Println(err)
		user.ID = -1
	}
	return user
}

// AddUser add user
func AddUser(user types.PCUserSignInfo) bool {
	currentTime := time.Now()
	dbuser := PCUser{
		Name:         user.Name,
		Email:        user.Email,
		Logo:         user.Logo,
		Evidence:     user.Evidence,
		Info:         user.Info,
		RegisterTime: &currentTime,
	}
	_, err := Engine.InsertOne(&dbuser)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetUserByEmail return user detail based on email
func GetUserByEmail(email string) *PCUser {
	user := new(PCUser)
	_, err := Engine.Where("email=?", email).Get(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

// GetUserByID return user detail based on id
func GetUserByID(id int) *PCUser {
	user := new(PCUser)
	ok, err := Engine.Where("id=?", id).Get(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !ok {
		user.ID = -1
	}
	return user
}

// VerifyUser set user verified status
func VerifyUser(id int, verify int, email string, password string, currentTime *time.Time) error {
	user := new(PCUser)
	user.Verified = verify
	user.Account = email
	user.Password = password
	user.RegisterTime = currentTime
	_, err := Engine.Where("id=?", id).Cols("verified").Cols("account").Cols("password").Cols("register_time").Update(user)
	return err
}

// GetUserList get all the user with given status
func GetUserList(verify int) []PCUser {
	ret := make([]PCUser, 0)
	err := Engine.Where("verified = ?", verify).Incr("id").Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}