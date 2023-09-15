package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func RspOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "",
		Data: data,
	})
}
func RspError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: -1,
		Msg:  msg,
	})
}