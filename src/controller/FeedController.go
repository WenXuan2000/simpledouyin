package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/middleware"
	"simpledouyin/src/service"
	"strconv"
	"time"
)

// Video
type Video struct {
	ID            uint   `json:"id"`             // 视频唯一标识
	Author        User   `json:"author"`         // 视频作者信息
	PlayURL       string `json:"play_url"`       // 视频播放地址
	CoverURL      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int    `json:"favorite_count"` // 视频的点赞总数
	CommentCount  int    `json:"comment_count"`  // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title         string `json:"title"`          // 视频标题
}

// User
type User struct {
	ID            uint   `json:"id"`             // 用户id
	Name          string `json:"name"`           // 用户名称
	FollowCount   int    `json:"follow_count"`   // 关注总数
	FollowerCount int    `json:"follower_count"` // 粉丝总数
	IsFollow      bool   `json:"is_follow"`      // true-已关注，false-未关注
}

// feedResponse
type FeedResponse struct {
	common.Response
	VideoList []Video `json:"video_list"` // 视频列表
	NextTime  int     `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

func Feed(c *gin.Context) {
	strToken := c.Query("token")
	var haveToken bool
	if strToken == "" {
		haveToken = false
	} else {
		haveToken = true
	}
	strLastTime := c.Query("lastest_time")
	lastTime, err := strconv.ParseInt(strLastTime, 10, 32)
	if err != nil {
		lastTime = 0
	}

	VideoList := make([]Video, 0)
	getVideoList, _ := service.FeedGet(lastTime)
	var newTime int64 = 0 //返回的视频的最久的一个的时间
	for _, x := range getVideoList {
		var tmp Video
		tmp.ID = x.ID
		tmp.PlayURL = x.PlayURL
		//tmp.Author = //依靠用户信息接口查询
		var user, err = service.GetUserInfoById(x.AuthorId)
		var feedUser User
		if err == nil { //用户存在
			feedUser.ID = user.ID
			feedUser.FollowerCount = user.FollowerCount
			feedUser.FollowCount = user.FollowCount
			feedUser.Name = user.Name
			feedUser.IsFollow = false
			if haveToken {
				// 查询是否关注
				_, tokenStruct, _ := middleware.ParseToken(strToken)
				if time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
					var uid1 = tokenStruct.Uid //用户id
					var uid2 = x.AuthorId      //视频发布者id
					if service.IsFollowed(uid1, uid2) {
						feedUser.IsFollow = true
					}
				}
			}
		}
		tmp.Author = feedUser
		tmp.CommentCount = x.CommentCount
		tmp.CoverURL = x.CoverURL
		tmp.FavoriteCount = x.FavoriteCount
		tmp.IsFavorite = false
		//if haveToken {
		//	//查询是否点赞过
		//	_, tokenStruct, _ := middleware.ParseToken(strToken)
		//	if time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
		//		var uid = tokenStruct.Uid            //用户id
		//		var vid = x.ID                       // 视频id
		//		if service.CheckFavorite(uid, vid) { //有点赞记录
		//			tmp.IsFavorite = true
		//		}
		//	}
		//}
		tmp.Title = x.Title
		VideoList = append(VideoList, tmp)
		newTime = x.CreatedAt.Unix()
	}
	if len(VideoList) > 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  common.Response{StatusCode: 0}, //成功
			VideoList: VideoList,
			NextTime:  int(newTime),
		})
	}

	//else {
	//	c.JSON(http.StatusOK, FeedNoVideoResponse{
	//		Response: common.Response{StatusCode: 0}, //成功
	//		NextTime: 0,                              //重新循环
	//	})
	//}

}
