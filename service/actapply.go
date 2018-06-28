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
	if actid <= 0 {
		applys = nil
		return
	}

	// Get applys from dao
	apply := new(dao.ActApplyInfo)
	apply.ActId = actid
	session := GetSession()
	defer DeleteSession(session, true)
	daoApplys := apply.ListByActID(session)
	// Construct ActApplySlice
	if daoApplys == nil {
		applys = nil
		DeleteSession(session, false)
		return
	}
	for _, v := range daoApplys {
		tmp := ActApplyInfo{v}
		applys = append(applys, tmp)
	}
}

// DeleteApply delete apply with act id and apply id
func (apply *ActApplyInfo) DeleteApply() bool {
	daoApply := new(dao.ActApplyInfo)
	daoApply.ID = apply.ID

	// Check if the apply matches
	session := GetSession()
	defer DeleteSession(session, true)
	daoApply.Get(session)
	if daoApply == nil || daoApply.ID < 0 || daoApply.ActId != apply.ActId {
		apply = nil
		return false
	}
	// Delete apply
	return daoApply.Delete(session)
}
