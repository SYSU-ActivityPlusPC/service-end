package service

import (
	"errors"

	"github.com/sysu-activitypluspc/service-end/dao"
)

// Activity stores json format the front-end wanted
type ActivityInfo struct {
	dao.ActivityInfo
}

type ActivitySlice []ActivityInfo

// AddActivity add activity
func (act *ActivityInfo) AddActivity() (int, error) {
	// TODO: Check data in activity
	daoAct := act.ActivityInfo
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := daoAct.Insert(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, nil
}

// GetActivityNumber receives club id, returns different status' number
// with order: audit, ongoing, over
func (act *ActivityInfo) GetActivityNumber(clubID int) ([]int, error) {
	session := GetSession()
	defer DeleteSession(session, true)
	a, b, c, err := act.ActivityInfo.ListStatusNumByClubID(session)
	if err != nil {
		return []int{a, b, c}, err
	}
	return []int{a, b, c}, nil
}

// GetActivityListByClud for club
func (acts ActivitySlice) GetActivityListByClud(page int, clubid int) (int, error) {
	act := new(dao.ActivityInfo)
	act.PCUserID = clubid
	session := GetSession()
	defer DeleteSession(session, true)
	daoActs, err := act.ListVerifiedByClubID(session, page)
	if err != nil {
		return 500, err
	}
	for _, v := range daoActs {
		tmp := ActivityInfo{v}
		acts = append(acts, tmp)
	}
	if len(acts) == 0 {
		return 204, nil
	}
	return 200, nil
}

// GetActivityListByadmin for administrator
func (acts ActivitySlice) GetActivityListByAdmin(page int, verified int) (int, error) {
	act := new(dao.ActivityInfo)
	act.Verified = verified
	session := GetSession()
	defer DeleteSession(session, true)
	daoActs, err := act.ListByVerifiedStatus(session, page)
	if err != nil {
		return 500, err
	}
	for _, v := range daoActs {
		tmp := ActivityInfo{v}
		acts = append(acts, tmp)
	}
	if len(acts) == 0 {
		return 204, nil
	}
	return 200, nil
}

// GetActivityInfor returns activity info with the given id
func (act *ActivityInfo) GetActivityInfo(id int) (int, error) {
	if act.ID <= 0 {
		act = nil
		return 400, errors.New("Invalid activity id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	has, err := act.Get(session)
	if err != nil {
		return 500, err
	}
	if !has {
		return 204, nil
	}
	return 200, nil
}

// ModifyActivity update activity information with id
func (act *ActivityInfo) ModifyActivity() (int, error) {
	if act.ID <= 0 {
		act = nil
		return 400, errors.New("Invalid activity id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	_, err := act.Update(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	return 200, nil
}

// DeleteActivity delete activity with id
// 0 - error
// 1 - ok
// 2 - no content
func (act *ActivityInfo) DeleteActivity() (int, error) {
	if act.ID <= 0 {
		return 400, errors.New("Invalid activity id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := act.Delete(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, nil
}

// AduitActivity aduit activity with id
func (act *ActivityInfo) AduitActivity() (int, error) {
	if act.ID <= 0 {
		return 400, errors.New("Invalid activity id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	_, err := act.UpdateVerifiedStatus(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	return 200, nil
}

// CloseActivity close activity with id
func (act *ActivityInfo) CloseActivity() (int, error) {
	if act.ID <= 0 {
		return 400, errors.New("Invalid activity id")
	}
	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := act.UpdateEnrolled(session)
	if err != nil {
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, nil
}

// Is act published by account
func (act *ActivityInfo) CheckMessageCorrectness() (bool, error) {
	userid := act.PCUserID
	session := GetSession()
	defer DeleteSession(session, true)
	has, err := act.Get(session)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	if userid != act.PCUserID {
		return false, nil
	}
	return true, nil
}
