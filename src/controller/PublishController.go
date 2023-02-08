package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"simpledouyin/src/common"
	"simpledouyin/src/entity"
	"simpledouyin/src/service"
	"strconv"
	"strings"
	"time"
)

func Publish(c *gin.Context) {
	//中间件验证token后，获取userId
	token := c.PostForm("token")
	userId := service.GetUserIdByToken(token)
	//接收请求参数信息
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//设置文件名称
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%d_%s", userId, time.Now().Unix(), fileName)
	playurl := "http://10.0.2.2:8080/static/" + finalName
	//存储到本地文件夹
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 封面路径
	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
	coverurl := "http://10.0.2.2:8080/static/" + coverName
	saveImage := filepath.Join("./public/", coverName)

	buf := service.ExampleReadFrameAsJpeg(saveFile, 3) //获取第3帧封面
	img, _ := jpeg.Decode(buf)                         //保存到本地时要用到
	imgw, _ := os.Create(saveImage)                    //先创建，后写入
	jpeg.Encode(imgw, img, &jpeg.Options{100})

	//保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	video := entity.Video{
		Model:         gorm.Model{},
		AuthorId:      userId,
		PlayURL:       playurl,
		CoverURL:      coverurl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	service.CreateVideo(&video)

	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "视频投稿成功",
	})
	return

}

// PublishListResponse
type PublishListResponse struct {
	common.Response
	VideoList []Video `json:"video_list"` // 视频列表
}

// 用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(c *gin.Context) {
	// 通过token获取hostid
	token := c.Query("token")
	hostid := service.GetUserIdByToken(token)
	// userid
	getuserid, _ := strconv.ParseUint(c.Query("user_id"), 10, 10)
	userid := uint(getuserid)
	// 得到视频
	videolist, err := service.GetPublicListByAuthorId(userid)
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.VideoGetWrong.Error(),
			},
		})
		return
	}
	if len(videolist) == 0 {
		c.JSON(http.StatusNoContent, PublishListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "该用户发布列表为空",
			},
			VideoList: nil,
		})
		return
	}
	// 得到userid的userinfo
	user, err := service.GetUserInfoById(userid)
	auther := User{
		ID:            user.ID,
		FollowerCount: user.FollowerCount,
		FollowCount:   user.FollowCount,
		Name:          user.Name,
		IsFollow:      false,
	}
	if userid != hostid {
		auther.IsFollow = service.IsFollowed(userid, hostid)
	}
	publishvideolist := make([]Video, 0)
	for _, video := range videolist {
		videotemp := Video{
			ID:            video.ID,
			Author:        auther,
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    service.CheckFavorite(hostid, video.ID),
			Title:         video.Title,
		}
		publishvideolist = append(publishvideolist, videotemp)
	}
	c.JSON(http.StatusOK, PublishListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "成功显示用户发布所有视频",
		},
		VideoList: publishvideolist,
	})
	return
}
