package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/middleware"
	"simpledouyin/src/service"
	"strconv"
)

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	uid := service.GetUserIdByToken(token)
	temp, err := strconv.ParseUint(to_user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "类型转换错误string->uint",
			},
		})
		return
	}
	// 获取查询用户uid
	touserid := uint(temp)
	err = service.FollowAction(uid, touserid, action_type)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 2,
			StatusMsg:  err.Error(),
		})
		return
	}
	var msg string
	switch action_type {
	case "1":
		msg = "成功完成关注操作"
	case "2":
		msg = "成功完成取消关注操作"
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  msg,
	})
	return
}

type FollowUser struct {
	FollowCount   int    `json:"follow_count"`   // 关注总数
	FollowerCount int    `json:"follower_count"` // 粉丝总数
	ID            uint   `json:"id"`             // 用户id
	IsFollow      bool   `json:"is_follow"`      // true-已关注，false-未关注
	Name          string `json:"name"`           // 用户名称
}
type FollowUserListResponse struct {
	common.Response
	FollowUserList []FollowUser `json:"user_list"` // 用户信息列表
}

func FollowList(c *gin.Context) {
	quid, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "类型转换错误string->uint",
			},
		})
		return
	}
	queryuid := uint(quid)
	strToken := c.Query("token")
	hostid := middleware.ParseTokenGetID(strToken)
	followlist := make([]FollowUser, 0)
	getfollowlist, err2 := service.FollowListGet(queryuid)
	if err2 != nil {
		c.JSON(http.StatusOK, FollowUserListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "获取关注列表失败",
			},
		})
		return
	}
	for _, user := range getfollowlist {
		followusertemp := FollowUser{
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			ID:            user.ID,
			Name:          user.Name,
			IsFollow:      service.IsFollowed(user.ID, hostid),
		}
		followlist = append(followlist, followusertemp)
	}
	c.JSON(http.StatusOK, FollowUserListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "成功获取关注列表",
		},
		FollowUserList: followlist,
	})
	return
}
