package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sysu-activitypluspc/service-end/model"
	"github.com/sysu-activitypluspc/service-end/types"
)

func ListActivityApplyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	actid := r.FormValue("act")

	intActID, err := strconv.Atoi(actid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	if intActID <= 0 {
		w.WriteHeader(400)
		return
	}

	list := model.GetApplyListByID(intActID)
	if list == nil {
		w.WriteHeader(500)
		return
	}
	if len(list) == 0 {
		w.WriteHeader(204)
		return
	}
	type RetType struct {
		TableTitle []string `json:"tableTitle"`
		Content []types.ActivityApplyMessage `json:"content"`
	}
	content := make([]types.ActivityApplyMessage, 0)
	for _, v := range(list) {
		tmp := types.ActivityApplyMessage{v.ID, v.UserName, v.StudentId, v.Phone, v.School}
		content = append(content, tmp)
	}
	retMsg := RetType{[]string{"名字", "学号", "联系方式", "学院"}, content}
	ret, err := json.Marshal(retMsg)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(ret)
}

func DeleteActivityApplyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	actid := r.FormValue("act")
	applyid := r.FormValue("actApply")
	intActID, err := strconv.Atoi(actid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	if intActID <= 0 {
		w.WriteHeader(400)
		return
	}
	intApplyID, err := strconv.Atoi(applyid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	if intApplyID <= 0 {
		w.WriteHeader(400)
		return
	}
	// Check if the user can access this activity
	apply := model.GetApplyByID(intApplyID)
	if apply.ID <= 0 || apply.ActId != intActID {
		w.WriteHeader(204)
		return
	}

	isRemoved := model.DeleteApplyByID(intActID, intApplyID)
	if isRemoved {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(500)
}

func CloseActivityApply(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actid"]
	intActID, err := strconv.Atoi(actid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	if intActID <= 0 {
		w.WriteHeader(400)
		return
	}

	isClosed := model.CloseActApplyByID(intActID)
	if isClosed {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(500)
}