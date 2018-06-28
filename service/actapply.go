package service

import (
	"errors"

	"github.com/sysu-activitypluspc/service-end/dao"
)

type ActApplyInfo struct {
	dao.ActApplyInfo
}

type ActApplySlice []ActApplyInfo

// GetApplyList returns apply list with given id
func (applys ActApplySlice) GetApplyList(actid int) (int, error) {
	if actid <= 0 {
		return 400, errors.New("Invalid activity id")
	}

	// Get applys from dao
	apply := new(dao.ActApplyInfo)
	apply.ActId = actid
	session := GetSession()
	defer DeleteSession(session, true)
	daoApplys, err := apply.ListByActID(session)
	// Construct ActApplySlice
	if err == nil {
		DeleteSession(session, false)
		return 500, err
	}
	for _, v := range daoApplys {
		tmp := ActApplyInfo{v}
		applys = append(applys, tmp)
	}
	return 200, nil
}

// DeleteApply delete apply with act id and apply id
func (apply *ActApplyInfo) DeleteApply() (int, error) {
	daoApply := new(dao.ActApplyInfo)
	daoApply.ID = apply.ID

	session := GetSession()
	defer DeleteSession(session, true)
	affected, err := daoApply.Delete(session)
	if err != nil {
		DeleteSession(session, false)
		return 500, err
	}
	if affected == 0 {
		return 204, nil
	}
	return 200, nil
}
