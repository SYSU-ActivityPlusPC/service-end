package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/types"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.PCUserSignInfo)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		w.WriteHeader(400)
	}
	if isEmailExist := CheckIfEmailExist(jsonBody.Email); isEmailExist == true {
		w.WriteHeader(400)
		return
	}
	ok := dao.AddUser(*jsonBody)
	if ok {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(400)
}

func GetPCUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	user := dao.GetUserByID(intID)
	if user.ID == -1 {
		w.WriteHeader(204)
	}
	var stringTime string
	if user.RegisterTime != nil {
		layout := "2006-01-02 15:04"
		stringTime = user.RegisterTime.Format(layout)
	}
	jsonRet := types.PCUserDetailedInfo{user.Name, user.Email, user.Logo, user.Evidence, user.Info, user.Account, stringTime}
	byteRet, err := json.Marshal(jsonRet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(byteRet)
}

func VerifyPCUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	verify := r.FormValue("verify")
	intID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	intVerify, err := strconv.Atoi(verify)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if intVerify != 2 && intVerify != 1 {
		w.WriteHeader(400)
		return
	}
	// Get body message
	refuseMessage := ""
	type RejectMsg struct {
		RefuseInfo string
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	jsonBody := new(RejectMsg)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		fmt.Println(err)
	} else {
		refuseMessage = jsonBody.RefuseInfo
	}

	// Update db, including time, password, account and status
	user := dao.GetUserByID(intID)
	if user.ID == -1 {
		w.WriteHeader(204)
		return
	}
	if user.Verified == intVerify {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	var password string
	if intVerify == 1 {
		now := time.Now().Add(time.Hour * 8)
		password = GeneratePassword(12)
		err = dao.VerifyUser(intID, intVerify, user.Email, getPassword(strconv.Itoa(user.ID), password), &now)
	} else {
		err = dao.VerifyUser(intID, intVerify, "", "", nil)
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(204)
		return
	}
	w.WriteHeader(200)
	// Send message to the user
	subject := "中大活动: 恭喜，您的账号注册请求被已通过"
	if intVerify == 2 {
		subject = "中大活动: 很抱歉，您的账号注册请求未被通过"
	}
	content := fmt.Sprintf("您的账户由于%s<br />导致未通过审核", refuseMessage)
	if intVerify == 1 {
		content = fmt.Sprintf("您的登录账户信息为: %s<br />您的登录密码为: %s<br />感谢您使用中大活动", user.Email, password)
	}
	msg := types.EmailContent{"admin@sysuactivity.com", user.Email, subject, content}
	go SendMail(msg.From, msg.To, msg.Content, msg.Subject)
	// byteContent, err := json.Marshal(&msg)
	// if err != nil {
	// 	fmt.Println(err)
	// 	w.WriteHeader(500)
	// 	return
	// }
	// WriteMessageQueue("email", byteContent)
}

func GetPCUserListHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stringType := r.FormValue("type")
	intType, err := strconv.Atoi(stringType)
	if err != nil || (intType != 0 && intType != 1) {
		w.WriteHeader(400)
		return
	}

	type pcuserList struct {
		Content []types.PCUserListInfo `json:"content"`
	}

	// Get user list from the db
	retContent := make([]types.PCUserListInfo, 0)
	if intType == 0 {
		// Get only the verified user
		userList := dao.GetUserList(1)
		if userList == nil {
			w.WriteHeader(500)
			return
		}
		if len(userList) == 0 {
			w.WriteHeader(204)
			return
		}
		for _, v := range userList {
			var stringTime string
			if v.RegisterTime != nil {
				layout := "2006-01-02 15:04"
				stringTime = v.RegisterTime.Format(layout)
			}
			tmp := types.PCUserListInfo{v.ID, v.Name, v.Logo, v.Verified, stringTime}
			retContent = append(retContent, tmp)
		}
	} else {
		// Get all of the user
		userList := dao.GetUserList(0)
		tmp := dao.GetUserList(2)
		if userList == nil || tmp == nil {
			w.WriteHeader(500)
			return
		}
		userList = append(userList, tmp...)
		if len(userList) == 0 {
			w.WriteHeader(204)
			return
		}
		for _, v := range userList {
			var stringTime string
			if v.RegisterTime != nil {
				layout := "2006-01-02 15:04"
				stringTime = v.RegisterTime.Format(layout)
			}
			tmp := types.PCUserListInfo{v.ID, v.Name, v.Logo, v.Verified, stringTime}
			retContent = append(retContent, tmp)
		}
	}
	// Write content back to the response
	jsonRet := pcuserList{retContent}
	byteRet, err := json.Marshal(&jsonRet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(byteRet)
}
