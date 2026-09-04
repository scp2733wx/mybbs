package statehandler

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"mybbs/config"
)

var File *os.File

func InitLogger() error {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return err
	}

	File, err = os.OpenFile(config.CFG.Server.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(io.MultiWriter(os.Stdout, File))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	gin.DefaultWriter = io.MultiWriter(os.Stdout, File)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, File)
	return nil
}

func AddLog(err error) {
	logFile, _ := os.OpenFile(config.CFG.Server.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logFile.WriteString(fmt.Sprintf("[SH-ERROR] | %s | %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
	defer logFile.Close()
}
