package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/middleware"
	"simpledouyin/src/service"
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
