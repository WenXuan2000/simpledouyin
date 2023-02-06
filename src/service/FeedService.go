package service

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
	"time"
)

func FeedGet(lastTime int64) ([]entity.Video, error) {
	if lastTime == 0 {
		lastTime = time.Now().Unix()
	}
	strTime := fmt.Sprint(time.Unix(lastTime, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("查询的时间", strTime)
	var VideoList []entity.Video
	VideoList = make([]entity.Video, 0)
	err := dao.SqlSession.Table("videos").Where("created_at < ?", strTime).Order("created_at desc").Limit(common.VideoFeedNum).Find(&VideoList).Error
	return VideoList, err
}

// ExampleReadFrameAsJpeg 获取封面
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

// CreateVideo 添加一条视频信息
func CreateVideo(video *entity.Video) {
	dao.SqlSession.Table("videos").Create(&video)
}

// 通过作者id 获取所有发布的视频信息
func GetPublicListByAuthorId(uid uint) (videos []entity.Video, err error) {
	err = dao.SqlSession.Model(&entity.Video{}).Where("author_id=?", uid).Order("created_at desc").Find(&videos).Error
	return
}
