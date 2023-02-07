package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/service"
	"strconv"
)

func FavoriteList(c *gin.Context) {

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
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "操作成功",
		})
		return
	}
}
