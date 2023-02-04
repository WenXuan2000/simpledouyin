package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func IsFollowed(uqid uint, uid uint) (ok bool) {
	// 自己不能follow自己
	if uqid == uid {
		return false
	}
	var followExist = &entity.Follow{}
	if err := dao.SqlSession.Model(&entity.Follow{}).Where("followed_id = ? AND follow_id = ?", "uqid", "uid").First(&followExist).Error; gorm.IsRecordNotFoundError(err) {
		return false
	} else {
		return true
	}
}
