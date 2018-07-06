package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/sysu-activitypluspc/service-end/model"
	"github.com/sysu-activitypluspc/service-end/types"
	"github.com/sysu-activitypluspc/service-end/services"
)

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
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.ActivityInfo)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		w.WriteHeader(400)
	}
	_, err = model.AddActivity(*jsonBody, ID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// ModifyActivityHandler change the message except id and verified
func ModifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actId"]
	if len(actid) <= 0 {
		w.WriteHeader(400)
		return
	}
	intActID, err := strconv.Atoi(actid)
	if err != nil || intActID <= 0 {
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
	_, err = model.UpdateActivity(intActID, *jsonBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// DeleteActivityHandler remove activity with given id
func DeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actId"]
	if len(actid) <= 0 {
		w.WriteHeader(400)
		return
	}
	intActID, err := strconv.Atoi(actid)
	if err != nil || intActID <= 0 {
		w.WriteHeader(400)
		return
	}

	_, err = model.DeleteActivity(intActID)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// VerifyActivityHandler change the verify status of an activity
func VerifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	actid := r.FormValue("act")
	verified := r.FormValue("verify")
	if len(actid)*len(verified) == 0 {
		w.WriteHeader(400)
		return
	}
	intActID, err := strconv.Atoi(actid)
	intVerify, err := strconv.Atoi(verified)
	if err != nil || intActID <= 0 || (intVerify != 0 && intVerify != 1) {
		w.WriteHeader(400)
		return
	}

	_, err = model.VerifyActivity(intActID, intVerify)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}

func GetNumberOfActStatusByClub(w http.ResponseWriter, r *http.Request) {
	clubId := mux.Vars(r)["clubId"]
	intClubId, err := strconv.Atoi(clubId)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	auditNum, ongoingNum, overNum := model.GetActStatusNumByClub(intClubId)

	tmp := types.NumOfActStatus{auditNum, ongoingNum, overNum}
	
	stringList, err := json.Marshal(tmp)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(500)
		return
	}
	w.Write(stringList)
}

// ShowActivitiesListByClubHandler display activity with given page number, only club use
func ShowActivitiesListByClubHandler(w http.ResponseWriter, r *http.Request) {
	clubId := mux.Vars(r)["clubId"]
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
	intClubId, err := strconv.Atoi(clubId)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intPageNum > 0 {
		// Get club's all activities
		activityList := model.GetActivityListByClub(intPageNum-1, intClubId)

		// Change each element to the format that we need
		infoArr := make([]types.ActivityIntroductionForClub, 0)
		for i := 0; i < len(activityList); i++ {
			var actType int
			now := time.Now().Add(time.Hour * 8)
			layout := "2006-01-02 15:04"

			if activityList[i].Verified == 2 {
				continue
			} else if activityList[i].Verified == 0 {
				actType = 0
			} else if activityList[i].Verified == 1 && activityList[i].PubEndTime.After(now) {
				actType = 1
			} else { 
				actType = 2
			}

			// pageViews := GetPageViewsByActId(activityList[i].ID)
			pageViews := 0
			registerNum := 0
			// activity registration is on or has done
			if activityList[i].CanEnrolled == 1 || activityList[i].CanEnrolled == 2 {
				registerNum = model.GetRegisterNumByActId(activityList[i].ID)
			}
			
			tmp := types.ActivityIntroductionForClub{
				ID:              	activityList[i].ID,
				Name:            	activityList[i].Name,
				PubStartTime:       activityList[i].PubStartTime.Format(layout),
				PageViews:          pageViews,
				RegisterNum:        registerNum,
				Type: 				actType,
				CanEnrolled:        activityList[i].CanEnrolled,
			}
			infoArr = append(infoArr, tmp)
		}
		returnList := types.ActivityListForClub{
			Content: infoArr,
		}

		// Transfer it to json
		stringList, err := json.Marshal(returnList)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(500)
			return
		}
		if len(activityList) <= 0 {
			w.WriteHeader(204)
		} else {
			w.Write(stringList)
		}
	} else {
		w.WriteHeader(400)
	}
}

// ShowActivitiesListHandler display activity with given page number and verify status
func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber string
	if len(r.Form["page"]) > 0 {
		pageNumber = r.Form["page"][0]
	} else {
		pageNumber = "1"
	}
	var verified string
	if len(r.Form["verify"]) > 0 {
		verified = r.Form["verify"][0]
	} else {
		w.WriteHeader(400)
		return
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	intVerified, err := strconv.Atoi(verified)
	if err != nil || (intVerified != 0 && intVerified != 1) {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intPageNum > 0 {
		// Get activity list
		activityList := model.GetActivityList(intPageNum-1, intVerified)

		// Change each element to the format that we need
		infoArr := make([]types.ActivityIntroduction, 0)
		for i := 0; i < len(activityList); i++ {
			tmp := types.ActivityIntroduction{
				ID:              activityList[i].ID,
				Name:            activityList[i].Name,
				StartTime:       activityList[i].StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:         activityList[i].EndTime.UnixNano() / int64(time.Millisecond),
				Campus:          activityList[i].Campus,
				EnrollCondition: activityList[i].EnrollCondition,
				Sponsor:         activityList[i].Sponsor,
				PubStartTime:    activityList[i].PubStartTime.UnixNano() / int64(time.Millisecond),
				PubEndTime:      activityList[i].PubEndTime.UnixNano() / int64(time.Millisecond),
				Verified:        activityList[i].Verified,
				Type:            activityList[i].Type,
			}
			infoArr = append(infoArr, tmp)
		}
		returnList := types.ActivityList{
			Content: infoArr,
		}

		// Transfer it to json
		stringList, err := json.Marshal(returnList)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(500)
			return
		}
		if len(activityList) <= 0 {
			w.WriteHeader(204)
		} else {
			w.Write(stringList)
		}
	} else {
		w.WriteHeader(400)
	}
}

// ShowActivityDetailHandler return required activity details with given activity id
func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intID > 0 {
		ok, activityInfo := model.GetActivityInfo(intID)
		if ok {
			layout := "2006-01-02 15:04"
			// Convert to ms
			retMsg := types.ActivityInfo{
				ID:              activityInfo.ID,
				Name:            activityInfo.Name,
				StartTime:       activityInfo.StartTime.Format(layout),
				EndTime:         activityInfo.EndTime.Format(layout),
				Campus:          activityInfo.Campus,
				Location:        activityInfo.Location,
				EnrollCondition: activityInfo.EnrollCondition,
				Sponsor:         activityInfo.Sponsor,
				Type:            activityInfo.Type,
				PubStartTime:    activityInfo.PubStartTime.Format(layout),
				PubEndTime:      activityInfo.PubEndTime.Format(layout),
				Detail:          activityInfo.Detail,
				Reward:          activityInfo.Reward,
				Introduction:    activityInfo.Introduction,
				Requirement:     activityInfo.Requirement,
				Poster:          activityInfo.Poster,
				Qrcode:          activityInfo.Qrcode,
				Email:           activityInfo.Email,
				Verified:        activityInfo.Verified,
			}
			retMsg.Poster = services.GetPoster(retMsg.Poster, retMsg.Type)
			stringInfo, err := json.Marshal(retMsg)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				w.WriteHeader(500)
				return
			}
			w.Write(stringInfo)
		} else {
			w.WriteHeader(204)
		}
	} else {
		w.WriteHeader(400)
	}
}