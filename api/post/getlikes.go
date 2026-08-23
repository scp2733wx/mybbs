package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	SH "mybbs/statehandler"
)

type LikeStatus struct {
	PostID uint `json:"post_id"`
	Liked  bool `json:"liked"`
}

func GetLikes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exsit := c.Get("user_id")
		if !exsit {
			SH.Error(c, http.StatusUnauthorized, "身份信息验证失败")
			return
		}

		var Likes struct {
			Likes   []Like `json:"likes"`
			PostIDs []uint `json:"post_ids"`
		}

		if err := c.ShouldBindJSON(&Likes); err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子列表错误："+err.Error())
			return
		}

		if err := db.Where("user_id = ? AND post_id IN ?", userID, Likes.PostIDs).Find(&Likes.Likes).Error; err != nil {
			SH.Error(c, http.StatusBadRequest, "查找点赞列表错误："+err.Error())
			return
		}

		islikeds := make(map[uint]bool, len(Likes.Likes))
		for _, l := range Likes.Likes {
			islikeds[uint(l.PostID)] = true
		}

		status := make([]LikeStatus, 0, len(Likes.PostIDs))
		for _, p := range Likes.PostIDs {
			status = append(status, LikeStatus{
				PostID: p,
				Liked:  islikeds[p],
			})
		}

		SH.Success(c, gin.H{"status": status})
	}
}
