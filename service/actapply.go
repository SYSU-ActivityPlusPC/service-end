package service

import (
	"github.com/sysu-activitypluspc/service-end/dao"
)

type ActApplyInfo struct {
	dao.ActApplyInfo
}

type ActApplySlice []ActApplyInfo

// GetApplyList returns apply list with given id
func (applys ActApplySlice) GetApplyList(actid int) {

}

// DeleteApply delete apply with act id and apply id
func (apply *ActApplyInfo) DeleteApply() {

}