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
)

// AddActivityHandler add activity to the db
func AddActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.StringActivityInfo)
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

	jsonBody := new(types.StringActivityInfo)
	json.Unmarshal(body, jsonBody)
	num, err := model.UpdateActivity(intActID, *jsonBody)
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

	num, err := model.DeleteActivity(intActID)
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

	num, err := model.VerifyActivity(intActID, intVerify)
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
			// Convert to ms
			retMsg := types.IntActivityInfo{
				ID:              activityInfo.ID,
				Name:            activityInfo.Name,
				StartTime:       activityInfo.StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:         activityInfo.EndTime.UnixNano() / int64(time.Millisecond),
				Campus:          activityInfo.Campus,
				Location:        activityInfo.Location,
				EnrollCondition: activityInfo.EnrollCondition,
				Sponsor:         activityInfo.Sponsor,
				Type:            activityInfo.Type,
				PubStartTime:    activityInfo.PubStartTime.UnixNano() / int64(time.Millisecond),
				PubEndTime:      activityInfo.PubEndTime.UnixNano() / int64(time.Millisecond),
				Detail:          activityInfo.Detail,
				Reward:          activityInfo.Reward,
				Introduction:    activityInfo.Introduction,
				Requirement:     activityInfo.Requirement,
				Poster:          activityInfo.Poster,
				Qrcode:          activityInfo.Qrcode,
				Email:           activityInfo.Email,
				Verified:        activityInfo.Verified,
			}
			retMsg.Poster = GetPoster(retMsg.Poster, retMsg.Type)
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
