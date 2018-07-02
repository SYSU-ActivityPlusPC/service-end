package dao

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

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
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

func (act *ActivityInfo) UpdateVerifiedStatus(session *xorm.Session) (int, error) {
	id := act.ID
	affected, err := session.Id(id).Cols("verified").Update(act)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

// TODO: Use map to update
func (act *ActivityInfo) Update(session *xorm.Session) (int, error) {
	id := act.ID
	affected, err := session.Id(id).Update(&act)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

func (act *ActivityInfo) UpdateEnrolled(session *xorm.Session) (int, error) {
	id := act.ID
	affected, err := session.Where("id=?", id).Cols("can_enrolled").Update(&act)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

func (act *ActivityInfo) Get(session *xorm.Session) (bool, error){
	id := act.ID
	has, err := session.ID(id).Get(act)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return has, nil
}

func (act *ActivityInfo) ListStatusNumByClubID(session *xorm.Session) (int, int, int, error) {
	clubId := act.PCUserID
	activityList := make([]ActivityInfo, 0)
	var auditNum, ongoingNum, overNum int = 0, 0, 0
	now := time.Now().Add(time.Hour * 8)
	// Search clubId's activity
	err := session.Desc("id").Where("pcuser_id = ?", clubId).Find(&activityList)
	if err != nil {
		fmt.Println(err)
		return -1, -1, -1, err
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
	return auditNum, ongoingNum, overNum, nil
}

func (act *ActivityInfo) ListVerifiedByClubID(session *xorm.Session, pageNum int) ([]ActivityInfo, error) {
	clubId := act.PCUserID
	activityList := make([]ActivityInfo, 0)
	// Search clubId's activity
	err := session.Desc("id").Where("pcuser_id = ?", clubId).And("verified = 1").Find(&activityList)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	from := pageNum * 10
	if from >= len(activityList) {
		return []ActivityInfo{}, nil
	}
	if from+10 > len(activityList) {
		return activityList[from:], nil
	}
	return activityList[from : from+10], nil
}

func (act *ActivityInfo) ListByVerifiedStatus(session *xorm.Session, pageNum int) ([]ActivityInfo, error) {
	verified := act.Verified
	activityList := make([]ActivityInfo, 0)
	// Search verified activity
	// 0 stands for no pass
	// 1 stands for pass
	// 2 stands for not yet verified
	err := session.Desc("id").Where("verified = ?", verified).Find(&activityList)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	from := pageNum * 10
	if from >= len(activityList) {
		return []ActivityInfo{}, nil
	}
	if from+10 > len(activityList) {
		return activityList[from:], nil
	}
	return activityList[from : from+10], nil
}
