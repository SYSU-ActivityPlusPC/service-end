package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	err = model.AddActivity(*jsonBody)
	if err != nil {
		w.WriteHeader(500)
	}
}
