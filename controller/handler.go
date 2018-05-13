package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sysu-activitypluspc/service-end/model"
	"github.com/sysu-activitypluspc/service-end/types"
)

func AddActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.ActivityInfo)
	json.Unmarshal(body, jsonBody)
	num, err := model.AddActivity(*jsonBody)
	if num == 0 {
		w.WriteHeader(204)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func ModifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actId"]
	if len(actid) <= 0 {
		w.WriteHeader(400)
		return
	}
	intActId, err := strconv.Atoi(actid)
	if err != nil || intActId <= 0 {
		w.WriteHeader(400)
		return
	}

	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.ActivityInfo)
	json.Unmarshal(body, jsonBody)
	num, err := model.UpdateActivity(intActId, *jsonBody)
	if num == 0 {
		w.WriteHeader(204)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func DeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actId"]
	if len(actid) <= 0 {
		w.WriteHeader(400)
		return
	}
	intActId, err := strconv.Atoi(actid)
	if err != nil || intActId <= 0 {
		w.WriteHeader(400)
		return
	}

	num, err := model.DeleteActivity(intActId)
	if num == 0 {
		w.WriteHeader(204)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func VerifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	actid := r.FormValue("act")
	verified := r.FormValue("verify")
	if len(actid)*len(verified) == 0 {
		w.WriteHeader(400)
		return
	}
	intActId, err := strconv.Atoi(actid)
	intVerify, err := strconv.Atoi(verified)
	if err != nil || intActId <= 0 || (intVerify != 0 && intVerify != 1) {
		w.WriteHeader(400)
		return
	}

	num, err := model.VerifyActivity(intActId, intVerify)
	if num == 0 {
		w.WriteHeader(204)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}
