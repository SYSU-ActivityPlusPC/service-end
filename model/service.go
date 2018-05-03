package model

import (
	"errors"
	"time"

	"github.com/sysu-activitypluspc/service-end/types"
)

// AddActivity insert a activity into the db
func AddActivity(activityInfo types.ActivityInfo) error {
	activity := ActivityInfo{
		ID:              activityInfo.ID,
		Name:            activityInfo.Name,
		StartTime:       msToTime(activityInfo.StartTime),
		EndTime:         msToTime(activityInfo.EndTime),
		Campus:          activityInfo.Campus,
		Location:        activityInfo.Location,
		EnrollCondition: activityInfo.EnrollCondition,
		Sponsor:         activityInfo.Sponsor,
		Type:            activityInfo.Type,
		PubStartTime:    msToTime(activityInfo.PubStartTime),
		PubEndTime:      msToTime(activityInfo.PubEndTime),
		Detail:          activityInfo.Detail,
		Reward:          activityInfo.Reward,
		Introduction:    activityInfo.Introduction,
		Requirement:     activityInfo.Requirement,
		Poster:          activityInfo.Poster,
		Qrcode:          activityInfo.Qrcode,
		Email:           activityInfo.Email,
		Verified:        activityInfo.Verified,
	}
	affected, err := Engine.InsertOne(&activity)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Failed to insert a new one, perhaps it has existed.")
	}
	return nil
}

func msToTime(ms int64) *time.Time {
	ret := time.Unix(0, ms*int64(time.Millisecond))
	return &ret
}
