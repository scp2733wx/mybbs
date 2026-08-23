package post

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	SH "mybbs/statehandler"
)

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_post_like" json:"user_id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_user_post_like" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

func ClickLike(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子信息失败："+err.Error())
			return
		}

		userID, exsit := c.Get("user_id")
		if !exsit {
			SH.Error(c, http.StatusUnauthorized, "身份信息验证失败")
			return
		}

		isliked := false
		var like Like
		switch db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error {
		case nil:
			if err := db.Delete(&like).Error; err != nil {
				SH.Error(c, 400, "取消点赞失败："+err.Error())
				return
			}
			db.Model(&Post{}).Where("id = ? AND like_count > 0", postID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
		default:
			newLike := Like{
				UserID: userID.(uint),
				PostID: uint(postID),
			}
			if err := db.Create(&newLike).Error; err != nil {
				SH.Error(c, http.StatusNotFound, "点赞失败："+err.Error())
			}
			isliked = true
			db.Model(&Post{}).Where("id = ? AND like_count > 0", postID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
		}

		SH.Success(c, gin.H{
			"post_id":  postID,
			"is_liked": isliked,
		})
	}
}
