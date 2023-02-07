package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

func FavoriteAction(userid uint, videoid uint, actiontype string) (err error) {
	// actiontype == 1 点赞
	// actiontype == 2 取消点赞
	switch actiontype {
	case "1":
		err = DoFavorite(userid, videoid)
	case "2":
		err = UnFavorite(userid, videoid)
	default:
		err = common.ActionTypeWrong
	}
	return
}

func DoFavorite(userid, videoid uint) (err error) {
	favoriteInfo := entity.Favorite{
		UserId:  userid,
		VideoId: videoid,
	}
	// favorite表 创建
	if err = dao.SqlSession.Model(&entity.Favorite{}).Create(&favoriteInfo).Error; err != nil {
		return err
	}
	// 对应的video数据 favorite_count+1
	if err = dao.SqlSession.Model(&entity.Video{}).
		Where("id = ?", videoid).
		Update("favorite_count", gorm.Expr("favorite_count + 1")).
		Error; err != nil {
		return err
	}
	// 作者的user表的total_favorite增加
	// 获取作者的id
	var authorId uint
	if authorId, err = GetVideoAuthor(videoid); err != nil {
		return err
	}
	// +1
	if err = dao.SqlSession.Model(&entity.User{}).
		Where("id=?", authorId).
		Update("total_favorited", gorm.Expr("total_favorited + 1")).
		Error; err != nil {
		return err
	}
	// userid对应的favorite_count +1
	if err = dao.SqlSession.Model(&entity.User{}).
		Where("id=?", userid).
		Update("favorite_count", gorm.Expr("favorite_count + 1")).
		Error; err != nil {
		return err
	}
	return nil
}

func UnFavorite(userid, videoid uint) (err error) {
	// 删除
	if err = dao.SqlSession.Model(&entity.Favorite{}).
		Where("user_id = ? AND video_id = ?", userid, videoid).
		Delete(&entity.Favorite{}).Error; err != nil {
		return err
	}
	// 对应的video数据 favorite_count-1
	if err = dao.SqlSession.Model(&entity.Video{}).
		Where("id = ?", videoid).
		Update("favorite_count", gorm.Expr("favorite_count - 1")).
		Error; err != nil {
		return err
	}
	// 作者的user表的total_favorite增加
	// 获取作者的id
	var authorId uint
	if authorId, err = GetVideoAuthor(videoid); err != nil {
		return err
	}
	if err = dao.SqlSession.Model(&entity.User{}).
		Where("id=?", authorId).
		Update("total_favorited", gorm.Expr("total_favorited - 1")).
		Error; err != nil {
		return err
	}
	// userid对应的favorite_count - 1
	if err = dao.SqlSession.Model(&entity.User{}).
		Where("id=?", userid).
		Update("favorite_count", gorm.Expr("favorite_count - 1")).
		Error; err != nil {
		return err
	}
	return nil
}

func CheckFavorite(userid, videoid uint) bool {
	if err := dao.SqlSession.Model(&entity.Favorite{}).
		Where("user_id = ? AND video_id = ?", userid, videoid).
		First(&entity.Favorite{}).Error; err != nil {
		return false
	}
	return true
}