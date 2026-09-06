package file

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mybbs/config"
	db "mybbs/database"
	SH "mybbs/statehandler"
)

func UpLoad() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			SH.Error(c, 400, "获取文件失败", err)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			SH.Error(c, 401, "未认证", nil)
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			SH.Error(c, 400, "获取帖子信息失败", err)
			return
		}

		Fpath := filepath.Join(config.CFG.Server.Assentpath, fmt.Sprintf("%s_%d_%d", time.Now().Format("2004-01-02_15-04-05"), userID.(uint), postID))

		new_file := File{
			FileName: file.Filename,
			PostID:   uint(postID),
			UserID:   userID.(uint),
			FilePath: Fpath,
			Size:     file.Size,
		}

		if err := c.SaveUploadedFile(file, Fpath); err != nil {
			SH.Error(c, 500, "保存文件失败", err)
			return
		}

		if err := db.DB.Create(&new_file).Error; err != nil {
			SH.Error(c, 500, "保存文件信息失败", err)
			return

		}
		SH.Success(c, new_file)
	}
}
