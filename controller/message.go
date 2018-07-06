package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sysu-activitypluspc/service-end/model"
	"github.com/sysu-activitypluspc/service-end/types"
)

// AddMessageHandler add message to the db
func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)

	// Only manager can add message
	// stringID := r.Header.Get("X-ID")
	// ID, err := strconv.Atoi(stringID)
	// if err != nil {
	// 	w.WriteHeader(500)
	// 	return
	// }

	jsonBody := new(types.MessageInfo)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		w.WriteHeader(400)
	}
	_, err = model.AddMessage(*jsonBody)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// DeleteMessageHandler remove message with given id
func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) <= 0 {
		w.WriteHeader(400)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil || intId <= 0 {
		w.WriteHeader(400)
		return
	}

	_, err = model.DeleteMessage(intId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// ShowMessagesListHandler display message with given page number
func ShowMessagesListHandler(w http.ResponseWriter, r *http.Request) {
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

	// Judge if the passed param is valid
	if intPageNum > 0 {
		// Get message list
		messageList := model.GetMessageList(intPageNum - 1)

		// Find SendTo club
		// Change each element to the format that we need
		infoArr := make([]types.MessageIntroduction, 0)
		for i := 0; i < len(messageList); i++ {
			sendToList, err := model.GetMessageSendToList(messageList[i].ID)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				w.WriteHeader(500)
				return
			}

			layout := "2006-01-02 15:04"
			tmp := types.MessageIntroduction{
				ID:      messageList[i].ID,
				Subject: messageList[i].Subject,
				Body:    messageList[i].Body,
				PubTime: messageList[i].PubTime.Format(layout),
				SendTo:  sendToList,
			}
			infoArr = append(infoArr, tmp)
		}
		returnList := types.MessageList{
			Content: infoArr,
		}

		// Transfer it to json
		stringList, err := json.Marshal(returnList)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(500)
			return
		}
		if len(messageList) <= 0 {
			w.WriteHeader(204)
		} else {
			w.Write(stringList)
		}
	} else {
		w.WriteHeader(400)
	}
}

// ShowMessageDetailHander return required message details with given message id
func ShowMessageDetailHandler(w http.ResponseWriter, r *http.Request) {
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
		ok, messageInfo := model.GetMessageInfo(intID)
		if ok {
			sendToList, err := model.GetMessageSendToList(intID)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				w.WriteHeader(500)
				return
			}

			layout := "2006-01-02 15:04"
			// Convert to ms
			retMsg := types.MessageIntroduction{
				ID:      messageInfo.ID,
				Subject: messageInfo.Subject,
				Body:    messageInfo.Body,
				PubTime: messageInfo.PubTime.Format(layout),
				SendTo:  sendToList,
			}
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