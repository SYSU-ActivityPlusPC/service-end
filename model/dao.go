package model

import (
	"strconv"
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

// AddActivity insert a activity into the db
func AddActivity(activityInfo types.ActivityInfo, id int) (int, error) {
	layout := "2006-01-02 15:04"
	starttime, err := time.Parse(layout, activityInfo.StartTime)
	if err != nil {
		return 0, err
	}
	endtime, err := time.Parse(layout, activityInfo.EndTime)
	if err != nil {
		return 0, err
	}
	pubstarttime, err := time.Parse(layout, activityInfo.PubStartTime)
	if err != nil {
		return 0, err
	}
	pubendtime, err := time.Parse(layout, activityInfo.PubEndTime)
	if err != nil {
		return 0, err
	}
	var Enrollendtime *time.Time
	if len(activityInfo.EnrollEndTime) != 0 {
		enrollendtime, err := time.Parse(layout, activityInfo.EnrollEndTime)
		Enrollendtime = &enrollendtime
		if err != nil {
			return 0, err
		}
	}
	activity := ActivityInfo{
		Name:            activityInfo.Name,
		StartTime:       &starttime,
		EndTime:         &endtime,
		Campus:          activityInfo.Campus,
		Location:        activityInfo.Location,
		EnrollCondition: activityInfo.EnrollCondition,
		Sponsor:         activityInfo.Sponsor,
		Type:            activityInfo.Type,
		PubStartTime:    &pubstarttime,
		PubEndTime:      &pubendtime,
		Detail:          activityInfo.Detail,
		Reward:          activityInfo.Reward,
		Introduction:    activityInfo.Introduction,
		Requirement:     activityInfo.Requirement,
		Poster:          activityInfo.Poster,
		Qrcode:          activityInfo.Qrcode,
		Email:           activityInfo.Email,
		Verified:        0,
		EnrollWay:       activityInfo.EnrollWay,
		EnrollEndTime:   Enrollendtime,
		CanEnrolled:     activityInfo.CanEnrolled,
		PCUserID:        id,
	}
	affected, err := Engine.InsertOne(&activity)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

// UpdateActivity update activity
func UpdateActivity(id int, activityInfo types.ActivityInfo) (int, error) {
	layout := "2006-01-02 15:04"
	starttime, err := time.Parse(layout, activityInfo.StartTime)
	if err != nil {
		return 0, err
	}
	endtime, err := time.Parse(layout, activityInfo.EndTime)
	if err != nil {
		return 0, err
	}
	pubstarttime, err := time.Parse(layout, activityInfo.PubStartTime)
	if err != nil {
		return 0, err
	}
	pubendtime, err := time.Parse(layout, activityInfo.PubEndTime)
	if err != nil {
		return 0, err
	}
	var Enrollendtime *time.Time
	if len(activityInfo.EnrollEndTime) != 0 {
		enrollendtime, err := time.Parse(layout, activityInfo.EnrollEndTime)
		Enrollendtime = &enrollendtime
		if err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}
	activity := ActivityInfo{
		Name:            activityInfo.Name,
		StartTime:       &starttime,
		EndTime:         &endtime,
		Campus:          activityInfo.Campus,
		Location:        activityInfo.Location,
		EnrollCondition: activityInfo.EnrollCondition,
		Sponsor:         activityInfo.Sponsor,
		Type:            activityInfo.Type,
		PubStartTime:    &pubstarttime,
		PubEndTime:      &pubendtime,
		Detail:          activityInfo.Detail,
		CanEnrolled:     activityInfo.CanEnrolled,
		EnrollWay:       activityInfo.EnrollWay,
		EnrollEndTime:   Enrollendtime,
		Reward:          activityInfo.Reward,
		Introduction:    activityInfo.Introduction,
		Requirement:     activityInfo.Requirement,
		Poster:          activityInfo.Poster,
		Qrcode:          activityInfo.Qrcode,
		Email:           activityInfo.Email,
		Verified:        0,
	}
	affected, err := Engine.Id(id).Update(&activity)
	return int(affected), err
}

// DeleteActivity remove an activity according to the id
func DeleteActivity(id int) (int, error) {
	affected, err := Engine.Id(id).Delete(&ActivityInfo{})
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		fmt.Println("Failed to delete an activity")
	}
	return int(affected), nil
}

// VerifyActivity set activity verified status
func VerifyActivity(id int, status int) (int, error) {
	activity := new(ActivityInfo)
	activity.Verified = status
	affected, err := Engine.Id(id).Cols("verified").Update(activity)
	if err != nil {
		fmt.Println(err)
	}
	if affected == 0 {
		fmt.Println("Activity status does not need to be changed.")
	}
	return int(affected), nil
}

func msToTime(ms int64) *time.Time {
	ret := time.Unix(0, ms*int64(time.Millisecond))
	return &ret
}

// GetActivityList return wanted activity list with given page number
func GetActivityList(pageNum int, verified int) []ActivityInfo {
	activityList := make([]ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	Engine.Desc("id").Where("verified = ?", verified).Find(&activityList)
	from := pageNum * 10
	if from >= len(activityList) {
		return []ActivityInfo{}
	}
	if from+10 > len(activityList) {
		return activityList[from:]
	}
	return activityList[from : from+10]
}

func IsPublishedByClub(clubId int, intActId int) (bool, error) {
	has, err := Engine.Where("pcuser_id = ? && id = ?", clubId, intActId).Exist(&ActivityInfo{})
	return has, err
}

// GetActStatusNumByClub return the number of activity status
func GetActStatusNumByClub(clubId int) (int, int, int) {
	activityList := make([]ActivityInfo, 0)
	var auditNum, ongoingNum, overNum int = 0, 0, 0
	now := time.Now().Add(time.Hour * 8)
	// Search clubId's activity
	Engine.Desc("id").Where("pcuser_id = ?", clubId).Find(&activityList)
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

// GetActivityListByClub return wanted activity list with given page number
func GetActivityListByClub(pageNum int, clubId int) []ActivityInfo {
	activityList := make([]ActivityInfo, 0)
	// Search clubId's activity
	Engine.Desc("id").Where("pcuser_id = ?", clubId).Find(&activityList)
	from := pageNum * 10
	if from >= len(activityList) {
		return []ActivityInfo{}
	}
	if from+10 > len(activityList) {
		return activityList[from:]
	}
	return activityList[from : from+10]
}

func GetRegisterNumByActId(actId int) int {
	counts, _ := Engine.Where("actid = ?", actId).Count(&ActApplyInfo{})
	s := strconv.FormatInt(counts, 10)
	result, _ := strconv.Atoi(s)
	return result
}

// GetActivityInfo return wanted activity detail information which is given by id
func GetActivityInfo(id int) (bool, ActivityInfo) {
	var activity ActivityInfo

	ok, _ := Engine.ID(id).Get(&activity)
	return ok, activity
}

// CheckPCUser check if the account exists
func CheckPCUser(username string) bool {
	has, _ := Engine.Where("account = ?", username).Exist(&PCUser{})
	return has
}

// GetUserInfo returns user information with given username
func GetUserInfo(account string) PCUser {
	var user PCUser
	_, err := Engine.Where("account = ?", account).Get(&user)
	if err != nil {
		fmt.Println(err)
		user.ID = -1
	}
	return user
}

// AddUser add user
func AddUser(user types.PCUserSignInfo) bool {
	currentTime := time.Now()
	dbuser := PCUser{
		Name:         user.Name,
		Email:        user.Email,
		Logo:         user.Logo,
		Evidence:     user.Evidence,
		Info:         user.Info,
		RegisterTime: &currentTime,
	}
	_, err := Engine.InsertOne(&dbuser)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetUserByEmail return user detail based on email
func GetUserByEmail(email string) *PCUser {
	user := new(PCUser)
	_, err := Engine.Where("email=?", email).Get(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

// GetUserByID return user detail based on id
func GetUserByID(id int) *PCUser {
	user := new(PCUser)
	ok, err := Engine.Where("id=?", id).Get(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !ok {
		user.ID = -1
	}
	return user
}

// VerifyUser set user verified status
func VerifyUser(id int, verify int, email string, password string, currentTime *time.Time) error {
	user := new(PCUser)
	user.Verified = verify
	user.Account = email
	user.Password = password
	user.RegisterTime = currentTime
	_, err := Engine.Where("id=?", id).Cols("verified").Cols("account").Cols("password").Cols("register_time").Update(user)
	return err
}

// GetUserList get all the user with given status
func GetUserList(verify int) []PCUser {
	ret := make([]PCUser, 0)
	err := Engine.Where("verified = ?", verify).Incr("id").Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

// GetApplyListByID return list of apply whose act id is given
func GetApplyListByID(id int) []ActApplyInfo{
	ret := make([]ActApplyInfo, 0)
	err := Engine.Where("actid=?", id).Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

// DeleteApplyByID remove apply with given id
func DeleteApplyByID(actid int, applyid int) bool{
	_, err := Engine.Where("actid=?", actid).And("id=?", applyid).Delete(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// CloseActApplyByID close applicance of act whose id is given
func CloseActApplyByID(id int) bool{
	edit := ActivityInfo{
		CanEnrolled: 2,
	}
	_, err := Engine.Where("id=?", id).Cols("can_enrolled").Update(&edit)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetApplyByID(id int) *ActApplyInfo {
	apply := new(ActApplyInfo)
	ok, err := Engine.Where("id=?", id).Get(apply)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !ok {
		apply.ID = -1
	}
	return apply
}