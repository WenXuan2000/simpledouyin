package entity

import "github.com/jinzhu/gorm"

func (Message) TableName() string {
	return "Messages"
}

type Message struct {
	gorm.Model
	ToUserID   uint   `json:"to_user_id"`
	UserID     uint   `json:"user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
}
