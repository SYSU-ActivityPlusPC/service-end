package dao

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

// Activity
func (act *ActivityInfo) Insert(session *xorm.Session) (int, error) {
	affected, err := session.InsertOne(&act)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

func (act *ActivityInfo) Delete(session *xorm.Session) (int, error) {
	id := act.ID
	affected, err := session.Id(id).Delete(&ActivityInfo{})
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		fmt.Println("Failed to delete an activity")
	}
	return int(affected), nil
}

func (act *ActivityInfo) UpdateVerifiedStatus(session *xorm.Session) (int, error) {
	id := act.ID
	affected, err := session.Id(id).Cols("verified").Update(act)
	if err != nil {
		fmt.Println(err)
	}
	return int(affected), nil
}

func (act *ActivityInfo) Update(session *xorm.Session) (int, error) {
	id := act.ID
	affected, err := session.Id(id).Update(&act)
	if err != nil {
		fmt.Println(err)
	}
	return int(affected), err
}

func (act *ActivityInfo) Get(session *xorm.Session) {
	id := act.ID
	_, err := session.ID(id).Get(act)
	if err != nil {
		act = nil
		fmt.Println(err)
	}
}

func (act *ActivityInfo) ListStatusNumByClubID(session *xorm.Session) (int, int, int) {
	clubId := act.PCUserID
	activityList := make([]ActivityInfo, 0)
	var auditNum, ongoingNum, overNum int = 0, 0, 0
	now := time.Now().Add(time.Hour * 8)
	// Search clubId's activity
	err := session.Desc("id").Where("pcuser_id = ?", clubId).Find(&activityList)
	if err != nil {
		fmt.Println(err)
		return -1, -1, -1
	}
	for i := 0; i < len(activityList); i++ {
		if activityList[i].Verified == 2 {
			continue
		} else if activityList[i].Verified == 0 {
			auditNum++
		} else if activityList[i].Verified == 1 && activityList[i].PubEndTime.After(now) {
			ongoingNum++
		} else {
			overNum++
		}
	}
	return auditNum, ongoingNum, overNum
}

func (act *ActivityInfo) ListByClubID(session *xorm.Session, pageNum int) []ActivityInfo {
	clubId := act.PCUserID
	activityList := make([]ActivityInfo, 0)
	// Search clubId's activity
	err := session.Desc("id").Where("pcuser_id = ?", clubId).Find(&activityList)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	from := pageNum * 10
	if from >= len(activityList) {
		return []ActivityInfo{}
	}
	if from+10 > len(activityList) {
		return activityList[from:]
	}
	return activityList[from : from+10]
}

func (act *ActivityInfo) ListByVerifiedStatus(session *xorm.Session, pageNum int) []ActivityInfo {
	verified := act.Verified
	activityList := make([]ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	err := session.Desc("id").Where("verified = ?", verified).Find(&activityList)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	from := pageNum * 10
	if from >= len(activityList) {
		return []ActivityInfo{}
	}
	if from+10 > len(activityList) {
		return activityList[from:]
	}
	return activityList[from : from+10]
}

// Apply
func (apply *ActApplyInfo) ListByActID(session *xorm.Session) []ActApplyInfo {
	actid := apply.ActId
	ret := make([]ActApplyInfo, 0)
	err := session.Where("actid=?", actid).Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

func (apply *ActApplyInfo) Delete(session *xorm.Session) bool {
	actid := apply.ActId
	applyid := apply.ID
	_, err := session.Where("actid=?", actid).And("id=?", applyid).Delete(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// PC user
func (user *PCUser) GetByAccount(session *xorm.Session) {
	account := user.Account
	_, err := session.Where("account = ?", account).Get(&user)
	if err != nil {
		fmt.Println(err)
		user.ID = -1
	}
}

func (user *PCUser) GetByEmail() {
	email := user.Email
	_, err := session.Where("email=?", email).Get(user)
	if err != nil {
		fmt.Println(err)
		user = nil
	}
}

func (user *PCUser) GetByID() {
	id := user.ID
	_, err := session.Where("id=?", id).Get(user)
	if err != nil {
		fmt.Println(err)
		user = nil
	}
}

func (user *PCUser) UpdateVerifiedStatus(session *xorm.Session) (int, error) {
	id := user.ID
	affected, err := session.Where("id=?", id).Cols("verified").Cols("account").Cols("password").Cols("register_time").Update(user)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), err
}

func (user *PCUser) ListByVerifiedStatus(session *xorm.Session) []PCUser{
	verify := user.Verified
	ret := make([]PCUser, 0)
	err := session.Where("verified = ?", verify).Incr("id").Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

func (user *PCUser) Insert(session *xorm.Session) bool {
	_, err := session.InsertOne(&user)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Message
// TODO: Add pcuser to message_pcuser table
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

func (msg *Message) ListSentTo(session *xorm.Session) ([]string, error) {
	messageId := msg.ID
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
