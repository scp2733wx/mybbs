package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	SH "mybbs/statehandler"
)

func CreatPost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var content struct {
			Content string `json:"content" binding:"required,min=1,max=2000"`
		}
		if err := c.ShouldBindJSON(&content); err != nil {
			SH.Error(c, http.StatusBadRequest, "获取帖子内容失败："+err.Error())
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			SH.Error(c, http.StatusUnauthorized, "未认证")
			return
		}

		post := Post{
			Content: content.Content,
			UserID:  userID.(uint),
		}
		if err := db.Create(&post).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "发布失败："+err.Error())
			return
		}

		SH.Success(c, post)
	}
}
