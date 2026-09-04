package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "mybbs/database"
	SH "mybbs/statehandler"
)

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		userID, exsit := c.Get("user_id")
		if post.UserID != userID.(uint) || !exsit {
			SH.Error(c, http.StatusNonAuthoritativeInfo, "身份信息验证失败", nil)
			return
		}

		if err := db.DB.Where("post_id = ?", postID).Delete(&Comment{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除评论失败", err)
			return
		}
		if err := db.DB.Where("post_id = ?", postID).Delete(&PostLike{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除点赞失败", err)
			return
		}
		if err := db.DB.Delete(&Post{}, postID).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除帖子失败", err)
			return
		}

		SH.Success(c, "null")
	}
}
