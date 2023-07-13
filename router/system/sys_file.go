package system

import (
	v1 "fileCollect/api/v1"

	"github.com/gin-gonic/gin"
)

type FileApi struct{}

func (f *FileApi) InitFileApiRouter(rg *gin.RouterGroup) {
	fileRouter := rg.Group("file")
	fileApis := v1.ApiGroupApp.SystemApiGroup
	{
		fileRouter.POST("uploadFiles", fileApis.UploadFiles)
		fileRouter.DELETE("deleteFile", fileApis.DeleteFiles)
		fileRouter.POST("updateFile", fileApis.UpdateFileName)
	}
}