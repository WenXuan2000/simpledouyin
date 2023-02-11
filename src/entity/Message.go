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
	CreateTime int64  `json:"create_time"`
}

func (MessageHistory) TableName() string {
	return "MessageHistorys"
}

type MessageHistory struct {
	gorm.Model
	ToUserID uint  `json:"to_user_id"`
	UserID   uint  `json:"user_id"`
	LastTime int64 `json:"Last_Time"`
}
