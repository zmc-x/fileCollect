package system

import (
	"fileCollect/model/common/response"
	"log"

	"github.com/gin-gonic/gin"
)

// Processing error
func processError(c *gin.Context, msg string, err error) {
	log.Println(msg + err.Error())
	response.Fail(c)
}