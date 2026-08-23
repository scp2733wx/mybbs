package power

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mybbs/api/post"
	SH "mybbs/statehandler"
)

func DeletePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子信息失败："+err.Error())
			return
		}

		if err := db.Where("post_id = ?", postID).Delete(&post.Comment{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除评论失败："+err.Error())
			return
		}
		if err := db.Where("post_id = ?", postID).Delete(&post.PostLike{}).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除点赞失败："+err.Error())
			return
		}
		if err := db.Delete(&post.Post{}, postID).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "删除帖子失败："+err.Error())
			return
		}

		SH.Success(c, "null")
	}
}
