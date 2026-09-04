package router

import (
	"github.com/gin-gonic/gin"

	power "mybbs/api/admin"
	"mybbs/api/post"
	"mybbs/api/user"
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
	protected.GET("/posts", post.GetPostList())
	protected.GET("/posts/:post_id", post.VeiwPost())
	protected.DELETE("/posts/:post_id", post.DeletePost())
	protected.POST("/posts/:post_id/like", post.ClickLike())
	protected.POST("/posts/likes", post.GetLikes())
	protected.POST("/posts/:post_id/comment", post.CreatComment())

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	admin.DELETE("/posts/:post_id", power.DeletePost())
	return engine
}
