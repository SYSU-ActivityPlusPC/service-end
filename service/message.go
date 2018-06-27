package service

import (
	"github.com/sysu-activitypluspc/service-end/dao"
)

type Message struct {
	dao.Message
}

type MessageSlice []Message

// PublishMessage publish a message
func (msg *Message) PublishMessage() {

}

// ListMessages list messages with given page
func (msgs MessageSlice) ListMessages(page int) {

}

// GetMessageInformation returns message information with given message id
func (msg *Message) GetMessageInformation() {

}

// DeleteMessage delete message with id
func (msg *Message) DeleteMessage() {

}