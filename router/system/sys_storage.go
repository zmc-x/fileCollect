package system

import (
	v1 "fileCollect/api/v1"

	"github.com/gin-gonic/gin"
)

type StorageApi struct{}

func (s *StorageApi) InitStorageRouter(rg *gin.RouterGroup) {
	storageRouter := rg.Group("storage")
	storageApi := v1.ApiGroupApp.SystemApiGroup.SystemStorageApi
	{
		storageRouter.POST("createStorage", storageApi.CreateStorage)
	}
	update := storageRouter.Group("update")
	{
		update.POST("storageName", storageApi.UpdateStorageName)
		update.POST("storageUrlName", storageApi.UpdateStorageUrl)
		update.POST("realPath", storageApi.UpdateStoragePath)
		update.POST("status", storageApi.UpdateStorageStatus)
	}
	query := storageRouter.Group("query")
	{
		query.GET("storageInfo", storageApi.QueryStorageInfo)
		query.GET("list", storageApi.QueryFilesList)
	}
}