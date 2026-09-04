package power

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mybbs/api/post"
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

		if err := db.DB.Where("post_id = ?", postID).Delete(&post.Comment{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除评论失败", err)
			return
		}
		if err := db.DB.Where("post_id = ?", postID).Delete(&post.Like{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除点赞失败", err)
			return
		}
		if err := db.DB.Delete(&post.Post{}, postID).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除帖子失败", err)
			return
		}

		SH.Success(c, "null")
	}
}
