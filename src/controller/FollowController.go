package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
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
