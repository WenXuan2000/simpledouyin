package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/middleware"
	"simpledouyin/src/service"
	"strconv"
)

type UidTokenResponse struct {
	common.Response
	Uid   uint   `json:"user_id"`
	Token string `json:"token"`
}

// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	err := service.IsUserPasswordLegal(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UidTokenResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	uid, err := service.UserRegister(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UidTokenResponse{
			Response: common.Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	token, err := middleware.CreateToken(uid)
	if err != nil {
		c.JSON(http.StatusOK, UidTokenResponse{
			Response: common.Response{
				StatusCode: 3,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UidTokenResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		Uid:   uid,
		Token: token,
	})
	return
}

// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 判断用户名存不存在
	if !service.IsUserNameUnique(username) {
		c.JSON(http.StatusOK, UidTokenResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.UserNoexist.Error(),
			},
		})
		return
	}
	// 判断用户名和密码对不对应
	uid, err := service.IsPasswordRight(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UidTokenResponse{
			Response: common.Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	token, err := middleware.CreateToken(uid)
	if err != nil {
		c.JSON(http.StatusOK, UidTokenResponse{
			Response: common.Response{
				StatusCode: 3,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UidTokenResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		Uid:   uid,
		Token: token,
	})
	return
}

type UserInfoResponse struct {
	common.Response
	User UserQueryResponse `json:"user"` // 用户信息
}

// User
type UserQueryResponse struct {
	ID            uint   `json:"id"`             // 用户id
	FollowCount   int    `json:"follow_count"`   // 关注总数
	FollowerCount int    `json:"follower_count"` // 粉丝总数
	Name          string `json:"name"`           // 用户名称
	IsFollow      bool   `json:"is_follow"`      // true-已关注，false-未关注
}

// 用户信息
func UserInfo(c *gin.Context) {
	// 获取请求的uid
	uqid := c.Query("user_id")
	// 请求的用户信息存不存在
	if !service.IsUserLiveById(uqid) {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.UserLiveWrong.Error(),
			},
		})
		return
	}
	// 获得请求的用户信息
	//var userq = entity.User{}
	temp, err := strconv.ParseUint(uqid, 10, 64)
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
	uqiduint := uint(temp)
	userinfo, err := service.GetUserInfoById(uqiduint)
	// 获取自己的token
	clientoken := c.Query("token")
	// 获取自己的uid
	uid := service.GetUserIdByToken(clientoken)
	// 计算uid是否关注uqiduint
	isfollow := service.IsFollowed(uqiduint, uid)

	c.JSON(http.StatusOK, UserInfoResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "请求用户信息成功",
		},
		User: UserQueryResponse{
			ID:            userinfo.ID,
			FollowCount:   userinfo.FollowCount,
			FollowerCount: userinfo.FollowerCount,
			Name:          userinfo.Name,
			IsFollow:      isfollow,
		},
	})
	return

}
