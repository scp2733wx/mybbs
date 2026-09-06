package main

import (
	"fmt"
	"log"
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

	engine := router.InitRouter()
	addr := fmt.Sprintf(":%d", config.CFG.Server.Port)

	if config.CFG.TLS.Enabled {
		log.Printf("HTTPS 服务启动: https://127.0.0.1%s", addr)
		err = engine.RunTLS(addr, config.CFG.TLS.Cert, config.CFG.TLS.Key)
	} else {
		log.Printf("HTTP 服务启动: http://127.0.0.1%s", addr)
		err = engine.Run(addr)
	}
	SH.PrintError("服务运行失败:", err)
}
