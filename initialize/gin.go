package initialize

import (
	"fileCollect/router"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	// set the upload file the max size
	r.MaxMultipartMemory = 1 << 30	// 1G
	// init routers
	systemRouters := router.RouterGroupApp.SystemRouter
	systemApi := r.Group("api")
	{
		systemRouters.InitFileApiRouter(systemApi)
		systemRouters.InitStorageRouter(systemApi)
		systemRouters.InitFolderRouter(systemApi)
	}
	return r
}