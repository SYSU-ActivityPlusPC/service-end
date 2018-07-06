package model

import (
	"fmt"
	"time"

	"github.com/sysu-activitypluspc/service-end/types"
)

// AddMessage insert a message and message_pcuser into the db
func AddMessage(messageInfo types.MessageInfo) (int, error) {
	// add message
	currentTime := time.Now()
	dbMessage := Message{
		Subject: messageInfo.Subject,
		Body: 	messageInfo.Body,
		PubTime: &currentTime,
	}
	_, err := Engine.InsertOne(&dbMessage)
	if err != nil {
		return -1, err
	}
	fmt.Println(dbMessage.ID)

	// add message_pcuser
	for _, v := range messageInfo.PCUserId {
		dbMessagePCUser := MessagePCUser{
			PCUserId:  v,
			MessageId: dbMessage.ID,
		}
		_, err := Engine.InsertOne(&dbMessagePCUser)
		if err != nil {
			return -1, err
		}
	}
	return 0, nil
}

// GetMessageInfo return wanted message detail information which is given by id
func GetMessageInfo(id int) (bool, Message) {
	var message Message

	ok, _ := Engine.ID(id).Get(&message)
	return ok, message
}

// GetMessageList return wanted message list with given page number
func GetMessageList(pageNum int) []Message {
	messageList := make([]Message, 0)
	Engine.Desc("id").Find(&messageList)
	from := pageNum * 10
	if from >= len(messageList) {
		return []Message{}
	}
	if from+10 > len(messageList) {
		return messageList[from:]
	}
	return messageList[from : from+10]
}

// GetMessageSendToList return []string which is the name list of sendto clubs
func GetMessageSendToList(messageId int) ([]string, error) {
	messagePCUserList := make([]MessagePCUser, 0)
	err := Engine.Where("message_id=?", messageId).Find(&messagePCUserList)
	if err != nil {
		return []string{}, err
	}

	clubNameList := make([]string, 0)
	for _, v := range messagePCUserList {
		var clubName string
		_, err := Engine.Table("pcuser").Where("id=?", v.PCUserId).Cols("name").Get(&clubName)
		if err != nil {
			return []string{}, err
		}
		clubNameList = append(clubNameList, clubName)
	}
	return clubNameList, nil
}

// DeleteMessage remove an activity according to the id
func DeleteMessage(id int) (int, error) {
	affected, err := Engine.Id(id).Delete(&Message{})
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		fmt.Println("Failed to delete an message")
	}
	return int(affected), nil
}