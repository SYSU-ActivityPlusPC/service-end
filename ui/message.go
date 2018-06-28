package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// AddMessageHandler add message to the db
func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
}

// DeleteMessageHandler remove message with given id
func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	id := mux.Vars(r)["id"]
	intId, err := strconv.Atoi(id)
	if err != nil || intId <= 0 {
		w.WriteHeader(400)
		return
	}
}

// ShowMessagesListHandler display message with given page number
func ShowMessagesListHandler(w http.ResponseWriter, r *http.Request) {
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

	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
}

// ShowMessageDetailHander return required message details with given message id
func ShowMessageDetailHandler(w http.ResponseWriter, r *http.Request) {
	role, _, _, _ := GetHeaderMessage(r)
	if role != 2 {
		w.WriteHeader(401)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
}
