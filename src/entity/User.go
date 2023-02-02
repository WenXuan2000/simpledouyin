package entity

import "github.com/jinzhu/gorm"

func (User) TableName() string {
	return "users"
}

type User struct {
	gorm.Model
	Name           string `json:"name"`
	Password       string `json:"password"`
	FollowCount    int    `json:"follow_count"`
	FollowerCount  int    `json:"follower_count"`
	TotalFavorited int    `json:"total_favorited"`
	FavoriteCount  int    `json:"favorite_count"`
}
