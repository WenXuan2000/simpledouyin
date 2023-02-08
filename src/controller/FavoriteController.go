package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/service"
	"strconv"
)

type FavoriteResponse struct {
	common.Response
	VideoList []Video `json:"video_list"` // 视频列表
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	hostid := service.GetUserIdByToken(token)
	getuserid, _ := strconv.ParseUint(c.Query("user_id"), 10, 10)
	userid := uint(getuserid)
	getFavoriteList, err := service.FavoriteListGet(userid)
	if err != nil {
		c.JSON(http.StatusNoContent, FavoriteResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "获取点赞列表失败" + err.Error(),
			},
		})
		return
	}
	if len(getFavoriteList) == 0 {
		c.JSON(http.StatusNoContent, FavoriteResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "喜欢列表为空",
			},
			VideoList: nil,
		})
		return
	}
	var favoriteList []Video
	for _, v := range getFavoriteList {
		favoriteUserTemp, _ := service.GetUserInfoById(v.AuthorId)
		favoriteVideoTemp := Video{
			ID: v.ID,
			Author: User{
				ID:            favoriteUserTemp.ID,
				Name:          favoriteUserTemp.Name,
				FollowCount:   favoriteUserTemp.FollowCount,
				FollowerCount: favoriteUserTemp.FollowerCount,
				IsFollow:      service.IsFollowed(favoriteUserTemp.ID, hostid),
			},
			PlayURL:       v.PlayURL,
			CoverURL:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    service.CheckFavorite(hostid, v.ID),
			Title:         v.Title,
		}
		favoriteList = append(favoriteList, favoriteVideoTemp)
	}
	c.JSON(http.StatusOK, FavoriteResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "获取点赞列表成功",
		},
		VideoList: favoriteList,
	})
	return
}

func FavoriteAction(c *gin.Context) {
	// token验证后获取userid
	getToken := c.Query("token")
	userId := service.GetUserIdByToken(getToken)
	// 获取请求参数
	getVideoId := c.Query("video_id")
	actionType := c.Query("action_type")
	// 转换参数类型 string -> uint64
	videoId, err := strconv.ParseUint(getVideoId, 10, 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  "string -> uint64 类型转换失败",
		})
		return
	}
	// 调用处理函数
	// uint64 -> uint
	err = service.FavoriteAction(userId, uint(videoId), actionType)
	var msg string
	switch actionType {
	case "1":
		msg = "成功完成点赞操作"
	case "2":
		msg = "成功完成取消点赞操作"
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  msg,
		})
		return
	}
}
