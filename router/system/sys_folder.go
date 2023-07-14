package system

import (
	v1 "fileCollect/api/v1"

	"github.com/gin-gonic/gin"
)

type FolderApi struct{}

func (f *FolderApi) InitFolderRouter(rg *gin.RouterGroup) {
	folderApis := v1.ApiGroupApp.SystemApiGroup.SystemFolderApi
	folderRouter := rg.Group("folder")
	{
		folderRouter.POST("updatefolder", folderApis.UpdateFolder)
		folderRouter.POST("createFolder", folderApis.CreateFolder)
		folderRouter.DELETE("deletefolder", folderApis.DeleteFolders)
	}
}