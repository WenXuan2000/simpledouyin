package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func SendMessage(uid uint, touid uint, createtime int64, content string) error {
	if uid == touid {
		return common.SendMessageActionWrong
	}
	newMessage := entity.Message{
		ToUserID:   touid,
		UserID:     uid,
		Content:    content,
		CreateTime: createtime,
	}
	err := dao.SqlSession.Model(&entity.Message{}).Create(&newMessage).Error
	return err
}

func GetMessageList(uid uint, touid uint) ([]entity.Message, error) {
	var messagelist []entity.Message
	lasttime := GetLastCheatTime(uid, touid)

	err := dao.SqlSession.Table("Messages").
		Where(" to_user_id in (?) and user_id in (?) and create_time > ?", []uint{touid, uid}, []uint{touid, uid}, lasttime).
		Order("created_at desc").
		Limit(common.MessageNum).
		Find(&messagelist).Error
	if err != nil {
		return messagelist, err
	}
	return messagelist, err
}

func GetLastCheatTime(uid uint, touid uint) (LastTime int64) {
	var messagehistory entity.MessageHistory
	err := dao.SqlSession.Table("MessageHistorys").
		Where(" to_user_id = ? and user_id = ?", touid, uid).
		Find(&messagehistory).Error
	LastTime = messagehistory.LastTime
	if err != nil {
		LastTime = 0
	}
	return LastTime
}

func UpdateLastCheatTime(uid uint, touid uint, lasttime int64) (err error) {
	if err = dao.SqlSession.Model(&entity.MessageHistory{}).
		Where(" to_user_id = ? and user_id = ?", touid, uid).
		First(&entity.MessageHistory{}).Error; gorm.IsRecordNotFoundError(err) {
		newHistory := entity.MessageHistory{
			ToUserID: touid,
			UserID:   uid,
			LastTime: lasttime,
		}
		dao.SqlSession.Model(&entity.User{}).Create(&newHistory)
		return nil
	}
	if err = dao.SqlSession.Model(&entity.MessageHistory{}).
		Where(" to_user_id = ? and user_id = ?", touid, uid).
		Update("last_time", lasttime).
		Error; err != nil {
		return err
	}
	return nil
}
