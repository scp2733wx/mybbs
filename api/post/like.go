package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "mybbs/database"
	SH "mybbs/statehandler"
)

func ClickLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子信息失败", err)
			return
		}

		userID, exsit := c.Get("user_id")
		if !exsit {
			SH.Error(c, http.StatusUnauthorized, "身份信息验证失败", nil)
			return
		}

		isliked := false
		var like Like
		switch db.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error {
		case nil:
			if err := db.DB.Delete(&like).Error; err != nil {
				SH.Error(c, 400, "取消点赞失败", err)
				return
			}
		case gorm.ErrRecordNotFound:
			newLike := Like{
				UserID: userID.(uint),
				PostID: uint(postID),
			}
			if err := db.DB.Create(&newLike).Error; err != nil {
				SH.Error(c, http.StatusNotFound, "点赞失败", err)
				return
			}
			isliked = true
		default:
			SH.Error(c, http.StatusInternalServerError, "数据库查询失败", err)
			return
		}

		SH.Success(c, gin.H{
			"post_id":  postID,
			"is_liked": isliked,
		})
	}
}
