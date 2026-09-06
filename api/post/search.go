package post

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	db "mybbs/database"
	SH "mybbs/statehandler"
)

func SearchList() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if err != nil || pageSize < 1 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}

		baseQuery := db.DB.Model(&Post{}).Joins("JOIN users ON users.id = posts.user_id")

		if q != "" {
			keyword := "%" + q + "%"
			baseQuery = baseQuery.Where(
				"posts.content LIKE ? OR users.username LIKE ? OR users.name LIKE ?",
				keyword, keyword, keyword,
			)
		}

		var total int64
		if err := baseQuery.Count(&total).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "获取搜索结果总数失败", err)
			return
		}

		var posts []Post
		if err := baseQuery.Preload("User").Order("posts.created_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&posts).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "搜索帖子失败", err)
			return
		}

		likeMap := make(map[uint]int, len(posts))
		if len(posts) > 0 {
			postIDs := make([]uint, 0, len(posts))
			for _, p := range posts {
				postIDs = append(postIDs, p.ID)
			}
			var likeRows []struct {
				PostID uint
				Count  int64
			}
			db.DB.Model(&Like{}).
				Select("post_id, COUNT(*) as count").
				Where("post_id IN ?", postIDs).
				Group("post_id").
				Scan(&likeRows)
			for _, r := range likeRows {
				likeMap[r.PostID] = int(r.Count)
			}
		}

		items := make([]PostResponse, 0, len(posts))
		for _, p := range posts {
			items = append(items, PostResponse{
				ID:           p.ID,
				Content:      p.Content,
				LikeCount:    likeMap[p.ID],
				ViewCount:    p.ViewCount,
				CommentCount: int64(len(p.Comments)),
				CreatedAt:    p.CreatedAt,
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
				"q":         q,
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		})
	}
}
