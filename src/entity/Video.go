package entity

import "github.com/jinzhu/gorm"

func (Video) TableName() string {
	return "videos"
}

type Video struct {
	gorm.Model
	AuthorId      uint   `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	Title         string `json:"title"`
}
