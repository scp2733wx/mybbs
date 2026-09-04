package statehandler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func JSON(c *gin.Context, httpStatus, code int, msg string, data any) {
	c.JSON(httpStatus, Envelope{Code: code, Msg: msg, Data: data})
}

func Success(c *gin.Context, data any) {
	JSON(c, http.StatusOK, 0, "Success", data)
	c.Abort()
}

func Error(c *gin.Context, code int, msg string, err error) {
	JSON(c, code, code, msg, nil)
	AddLog(err)
	c.Abort()
}

func PrintError(txt string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, txt)
		AddLog(err)
		os.Exit(1)
	}
}
