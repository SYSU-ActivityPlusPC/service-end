package ui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sysu-activitypluspc/service-end/service"
)

func ListActivityApplyHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	r.ParseForm()
	actid := r.FormValue("act")

	intActID, err := strconv.Atoi(actid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// Validate
	if role == 0 {
		w.WriteHeader(401)
		return
	}
	if role == 1 {
		act := new(service.ActivityInfo)
		act.PCUserID = id
		act.ID = intActID
		ok, err := act.CheckMessageCorrectness()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if !ok {
			w.WriteHeader(401)
		}
	}
}

func DeleteActivityApplyHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	r.ParseForm()
	actid := r.FormValue("act")
	applyid := r.FormValue("actApply")
	intActID, err := strconv.Atoi(actid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	intApplyID, err := strconv.Atoi(applyid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// Validate
	if role == 0 {
		w.WriteHeader(401)
		return
	}
	if role == 1 {
		act := new(service.ActivityInfo)
		act.PCUserID = id
		act.ID = intActID
		ok, err := act.CheckMessageCorrectness()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if !ok {
			w.WriteHeader(401)
		}
	}
}
