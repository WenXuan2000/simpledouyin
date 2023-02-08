package service

import (
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

// CreateVideo 添加一条视频信息
func CreateVideo(video *entity.Video) {
	dao.SqlSession.Table("videos").Create(&video)
}

// PubNumOfUser发布总数
