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

type ChatMessage struct {
	Id         uint   `json:"id,omitempty"`
	FromUserId uint   `json:"from_user_id"`
	ToUserId   uint   `json:"to_user_id"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}
type DouYinChatMessageResponse struct {
	common.Response
	ChatMessages []ChatMessage `json:"message_list"`
}

func MessageChat(c *gin.Context) {
	token := c.Query("token")
	to_user_id, _ := strconv.Atoi(c.Query("to_user_id"))
	uid := middleware.ParseTokenGetID(token)
	MessageList, err := service.GetMessageList(uid, uint(to_user_id))
	if err != nil {
		c.JSON(http.StatusOK, DouYinChatMessageResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "获取消息列表失败",
			},
		})
		return
	}
	var ChatMessages []ChatMessage
	for _, message := range MessageList {
		ChatMessages = append(ChatMessages, ChatMessage{
			Id:         message.ID,
			FromUserId: message.UserID,
			ToUserId:   message.ToUserID,
			Content:    message.Content,
			CreateTime: message.CreateTime,
		})
	}
	c.JSON(http.StatusOK, DouYinChatMessageResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "成功获得消息列表",
		},
		ChatMessages: ChatMessages,
	})
	return
}

func MessageAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id, _ := strconv.Atoi(c.Query("to_user_id"))
	action_type := c.Query("action_type")
	content := c.Query("content")
	createtime := time.Now().Format(time.Kitchen)
	uid := middleware.ParseTokenGetID(token)

	switch action_type {
	case "1":
		err := service.SendMessage(uid, uint(to_user_id), createtime, content)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}

	default:
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "没有定义的操作",
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "成功发送消息",
	})
	return
}
