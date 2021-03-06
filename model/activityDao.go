package model

import (
	"strconv"
	"fmt"
	"time"

	"github.com/sysu-activitypluspc/service-end/types"
)

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
