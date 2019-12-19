package util

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func Resp(c *gin.Context, httpCode, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}

func SuccessResp(c *gin.Context, message string) {
	if message == "" {
		message = "ok"
	}
	Resp(c, 0, 0, message)
}

func ErrorResp(c *gin.Context, httpCode, code int, message string) {
	Resp(c, httpCode, code, message)
}
