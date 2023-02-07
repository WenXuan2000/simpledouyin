package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func IsFollowed(uqid uint, uid uint) (ok bool) {
	// 自己不能follow自己
	if uqid == uid {
		return false
	}
	var followExist = &entity.Follow{}
	if err := dao.SqlSession.Model(&entity.Follow{}).Where("followed_id = ? AND follow_id = ?", uqid, uid).First(&followExist).Error; gorm.IsRecordNotFoundError(err) {
		return false
	} else {
		return true
	}
}

// 修改关注列表
func FollowAction(uid uint, touid uint, action_type string) (err error) {
	switch action_type {
	case "1":
		err = DoFollow(uid, touid)
	case "2":
		err = UnFollow(uid, touid)
	default:
		err = common.ActionTypeWrong
	}
	return err
}

func DoFollow(uid uint, touid uint) (err error) {
	// 查看数据库中有没有已经关注的记录
	if IsFollowed(touid, uid) {
		return common.FollowActionDuplicate
	}
	follow := &entity.Follow{
		FollowedId: touid,
		FollowId:   uid,
	}
	err = dao.SqlSession.Model(&entity.Follow{}).Create(&follow).Error
	if err != nil {
		return err
	}
	err = UpdateUserFollowCount(uid, "1")
	if err != nil {
		return err
	}
	err = UpdateUserFollowerCount(touid, "1")
	return err
}
func UnFollow(uid uint, touid uint) (err error) {
	//先查询有没有这条记录，防止重复删除
	if err = dao.SqlSession.Model(&entity.Follow{}).Where("followed_id = ? AND follow_id = ?", touid, uid).First(&entity.Follow{}).Error; err != nil {
		return err
	}
	// 这里删除不存在的记录也不会报错
	err = dao.SqlSession.Model(&entity.Follow{}).Where("followed_id = ? AND follow_id = ?", touid, uid).Delete(&entity.Follow{}).Error
	if err != nil {
		return err
	}
	err = UpdateUserFollowCount(uid, "-1")
	if err != nil {
		return err
	}
	err = UpdateUserFollowerCount(touid, "-1")
	return err
}

func FollowListGet(uid uint) ([]entity.User, error) {
	var FollowList []entity.User
	FollowList = make([]entity.User, 0)
	if err := dao.SqlSession.Model(&entity.User{}).
		Joins("left join "+entity.Follow{}.TableName()+" on "+entity.User{}.TableName()+".id = "+entity.Follow{}.TableName()+".followed_id").
		Where(entity.Follow{}.TableName()+".follow_id=? AND "+entity.Follow{}.TableName()+".deleted_at is null", uid).
		Scan(&FollowList).Error; err != nil {
		return FollowList, err
	}
	return FollowList, nil
}
