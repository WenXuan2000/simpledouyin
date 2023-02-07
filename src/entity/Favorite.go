package entity

import "github.com/jinzhu/gorm"

func (Favorite) TableName() string {
	return "favorites"
}

type Favorite struct {
	gorm.Model
	UserId  uint `json:"user_id"`
	VideoId uint `json:"video_id"`
}
