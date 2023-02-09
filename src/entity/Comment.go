package entity

import "github.com/jinzhu/gorm"

func (Comment) TableName() string {
	return "comments"
}

type Comment struct {
	gorm.Model
	UserId  uint   `json:"user_id"`
	VideoId uint   `json:"video_id"`
	Content string `json:"content"`
}
