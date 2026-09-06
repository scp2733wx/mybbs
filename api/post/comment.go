package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "mybbs/database"
	SH "mybbs/statehandler"
)

func CreatComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var content struct {
			Content string `json:"content" binding:"required,min=1,max=2000"`
		}
		if err := c.ShouldBindJSON(&content); err != nil {
			SH.Error(c, http.StatusBadRequest, "获取评论内容失败", err)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			SH.Error(c, http.StatusUnauthorized, "未认证", nil)
			return
		}

		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子信息失败", err)
			return
		}

		var post Post
		if err := db.DB.Preload("User").Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at asc")
		}).Preload("Comments.User").First(&post, postID).Error; err != nil {
			SH.Error(c, http.StatusNotFound, "帖子不存在", err)
			return
		}

		comment := Comment{
			Content: content.Content,
			UserID:  userID.(uint),
			PostID:  uint(postID),
		}
		if err := db.DB.Create(&comment).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "发布评论失败", err)
			return
		}

		if err := db.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "用户信息登记失败", err)
			return
		}

		SH.Success(c, comment)
	}
}

func GetCommentCount(postID uint) int64 {
	var count int64
	db.DB.Model(&Comment{}).Where("post_id = ?", postID).Count(&count)
	return count
}
