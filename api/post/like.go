package post

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "mybbs/database"
	SH "mybbs/statehandler"
)

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_post_like" json:"user_id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_user_post_like" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

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
			db.DB.Model(&Post{}).Where("id = ? AND like_count > 0", postID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))

		default:
			newLike := Like{
				UserID: userID.(uint),
				PostID: uint(postID),
			}
			if err := db.DB.Create(&newLike).Error; err != nil {
				SH.Error(c, http.StatusNotFound, "点赞失败", err)
			}
			isliked = true
			db.DB.Model(&Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
		}

		SH.Success(c, gin.H{
			"post_id":  postID,
			"is_liked": isliked,
		})
	}
}
