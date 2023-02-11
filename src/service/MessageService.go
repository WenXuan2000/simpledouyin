package service

import (
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func SendMessage(uid uint, touid uint, createtime string, content string) error {
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

func GetMessageList(uid uint, touid uint) (messagelist []entity.Message, err error) {
	err = dao.SqlSession.Table("Messages").Where(" to_user_id in (?) and user_id in (?)", []uint{touid, uid}, []uint{touid, uid}).
		Order("created_at desc").
		Limit(common.MessageNum).
		Find(&messagelist).Error
	return
}
