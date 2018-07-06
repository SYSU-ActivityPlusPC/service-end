package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sysu-activitypluspc/service-end/service"
)

type ActApplyListContent struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	StudentID string `json:"studentid"`
	School    string `json:"school"`
}

type ActApplyListResponse struct {
	TableTitle []string              `json:"tableTitle"`
	Content    []ActApplyListContent `json:"content"`
}

func ListActivityApplyHandler(w http.ResponseWriter, r *http.Request) {
	role, id, _, _ := GetHeaderMessage(r)
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

	// Get data
	var applys service.ActApplySlice
	code, _ := applys.GetApplyList(intActID)
	if code != 200 {
		w.WriteHeader(code)
		return
	}
	var res ActApplyListResponse
	for _, v := range applys {
		tmp := ActApplyListContent{v.ID, v.UserName, v.StudentId, v.School}
		res.Content = append(res.Content, tmp)
	}
	byteRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(byteRes)
}

func DeleteActivityApplyHandler(w http.ResponseWriter, r *http.Request) {
	role, id, _, _ := GetHeaderMessage(r)
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

	apply := new(service.ActApplyInfo)
	apply.ID = intApplyID
	apply.ActId = intActID
	code, _ := apply.DeleteApply()
	w.WriteHeader(code)
}
