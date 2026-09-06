package router

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	power "mybbs/api/admin"
	"mybbs/api/file"
	"mybbs/api/post"
	"mybbs/api/user"
	"mybbs/config"
	"mybbs/middleware"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()

	api := engine.Group("/api/v1")
	api.POST("/auth/register", user.Register())
	api.POST("/auth/login", user.Login())

	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	protected.POST("/posts", post.CreatPost())
	protected.POST("/posts/:post_id/assents", file.UpLoad())
	protected.GET("/posts", post.GetPostList())
	protected.GET("/posts/search", post.SearchList())
	protected.GET("/posts/:post_id", post.VeiwPost())
	protected.DELETE("/posts/:post_id", post.DeletePost())
	protected.POST("/posts/:post_id/like", post.ClickLike())
	protected.POST("/posts/likes", post.GetLikes())
	protected.POST("/posts/:post_id/comment", post.CreatComment())

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	admin.DELETE("/posts/:post_id", power.DeletePost())

	mountFrontend(engine)
	return engine
}

// mountFrontend 动态挂载前端静态资源。
// 仅当配置启用且目录存在时挂载，否则跳过。
// 支持 SPA：未匹配的路径回退到 index.html。
func mountFrontend(engine *gin.Engine) {
	cfg := config.CFG.Server.Frontend
	if !cfg.Enabled || cfg.Path == "" {
		return
	}

	root, err := filepath.Abs(cfg.Path)
	if err != nil {
		return
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return
	}

	indexFile := filepath.Join(root, "index.html")
	if _, err := os.Stat(indexFile); err != nil {
		return
	}

	engine.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}

		rel := strings.TrimPrefix(c.Request.URL.Path, "/")
		if rel == "" {
			rel = "index.html"
		}

		target := filepath.Join(root, filepath.Clean("/"+rel))
		if !strings.HasPrefix(target, root) {
			c.Next()
			return
		}

		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			c.File(target)
			c.Abort()
			return
		}

		// SPA 路由回退到 index.html
		c.File(indexFile)
		c.Abort()
	})
}
