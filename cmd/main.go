package main

import (
	"fmt"
	"mybbs/api/post"
	"mybbs/api/user"
	"mybbs/config"
	db "mybbs/database"
	"mybbs/router"
	SH "mybbs/statehandler"
)

func main() {
	err := config.Load()
	SH.PrintError("加载配置文件失败:", err)

	err = SH.InitLogger()
	SH.PrintError("初始化日志失败：", err)
	defer SH.File.Close()

	err = db.ConnectMySQL(config.CFG.Database)
	SH.PrintError("连接MySQL失败:", err)
	DB, err := db.DB.DB()
	SH.PrintError("获取数据失败:", err)
	defer DB.Close()
	db.DB.AutoMigrate(
		&user.User{},
		&post.Post{},
		&post.Comment{},
		&post.PostLike{},
	)

	engine := router.InitRouter()
	engine.Run(fmt.Sprintf(":%d", config.CFG.Server.Port))
}
