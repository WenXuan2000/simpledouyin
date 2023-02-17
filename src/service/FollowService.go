package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

// uqid：被关注的人
// uid： 自己
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
	if uid == touid {
		return common.FollowActionWrong
	}
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

func DoFollow(uid uint, touid uint) error {
	// 查看数据库中有没有已经关注的记录
	if IsFollowed(touid, uid) {
		return common.FollowActionDuplicate
	}
	follow := &entity.Follow{
		FollowedId: touid,
		FollowId:   uid,
	}
	// 事务操作
	err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		err := db.Model(&entity.Follow{}).Create(&follow).Error
		if err != nil {
			return err
		}
		err = UpdateUserFollowCount(uid, "1", db)
		if err != nil {
			return err
		}
		err = UpdateUserFollowerCount(touid, "1", db)
		return err
	})
	if err1 != nil {
		return err1
	}
	return nil
}
func UnFollow(uid uint, touid uint) error {
	//先查询有没有这条记录，防止重复删除
	if err := dao.SqlSession.Model(&entity.Follow{}).Where("followed_id = ? AND follow_id = ?", touid, uid).First(&entity.Follow{}).Error; err != nil {
		return err
	}
	// 事务操作
	err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		// 这里删除不存在的记录也不会报错
		err := db.Model(&entity.Follow{}).Where("followed_id = ? AND follow_id = ?", touid, uid).Delete(&entity.Follow{}).Error
		if err != nil {
			return err
		}
		err = UpdateUserFollowCount(uid, "-1", db)
		if err != nil {
			return err
		}
		err = UpdateUserFollowerCount(touid, "-1", db)
		return err
	})
	if err1 != nil {
		return err1
	}
	return nil

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
func FollowerListGet(uid uint) ([]entity.User, error) {
	var FollowerList []entity.User
	FollowerList = make([]entity.User, 0)
	if err := dao.SqlSession.Model(&entity.User{}).
		Joins("left join "+entity.Follow{}.TableName()+" on "+entity.User{}.TableName()+".id = "+entity.Follow{}.TableName()+".follow_id").
		Where(entity.Follow{}.TableName()+".followed_id=? AND "+entity.Follow{}.TableName()+".deleted_at is null", uid).
		Scan(&FollowerList).Error; err != nil {
		return FollowerList, err
	}
	return FollowerList, nil
}

func FriendListGet(uid uint) ([]entity.User, error) {
	var FollowerList []entity.User
	FollowerList = make([]entity.User, 0)
	// 首先取出互相关注的被关注者的id
	var idlist []uint
	if err := dao.SqlSession.Table("follow a").
		Joins("left join follow b on a.follow_id =b.followed_id And b.follow_id =a.followed_id").
		Where("a.follow_id=? AND a.deleted_at is null AND b.deleted_at is null", uid).
		Pluck("a.followed_id", &idlist).Error; err != nil {
		return FollowerList, err
	}
	if err := dao.SqlSession.Model(&entity.User{}).
		Where("id in (?)", idlist).
		Scan(&FollowerList).Error; err != nil {
		return FollowerList, err
	}
	return FollowerList, nil
}
