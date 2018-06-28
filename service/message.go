package service

import (
	"github.com/sysu-activitypluspc/service-end/dao"
)

type Message struct {
	dao.Message
	dao.MessagePCUser
}

type MessageSlice []Message

// PublishMessage publish a message
func (msg *Message) PublishMessage() bool {
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := msg.Message.Insert(session)
	if err != nil {
		DeleteSession(session, false)
		return false
	}
	affected = msg.MessagePCUser.Insert(session)
	if affected == -1 {
		DeleteSession(session, false)
		return false
	}
	return true
}

// ListMessages list messages with given page
func (msgs MessageSlice) ListMessages(page int) {
	msg := new(dao.Message)
	session := GetSession()
	defer DeleteSession(session, true)
	daoMsgs := msg.List(session, page)
	if daoMsgs == nil {
		msgs = nil
		return
	}
	for _, v := range daoMsgs {
		tmp := Message{}
		tmp.Message = v
		msgs = append(msgs, tmp)
	}
}

// GetMessageInformation returns message information with given message id
func (msg *Message) GetMessageInformation() {
	if msg.Message.ID <= 0 {
		msg = nil
		return
	}
	session := GetSession()
	defer DeleteSession(session, true)
	err := msg.Get(session)
	if err != nil {
		msg = nil
	}
}

// DeleteMessage delete message with id
func (msg *Message) DeleteMessage() bool {
	if msg.Message.ID <= 0 {
		msg = nil
		return false
	}
	session := GetSession()
	defer DeleteSession(session, true)
	_, err := msg.Delete(session)
	if err != nil {
		DeleteSession(session, false)
		return false
	}
	return true
}
