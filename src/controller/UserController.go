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

//func Login(c *gin.Context) {
//	//定义一个User变量
//	var user entity.User
//	//将调用后端的request请求中的body数据根据json格式解析到User结构变量中
//	c.BindJSON(&user)
//	//将被转换的user变量传给service层的CreateUser方法，进行User的新建
//	err := service.CreateUser(&user)
//	//判断是否异常，无异常则返回包含200和更新数据的信息
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//	} else {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 200,
//			"msg":  "success",
//			"data": user,
//		})
//	}
//}
