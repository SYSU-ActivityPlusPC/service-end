package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/types"
)

type PCUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Logo     string `json:"Logo"`
	Evidence string `json:"evidence"`
	Info     string `json:"info"`
}

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
}
