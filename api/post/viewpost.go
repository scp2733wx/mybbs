package post

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mybbs/api/file"
	db "mybbs/database"
	SH "mybbs/statehandler"
)

type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	Content   string    `json:"content"`
	Author    auther    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

func VeiwPost() gin.HandlerFunc {
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

		if err := db.DB.Model(&Post{}).Where("id = ?", postID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err == nil {
			post.ViewCount += 1
		}

		var files []file.File
		if err := db.DB.Where("post_id = ?", postID).Find(&files).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "获取文件信息失败", err)
			return
		}
		f_iles := make([]file.FileResponse, 0, len(files))
		for _, f := range files {
			f_iles = append(f_iles, file.FileResponse{
				ID:        f.ID,
				PostID:    f.PostID,
				UserID:    f.UserID,
				FileName:  f.FileName,
				Size:      f.Size,
				UpdatedAt: f.UpdatedAt,
			})
		}

		comments := make([]CommentResponse, 0, len(post.Comments))
		for _, c := range post.Comments {
			comments = append(comments, CommentResponse{
				ID:      c.ID,
				PostID:  c.PostID,
				Content: c.Content,
				Author: auther{
					ID:       c.User.ID,
					Username: c.User.Username,
					Name:     c.User.Name,
					Role:     c.User.Role,
				},
				CreatedAt: c.CreatedAt,
			})
		}

		SH.Success(c, PostResponse{
			ID:      post.ID,
			Content: post.Content,
			Author: auther{
				ID:       post.User.ID,
				Username: post.User.Username,
				Name:     post.User.Name,
				Role:     post.User.Role,
			},
			LikeCount: GetLikeCount(post.ID),
			ViewCount: post.ViewCount,
			CreatedAt: post.CreatedAt,
			Files:     f_iles,
			Comments:  comments,
		})
	}
}
