package dao

import (
	"fmt"
	"strconv"

	"github.com/go-xorm/xorm"
)

func (apply *ActApplyInfo) Get(session *xorm.Session) (bool, error) {
	id := apply.ID
	has, err := session.Where("id=?", id).Get(apply)
	if err != nil {
		fmt.Println(err)
		apply = nil
		return false, err
	}
	return has, nil
}

func (apply *ActApplyInfo) ListByActID(session *xorm.Session) ([]ActApplyInfo, error) {
	actid := apply.ActId
	ret := make([]ActApplyInfo, 0)
	err := session.Where("actid=?", actid).Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

func (apply *ActApplyInfo) GetRegisterNum(session *xorm.Session) (int, error) {
	actId := apply.ActId
	counts, err := session.Where("actid = ?", actId).Count(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	s := strconv.FormatInt(counts, 10)
	result, _ := strconv.Atoi(s)
	return result, nil
}

func (apply *ActApplyInfo) Delete(session *xorm.Session) (int, error) {
	actid := apply.ActId
	applyid := apply.ID
	affected, err := session.Where("actid=?", actid).And("id=?", applyid).Delete(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}
