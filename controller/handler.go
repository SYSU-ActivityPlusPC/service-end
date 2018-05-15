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

	jsonBody := new(types.StringActivityInfo)
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

func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber string
	if len(r.Form["pageNum"]) > 0 {
		pageNumber = r.Form["pageNum"][0]
	} else {
		pageNumber = "1"
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intPageNum > 0 {
		// Get activity list
		activityList := model.GetActivityList(intPageNum - 1)

		// Change each element to the format that we need
		infoArr := make([]types.ActivityIntroduction, 0)
		for i := 0; i < len(activityList); i++ {
			tmp := types.ActivityIntroduction{
				ID:        activityList[i].ID,
				Name:      activityList[i].Name,
				StartTime: activityList[i].StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:   activityList[i].EndTime.UnixNano() / int64(time.Millisecond),
				Campus:    activityList[i].Campus,
				Type:      activityList[i].Type,
				Poster:    activityList[i].Poster,
				Location:  activityList[i].Location,
				Verified:  activityList[i].Verified,
			}
			tmp.Poster = GetPoster(tmp.Poster, tmp.Type)
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
