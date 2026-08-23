package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	SH "mybbs/statehandler"
)

func DeletePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子信息失败："+err.Error())
			return
		}

		var post Post
		if err := db.Preload("User").Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at asc")
		}).Preload("Comments.User").First(&post, postID).Error; err != nil {
			SH.Error(c, http.StatusNotFound, "帖子不存在："+err.Error())
			return
		}
		userID, exsit := c.Get("user_id")
		if post.UserID != userID.(uint) || !exsit {
			SH.Error(c, http.StatusNonAuthoritativeInfo, "身份信息验证失败")
			return
		}

		if err := db.Where("post_id = ?", postID).Delete(&Comment{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除评论失败："+err.Error())
			return
		}
		if err := db.Where("post_id = ?", postID).Delete(&PostLike{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除点赞失败："+err.Error())
			return
		}
		if err := db.Delete(&Post{}, postID).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除帖子失败："+err.Error())
			return
		}

		SH.Success(c, "null")
	}
}
