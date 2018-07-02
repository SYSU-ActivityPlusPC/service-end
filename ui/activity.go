package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/service"
)

type AddActivityMessage struct {
	service.ActivityInfo
}

type ActivityStatus struct {
	AduitNumber   int `json:"aduitNum"`
	OngoingNumber int `json:"ongoingNum"`
	OverNumber    int `json:"overNum"`
}

// AddActivityHandler add activity to the db
func AddActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	stringID := r.Header.Get("X-ID")
	// Handle anonymous user
	if len(stringID) == 0 {
		stringID = "-1"
	}
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	act := new(AddActivityMessage)
	err = json.Unmarshal(body, act)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	act.PCUserID = ID
	code, err := act.AddActivity()
	w.WriteHeader(code)
}

// ModifyActivityHandler change the message except id and verified
func ModifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	actid := mux.Vars(r)["actId"]
	intActID, err := strconv.Atoi(actid)
	if err != nil {
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

	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	act := new(AddActivityMessage)
	err = json.Unmarshal(body, act)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	code, err := act.ModifyActivity()
	w.WriteHeader(code)
}

// DeleteActivityHandler remove activity with given id
func DeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	actid := mux.Vars(r)["actId"]
	intActID, err := strconv.Atoi(actid)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	act := new(AddActivityMessage)
	act.ID = intActID
	code, err := act.DeleteActivity()
	w.WriteHeader(code)
}

// VerifyActivityHandler change the verify status of an activity
func VerifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	r.ParseForm()
	actid := r.FormValue("act")
	verified := r.FormValue("verify")
	intActID, err := strconv.Atoi(actid)
	intVerify, err := strconv.Atoi(verified)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	act := new(AddActivityMessage)
	act.ID = intActID
	act.Verified = intVerify
	code, err := act.AduitActivity()
	w.WriteHeader(code)
}

func GetNumberOfActStatusByClubHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	clubId := mux.Vars(r)["clubId"]
	intClubId, err := strconv.Atoi(clubId)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	if id != intClubId {
		w.WriteHeader(401)
		return
	}

	act := new(AddActivityMessage)
	act.PCUserID = intClubId
	status, err := act.GetActivityNumber()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	res := ActivityStatus{status[0], status[1], status[2]}
	byteRes, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(byteRes)
}

// ShowActivitiesListByClubHandler display activity with given page number, only club use
func ShowActivitiesListByClubHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	clubId := mux.Vars(r)["clubId"]
	intClubId, err := strconv.Atoi(clubId)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	if id != intClubId {
		w.WriteHeader(401)
		return
	}
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber string
	if len(r.Form["page"]) > 0 {
		pageNumber = r.Form["page"][0]
	} else {
		pageNumber = "1"
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
}

// ShowActivitiesListHandler display activity with given page number and verify status
func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber string
	if len(r.Form["page"]) > 0 {
		pageNumber = r.Form["page"][0]
	} else {
		pageNumber = "1"
	}
	var verified string
	verified = r.Form["verify"][0]
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	intVerified, err := strconv.Atoi(verified)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
}

// ShowActivityDetailHandler return required activity details with given activity id
func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	vars := mux.Vars(r)
	actid := vars["id"]
	intActID, err := strconv.Atoi(actid)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
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

func CloseActivityHandler(w http.ResponseWriter, r *http.Request) {
	role, id, account, _ := GetHeaderMessage(r)
	actid := mux.Vars(r)["actid"]
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
