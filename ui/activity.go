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

type ClubActivityListInformation struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PubStartTime string `json:"pubStartTime"`
	PageViews    int    `json:"pageViews"`
	RegisterNum  int    `json:"registerNum"`
	Type         int    `json:"type"`
	CanEnrolled  int    `json:"canEnrolled"`
}

type AdminActivityListInformation struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	Campus          int    `json:"campus"`
	EnrollCondition string `json:"enrollCondition"`
	Sponsor         string `json:"sponsor"`
	PubStartTime    string `json:"pubStartTime"`
	PubEndTime      string `json:"pubEndTime"`
	Verified        int    `json:"verified"`
	Type            int    `json:"type"`
}

// TODO: Add `json:""`
type ActivityDetailInformation struct {
	ID              int
	Name            string
	StartTime       string
	EndTime         string
	Location        string
	Campus          int
	EnrollCondition string
	Sponsor         string
	Type            int
	PubStartTime    string
	PubEndTime      string
	Detail          string
	EnrollWay       string
	EnrollEndTime   string
	Reward          string
	Introduction    string
	Requirement     string
	Poster          string
	Qrcode          string
	CanEnrolled     int
	Verified        int
	Email           string
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
	role, id, _, _ := GetHeaderMessage(r)
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
	role, _, _, _ := GetHeaderMessage(r)
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
	role, _, _, _ := GetHeaderMessage(r)
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
	_, id, _, _ := GetHeaderMessage(r)
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
	_, id, _, _ := GetHeaderMessage(r)
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

	type ClubActivityList struct {
		Content []AdminActivityListInformation `json:"content"`
	}
	acts := new(service.ActivitySlice)
	code, err := acts.GetActivityListByClud(intPageNum, intClubId)
	if code != 200 {
		w.WriteHeader(code)
		return
	}
	layout := "2006-01-02 15:04"
	actList := ClubActivityList{}
	for _, v := range *acts {
		tmp := AdminActivityListInformation{v.ID, v.Name, v.StartTime.Format(layout), v.EndTime.Format(layout), v.Campus, v.EnrollCondition, v.Sponsor,
			v.PubStartTime.Format(layout), v.PubEndTime.Format(layout), v.Verified, v.Type}
		actList.Content = append(actList.Content, tmp)
	}
	ret, err := json.Marshal(actList)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(ret)
}

// ShowActivitiesListHandler display activity with given page number and verify status
func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
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

	acts := new(service.ActivitySlice)
	code, err := acts.GetActivityListByAdmin(intPageNum, intVerified)
	if code != 200 {
		w.WriteHeader(code)
		return
	}
	type ClubActivityList struct {
		Content []ClubActivityListInformation `json:"content"`
	}
	layout := "2006-01-02 15:04"
	actList := ClubActivityList{}
	for _, v := range *acts {
		tmp := ClubActivityListInformation{v.ID, v.Name, v.PubStartTime.Format(layout),
			0, 0, v.Type, v.CanEnrolled}
		actList.Content = append(actList.Content, tmp)
	}
	ret, err := json.Marshal(actList)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(ret)
}

// ShowActivityDetailHandler return required activity details with given activity id
func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	role, id, _, _ := GetHeaderMessage(r)
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

	act := new(service.ActivityInfo)
	act.ID = intActID
	code, err := act.GetActivityInfo()
	if code != 200 {
		w.WriteHeader(code)
	}
	layout := "2006-01-02 15:04"
	jsonRet := ActivityDetailInformation{act.ID, act.Name, act.StartTime.Format(layout),
		act.EndTime.Format(layout), act.Location, act.Campus, act.EnrollCondition, act.Sponsor,
		act.Type, act.PubStartTime.Format(layout), act.PubEndTime.Format(layout), act.Detail,
		act.EnrollWay, act.EnrollEndTime.Format(layout), act.Reward, act.Introduction, act.Requirement,
		act.Poster, act.Qrcode, act.CanEnrolled, act.Verified, act.Email}
	ret, err := json.Marshal(jsonRet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(ret)
}

func CloseActivityHandler(w http.ResponseWriter, r *http.Request) {
	role, id, _, _ := GetHeaderMessage(r)
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

	act := new(service.ActivityInfo)
	act.ID = intActID
	code, err := act.CloseActivity()
	w.WriteHeader(code)
}
