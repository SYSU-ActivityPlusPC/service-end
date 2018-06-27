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
func (act *ActivityInfo) AddActivity() {

}

// GetActivityNumber receives club id, returns different status' number
// with order: audit, ongoing, over
func (act *ActivityInfo) GetActivityNumber(clubID int) []int {

}

// GetActivityList receives the required page number and verified status
// and returns the list of activities
func (acts *ActivitySlice) getActivityList(page int, verified int) {

}

// GetActivityListByClud for club
func (acts *ActivitySlice) GetActivityListByClud(page int) {

}

// GetActivityListByadmin for administrator
func (acts *ActivitySlice) GetActivityListByAdmin(page int, verified int) {

}

// GetActivityInfor returns activity info with the given id
func (act *ActivityInfo) GetActivityInfor(id int) {

}

// ModifyActivity update activity information with id
func (act *ActivityInfo) ModifyActivity() {

}

// DeleteActivity delete activity with  id
func (act *ActivityInfo) DeleteActivity() {

}

// AduitActivity aduit activity with id
func (act *ActivityInfo) AduitActivity() {

}

// CloseActivity close activity with id
func (act *ActivityInfo) CloseActivity() {

}
