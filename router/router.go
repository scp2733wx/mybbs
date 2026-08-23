package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	power "mybbs/api/admin"
	"mybbs/api/post"
	"mybbs/api/user"
	"mybbs/middleware"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	engine := gin.Default()

	api := engine.Group("/api/v1")
	api.POST("/auth/register", user.Register(db))
	api.POST("/auth/login", user.Login(db))

	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	protected.POST("/posts", post.CreatPost(db))
	protected.GET("/posts", post.GetPostList(db))
	protected.GET("/posts/:post_id", post.VeiwPost(db))
	protected.DELETE("/posts/:post_id", post.DeletePost(db))
	protected.POST("/posts/:post_id/like", post.ClickLike(db))
	protected.POST("/posts/likes", post.GetLikes(db))
	protected.POST("/posts/:post_id/comment", post.CreatComment(db))

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	admin.DELETE("/posts/:post_id", power.DeletePost(db))
	return engine
}
