package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpledouyin/src/common"
	"simpledouyin/src/entity"
	"simpledouyin/src/service"
	"strconv"
)

// CommentResponse comment响应结构体
type CommentResponse struct {
	Id      uint              `json:"id"`
	User    UserQueryResponse `json:"user"`
	Content string            `json:"content"`
	Date    string            `json:"create_date"`
}

// CommentActionResponse 发布评论响应结构体
type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment"`
}

// CommentListResponse 获取评论列表响应结构体
type CommentListResponse struct {
	common.Response
	CommentList []CommentResponse `json:"comment_list"`
}

// CommentAction 发布评论处理函数
func CommentAction(c *gin.Context) {
	// 获取参数
	token := c.Query("token")
	hostid := service.GetUserIdByToken(token)
	videoid, _ := strconv.ParseUint(c.Query("video_id"), 10, 10)
	// 判断发布还是删除
	actiontype := c.Query("action_type")
	switch actiontype {
	case "1": // 发布评论
		commentText := c.Query("comment_text")
		AddComment(c, hostid, uint(videoid), commentText)
	case "2": // 删除评论
		commentId, _ := strconv.ParseUint(c.Query("comment_id"), 10, 10)
		DeleteComment(c, uint(videoid), uint(commentId))
	default:
		c.JSON(http.StatusMethodNotAllowed, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "action_type参数错误",
			},
		})
	}
	return
}

func AddComment(c *gin.Context, uid uint, vid uint, content string) {
	// 接收数据库内模型
	newComment := entity.Comment{
		UserId:  uid,
		VideoId: vid,
		Content: content,
	}
	// 发布评论，comments表新增
	if err := service.CommentsCreate(&newComment); err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "comments表create评论失败",
			},
		})
		return
	}
	// videos表字段comment_count+1
	if err := service.VideosCommentCountAdd(vid); err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "videos表更新comment_count失败" + err.Error(),
			},
		})
		return
	}
	// 获取评论发布者的userinfo
	getUserInfo, err1 := service.GetUserInfoById(uid)
	if err1 != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "获取评论发布者的userinfo失败" + err1.Error(),
			},
		})
		return
	}
	// 获取视屏作者的id
	authorId, err2 := service.GetVideoAuthor(vid)
	if err2 != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "获取视频发布者的id失败" + err2.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "发布评论成功",
		},
		Comment: CommentResponse{
			Id: newComment.ID,
			User: UserQueryResponse{
				ID:            getUserInfo.ID,
				FollowCount:   getUserInfo.FollowCount,
				FollowerCount: getUserInfo.FollowerCount,
				Name:          getUserInfo.Name,
				IsFollow:      service.IsFollowed(authorId, uid),
			},
			Content: content,
			Date:    newComment.CreatedAt.Format("01-02"),
		},
	})
	return
}

func DeleteComment(c *gin.Context, vid uint, cid uint) {

	// 删除comments表内数据
	if err := service.CommentsDelete(cid); err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "删除指定comment失败" + err.Error(),
			},
		})
		return
	}
	// videos表更新
	if err := service.VideosCommentCountReduce(vid); err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "videos表更新comment_count失败" + err.Error(),
			},
		})
		return
	}
	// 响应处理
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "删除评论成功",
		},
	})
	return
}

// CommentList 获取评论列表处理函数
func CommentList(c *gin.Context) {
	// 取参数
	token := c.Query("token")
	hostid := service.GetUserIdByToken(token)
	vid, _ := strconv.ParseUint(c.Query("video_id"), 10, 10)
	// 获取comment切片
	getCommentList, err := service.GetCommentListByVideoId(uint(vid))
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "通过video_id获取评论信息失败" + err.Error(),
			},
			CommentList: nil,
		})
		return
	}
	commentList := make([]CommentResponse, 0)
	for _, com := range getCommentList {
		// 通过comment的id获取发布者信息
		user, err := service.GetUserInfoById(com.UserId)
		if err != nil {
			c.JSON(http.StatusOK, CommentListResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "通过comment的id获取发布者信息失败" + err.Error(),
				},
				CommentList: nil,
			})
			return
		}
		comment := CommentResponse{
			Id: com.ID,
			User: UserQueryResponse{
				ID:            user.ID,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				Name:          user.Name,
				IsFollow:      service.IsFollowed(com.ID, hostid),
			},
			Content: com.Content,
			Date:    com.CreatedAt.Format("01-02"),
		}
		commentList = append(commentList, comment)
	}
	// 响应
	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "获取评论成功",
		},
		CommentList: commentList,
	})
	return
}
