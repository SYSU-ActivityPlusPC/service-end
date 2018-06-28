package dao

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

func (msg *Message) Insert(session *xorm.Session) (int, error) {
	affected, err := session.InsertOne(msg)
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

func (msg *Message) List(session *xorm.Session, pageNum int) []Message {
	messageList := make([]Message, 0)
	session.Desc("id").Find(&messageList)
	from := pageNum * 10
	if from >= len(messageList) {
		return []Message{}
	}
	if from+10 > len(messageList) {
		return messageList[from:]
	}
	return messageList[from : from+10]
}

func (msg *Message) Get(session *xorm.Session) error {
	id := msg.ID
	_, err := session.ID(id).Get(&msg)
	if err != nil {
		msg = nil
		return err
	}
	return nil
}

func (msg *Message) Delete(session *xorm.Session) (int, error) {
	id := msg.ID
	affected, err := session.Id(id).Delete(&Message{})
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		fmt.Println("Failed to delete an message")
	}
	return int(affected), nil
}

// MessagePCUser
func (msg *MessagePCUser) ListSentTo(session *xorm.Session) ([]string, error) {
	messageId := msg.MessageId
	messagePCUserList := make([]MessagePCUser, 0)
	err := session.Where("message_id=?", messageId).Find(&messagePCUserList)
	if err != nil {
		return []string{}, err
	}

	clubNameList := make([]string, 0)
	for _, v := range messagePCUserList {
		var clubName string
		_, err := session.Table("pcuser").Where("id=?", v.PCUserId).Cols("name").Get(&clubName)
		if err != nil {
			return []string{}, err
		}
		clubNameList = append(clubNameList, clubName)
	}
	return clubNameList, nil
}

func (msg *MessagePCUser) Insert(session *xorm.Session) int {
	affected, err := session.InsertOne(msg)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return int(affected)
}
