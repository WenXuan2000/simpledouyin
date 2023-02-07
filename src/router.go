package main

import (
	"github.com/gin-gonic/gin"
	"simpledouyin/src/controller"
	"simpledouyin/src/middleware"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.JwtHandler(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.JwtHandler(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.JwtHandler(), controller.PublishList)
	//
	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JwtHandler(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.JwtHandler(), controller.FavoriteList)
	//apiRouter.POST("/comment/action/", middleware.JwtHandler(), controller.CommentAction)
	//apiRouter.GET("/comment/list/", middleware.JwtHandler(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JwtHandler(), controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", middleware.JwtHandler(), controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", middleware.JwtHandler(), controller.FollowerList)
	//apiRouter.GET("/relation/friend/list/", middleware.JwtHandler(), controller.FriendList)
	//apiRouter.GET("/message/chat/", middleware.JwtHandler(), controller.MessageChat)
	//apiRouter.POST("/message/action/", middleware.JwtHandler(), controller.MessageAction)
}
