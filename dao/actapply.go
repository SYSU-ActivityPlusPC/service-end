package dao

import (
	"fmt"
	"strconv"

	"github.com/go-xorm/xorm"
)

func (apply *ActApplyInfo) Get(session *xorm.Session) {
	id := apply.ID
	has, err := session.Where("id=?", id).Get(apply)
	if err != nil {
		fmt.Println(err)
		apply = nil
	}
	if !has {
		apply.ID = -1
	}
}

func (apply *ActApplyInfo) ListByActID(session *xorm.Session) []ActApplyInfo {
	actid := apply.ActId
	ret := make([]ActApplyInfo, 0)
	err := session.Where("actid=?", actid).Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

func (apply *ActApplyInfo) GetRegisterNum(session *xorm.Session) int {
	actId := apply.ActId
	counts, err := session.Where("actid = ?", actId).Count(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return -1
	}
	s := strconv.FormatInt(counts, 10)
	result, _ := strconv.Atoi(s)
	return result
}

func (apply *ActApplyInfo) Delete(session *xorm.Session) bool {
	actid := apply.ActId
	applyid := apply.ID
	_, err := session.Where("actid=?", actid).And("id=?", applyid).Delete(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
