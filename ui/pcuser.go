package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/service"
)

type PCUserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Logo     string `json:"Logo"`
	Evidence string `json:"evidence"`
	Info     string `json:"info"`
}

type PCUserDetailInformation struct {
	Name         string `json:"name"`
	Account      string `json:"account"`
	RegisterTime string `json:"registerTime"`
	Email        string `json:"email`
	Logo         string `json:"logo"`
	Evidence     string `json:"evidence"`
	Info         string `json:"info"`
}

type PCUserListInformation struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	Verified     int    `json:"verified"`
	RegisterTime string `json:"register_time"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(PCUserSignUp)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		w.WriteHeader(400)
	}

	user := new(service.PCUser)
	user.Name = jsonBody.Name
	user.Email = jsonBody.Email
	user.Logo = jsonBody.Logo
	user.Evidence = jsonBody.Evidence
	user.Info = jsonBody.Info
	code, err := user.SignUp()
	w.WriteHeader(code)
}

func GetPCUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	userID := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	user := new(service.PCUser)
	user.ID = intID
	code, err := user.GetUserInformation()
	if code != 200 {
		w.WriteHeader(code)
		return
	}
	layout := "2006-01-02 15:04"
	jsonRet := PCUserDetailInformation{user.Name, user.Account, user.RegisterTime.Format(layout),
		user.Email, user.Logo, user.Evidence, user.Info}
	byteRet, err := json.Marshal(jsonRet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(byteRet)
}

func VerifyPCUserHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
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

	user := new(service.PCUser)
	user.ID = intID
	user.Verified = intVerify
	code, err := user.AduitUser(jsonBody.RefuseInfo)
	w.WriteHeader(code)
}

func GetPCUserListHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	r.ParseForm()
	stringType := r.FormValue("type")
	intType, err := strconv.Atoi(stringType)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	type PCUserList struct {
		Content []PCUserListInformation `json:"content"`
	}
	users := new(service.PCUserSlice)
	code, err := users.ListUsers(intType)
	if code != 200 {
		w.WriteHeader(code)
		return
	}
	list := PCUserList{}
	layout := "2006-01-02 15:04"
	for _, v:= range(*users) {
		tmp := PCUserListInformation {v.ID, v.Name, v.Logo, v.Verified, v.RegisterTime.Format(layout)}
		list.Content = append(list.Content, tmp)
	}
	ret, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(ret)
}
