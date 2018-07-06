package model

import (
	"fmt"
)

// GetApplyListByID return list of apply whose act id is given
func GetApplyListByID(id int) []ActApplyInfo{
	ret := make([]ActApplyInfo, 0)
	err := Engine.Where("actid=?", id).Find(&ret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

// DeleteApplyByID remove apply with given id
func DeleteApplyByID(actid int, applyid int) bool{
	_, err := Engine.Where("actid=?", actid).And("id=?", applyid).Delete(&ActApplyInfo{})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// CloseActApplyByID close applicance of act whose id is given
func CloseActApplyByID(id int) bool{
	edit := ActivityInfo{
		CanEnrolled: 2,
	}
	_, err := Engine.Where("id=?", id).Cols("can_enrolled").Update(&edit)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetApplyByID(id int) *ActApplyInfo {
	apply := new(ActApplyInfo)
	ok, err := Engine.Where("id=?", id).Get(apply)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !ok {
		apply.ID = -1
	}
	return apply
}