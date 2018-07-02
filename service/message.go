package service

import (
	"errors"
	"github.com/sysu-activitypluspc/service-end/dao"
)

type Message struct {
	dao.Message
	dao.MessagePCUser
}

type MessageSlice []Message

// PublishMessage publish a message
func (msg *Message) PublishMessage() (int, error) {
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := msg.Message.Insert(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	affected, err = msg.MessagePCUser.Insert(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, err
}

// ListMessages list messages with given page
func (msgs MessageSlice) ListMessages(page int) (int, error){
	msg := new(dao.Message)
	session := GetSession()
	defer DeleteSession(session, true)
	daoMsgs, err := msg.List(session, page)
	if err != nil {
		return 500, err
	}
	for _, v := range daoMsgs {
		tmp := Message{}
		tmp.Message = v
		msgs = append(msgs, tmp)
	}
	if len(msgs) == 0 {
		return 204, nil
	}
	return 200, nil
}

// GetMessageInformation returns message information with given message id
func (msg *Message) GetMessageInformation() (int, error){
	if msg.Message.ID <= 0 {
		return 400, errors.New("Invalid message id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	has, err := msg.Get(session)
	if err != nil {
		return 500, err
	}
	if !has {
		return 204, nil
	}
	return 200, nil
}

// DeleteMessage delete message with id
func (msg *Message) DeleteMessage() (int, error) {
	if msg.Message.ID <= 0 {
		return 400, errors.New("Invalid message id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := msg.Delete(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, nil
}
