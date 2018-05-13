package model

import (
	"fmt"
	"time"

	"github.com/sysu-activitypluspc/service-end/types"
)

// AddActivity insert a activity into the db
func AddActivity(activityInfo types.ActivityInfo) (int, error) {
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

func UpdateActivity(id int, activityInfo types.ActivityInfo) (int, error) {
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
		Verified:        0,
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
	activity := ActivityInfo{
		Verified: status,
	}
	affected, _ := Engine.Id(id).Update(&activity)
	if affected == 0 {
		fmt.Println("Activity status does not need to be changed.")
	}
	return int(affected), nil
}

func msToTime(ms int64) *time.Time {
	ret := time.Unix(0, ms*int64(time.Millisecond))
	return &ret
}
