package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func CommentsCreate(comment *entity.Comment) (err error) {
	err = dao.SqlSession.Model(&entity.Comment{}).Create(comment).Error
	return
}
func VideosCommentCountAdd(vid uint) error {
	if err := dao.SqlSession.Model(&entity.Video{}).Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count+1")).Error; err != nil {
		return err
	}
	return nil
}
func CommentsDelete(cid uint) error {
	if err := dao.SqlSession.Model(&entity.Comment{}).
		Where("id = ?", cid).
		Delete(&entity.Comment{}).Error; err != nil {
		return err
	}
	return nil
}
func VideosCommentCountReduce(vid uint) error {
	if err := dao.SqlSession.Model(&entity.Video{}).Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count-1")).Error; err != nil {
		return err
	}
	return nil
}

func GetCommentListByVideoId(vid uint) (commentList []entity.Comment, err error) {
	if err = dao.SqlSession.Model(&entity.Comment{}).
		Where("video_id=?", vid).
		Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}
