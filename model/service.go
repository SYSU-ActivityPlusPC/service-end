package model

import (
	"fmt"
	"time"

	"github.com/sysu-activitypluspc/service-end/types"
)

// AddActivity insert a activity into the db
func AddActivity(activityInfo types.StringActivityInfo) (int, error) {
	starttime, err := time.Parse("2006-01-01", activityInfo.StartTime)
	endtime, err := time.Parse("2006-01-01", activityInfo.EndTime)
	pubstarttime, err := time.Parse("2006-01-01", activityInfo.PubStartTime)
	pubendtime, err := time.Parse("2006-01-01", activityInfo.PubEndTime)
	if err != nil {
		return 0, err
	}
	activity := ActivityInfo{
		ID:              activityInfo.ID,
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
	}
	affected, err := Engine.InsertOne(&activity)
	return int(affected), nil
}

func UpdateActivity(id int, activityInfo types.StringActivityInfo) (int, error) {
	starttime, err := time.Parse("2006-01-01", activityInfo.StartTime)
	endtime, err := time.Parse("2006-01-01", activityInfo.EndTime)
	pubstarttime, err := time.Parse("2006-01-01", activityInfo.PubStartTime)
	pubendtime, err := time.Parse("2006-01-01", activityInfo.PubEndTime)
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
		Reward:          activityInfo.Reward,
		Introduction:    activityInfo.Introduction,
		Requirement:     activityInfo.Requirement,
		Poster:          activityInfo.Poster,
		Qrcode:          activityInfo.Qrcode,
		Email:           activityInfo.Email,
	}
	affected, err := Engine.Id(id).Update(&activity)
	return int(affected), err
}

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
func GetActivityList(pageNum int) []ActivityInfo {
	activityList := make([]ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	Engine.Desc("id").Limit(10, pageNum*10).Find(&activityList)
	return activityList
}

// GetActivityInfo return wanted activity detail information which is given by id
func GetActivityInfo(id int) (bool, ActivityInfo) {
	var activity ActivityInfo

	ok, _ := Engine.ID(id).Get(&activity)
	return ok, activity
}
