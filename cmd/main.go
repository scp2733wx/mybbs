package main

import (
	"fmt"
	"mybbs/config"
	"mybbs/database"
	"mybbs/router"
	SH "mybbs/statehandler"
)

func main() {
	cfg, err := config.Load()
	SH.PrintError("加载配置文件失败:", err)

	db, err := database.ConnectMySQL(cfg.Database)
	SH.PrintError("连接MySQL失败:", err)
	DB, err := db.DB()
	SH.PrintError("获取数据失败:", err)
	defer DB.Close()
	database.AutoMigrate(db)

	engine := router.InitRouter(db)
	engine.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
