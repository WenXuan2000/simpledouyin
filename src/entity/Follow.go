package entity

import "github.com/jinzhu/gorm"

func (Follow) TableName() string {
	return "follow"
}

// Follow关注逻辑： FollowId->FollowedId
type Follow struct {
	gorm.Model
	FollowedId uint `json:"followed_id"`
	FollowId   uint `json:"follow_id"`
}
