package post

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	db "mybbs/database"
	SH "mybbs/statehandler"
)

type auther struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type PostResponse struct {
	ID           uint              `json:"id"`
	Content      string            `json:"content"`
	Author       auther            `json:"author"`
	LikeCount    int               `json:"like_count"`
	CommentCount int64             `json:"comment_count"`
	ViewCount    int               `json:"view_count"`
	CreatedAt    time.Time         `json:"created_at"`
	Comments     []CommentResponse `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

func GetPostList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if err != nil {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}

		var total int64
		if err := db.DB.Model(&Post{}).Count(&total).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "获取帖子总数失败", err)
			return
		}

		var rule string
		switch c.Query("mode") {
		case "hot":
			rule = "nil"
		default:
			rule = "created_at desc"
		}

		var posts []Post
		if err := db.DB.Model(&Post{}).Preload("User").Order(rule).Limit(pageSize).Offset((page - 1) * pageSize).Find(&posts).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "获取帖子列表失败", err)
			return
		}

		items := make([]PostResponse, 0, len(posts))
		for _, p := range posts {
			items = append(items, PostResponse{
				ID:        p.ID,
				Content:   p.Content,
				LikeCount: GetLikeCount(p.ID),
				ViewCount: p.ViewCount,
				CreatedAt: p.CreatedAt,
				Author: auther{
					ID:       p.User.ID,
					Username: p.User.Username,
					Name:     p.User.Name,
					Role:     p.User.Role,
				},
			})
		}

		SH.Success(c, gin.H{
			"items": items,
			"meta": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		})
	}
}
