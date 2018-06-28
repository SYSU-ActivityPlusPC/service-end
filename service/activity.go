package service

import (
	"github.com/sysu-activitypluspc/service-end/dao"
)

// Activity stores json format the front-end wanted
type ActivityInfo struct {
	dao.ActivityInfo
}

type ActivitySlice []ActivityInfo

// AddActivity add activity
// 0 - already exist
// 1 - finish
// 2 - fail
func (act *ActivityInfo) AddActivity() int {
	// TODO: Check data in activity
	daoAct := act.ActivityInfo
	session := GetSession()
	defer DeleteSession(session, true)
	affected := daoAct.Insert(session)
	if affected == -1 {
		DeleteSession(session, false)
		return 2
	}
	if affected == 0 {
		return 0
	}
	return 1
}

// GetActivityNumber receives club id, returns different status' number
// with order: audit, ongoing, over
func (act *ActivityInfo) GetActivityNumber(clubID int) (bool, []int) {
	session := GetSession()
	defer DeleteSession(session, true)
	a, b, c := act.ActivityInfo.ListStatusNumByClubID(session)
	if a == -1 {
		return false, []int{a, b, c}
	}
	return true, []int{a, b, c}
}

// GetActivityListByClud for club
func (acts ActivitySlice) GetActivityListByClud(page int, clubid int) {
	act := new(dao.ActivityInfo)
	act.PCUserID = clubid
	session := GetSession()
	defer DeleteSession(session, true)
	daoActs := act.ListVerifiedByClubID(session, page)
	if daoActs == nil {
		acts = nil
		return
	}
	for _, v := range daoActs {
		tmp := ActivityInfo{v}
		acts = append(acts, tmp)
	}
}

// GetActivityListByadmin for administrator
func (acts ActivitySlice) GetActivityListByAdmin(page int, verified int) {
	act := new(dao.ActivityInfo)
	act.Verified = verified
	session := GetSession()
	defer DeleteSession(session, true)
	daoActs := act.ListByVerifiedStatus(session, page)
	if daoActs == nil {
		acts = nil
		return
	}
	for _, v := range daoActs {
		tmp := ActivityInfo{v}
		acts = append(acts, tmp)
	}
}

// GetActivityInfor returns activity info with the given id
func (act *ActivityInfo) GetActivityInfo(id int) {
	if act.ID <= 0 {
		act = nil
		return
	}
	session := GetSession()
	defer DeleteSession(session, true)
	act.Get(session)
}

// ModifyActivity update activity information with id
func (act *ActivityInfo) ModifyActivity() bool {
	if act.ID <= 0 {
		act = nil
		return false
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected := act.Update(session)
	if affected == -1 {
		act = nil
		DeleteSession(session, false)
		return false
	}
	return true
}

// DeleteActivity delete activity with id
// 0 - error
// 1 - ok
// 2 - no content
func (act *ActivityInfo) DeleteActivity() int {
	if act.ID <= 0 {
		act = nil
		return 0
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected := act.Delete(session)
	if affected == -1 {
		act = nil
		DeleteSession(session, false)
		return 0
	}
	if affected == 0 {
		return 2
	}
	return 1
}

// AduitActivity aduit activity with id
func (act *ActivityInfo) AduitActivity() bool {
	if act.ID <= 0 {
		act = nil
		return false
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected := act.UpdateVerifiedStatus(session)
	if affected == -1 {
		DeleteSession(session, false)
		return false
	}
	return true
}

// CloseActivity close activity with id
func (act *ActivityInfo) CloseActivity() bool {
	if act.ID <= 0 {
		act = nil
		return false
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected := act.UpdateEnrolled(session)
	if affected == -1 {
		return false
	}
	return true
}
