package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int	`json:"code"`
	Data interface{}	`json:"data"`
	Msg  string	`json:"msg"`
}

const (
	SUCCESS = 0
	ERROR   = -1
)

func result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg: msg,
	})
}

func Ok(c *gin.Context) {
	result(SUCCESS, map[string]interface{}{}, "Successful operation", c)
}

func OkWithData(c *gin.Context, data interface{}) {
	result(SUCCESS, data, "Query success", c)
}

func OkWithMsg(c *gin.Context, msg string) {
	result(SUCCESS, map[string]interface{}{}, msg, c)
}

func Fail(c *gin.Context) {
	result(ERROR, map[string]interface{}{}, "Operation failure", c)
}

func FailWithData(c *gin.Context, data interface{}) {
	result(ERROR, data, "Query fail", c)
}

func FailWithMsg(c *gin.Context, msg string) {
	result(ERROR, map[string]interface{}{}, msg, c)
}