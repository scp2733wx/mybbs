package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mybbs/api/file"
	db "mybbs/database"
	SH "mybbs/statehandler"
)

func SerchList() gin.HandlerFunc {
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
			var files []file.File
			if err := db.DB.Model(&file.File{}).Where("post_id = ?", p.ID).Find(&files).Error; err != nil {
				SH.Error(c, http.StatusInternalServerError, "获取帖子附件失败", err)
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

			items = append(items, PostResponse{
				ID:           p.ID,
				Content:      p.Content,
				LikeCount:    GetLikeCount(p.ID),
				ViewCount:    p.ViewCount,
				CommentCount: int64(len(p.Comments)),
				CreatedAt:    p.CreatedAt,
				Files:        f_iles,
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
