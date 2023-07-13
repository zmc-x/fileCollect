package system

import (
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"log"

	"github.com/gin-gonic/gin"
)

type SystemStorageApi struct{}

// router:/api/storage/createStorage
// method:post
func (s *SystemStorageApi) CreateStorage(c *gin.Context) {
	var info request.StorageInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	// store
	if err := storageService.CreateStorage(info.StorageName, info.StorageURLName, info.StorageRealPath); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	response.Ok(c)
}

// router:/api/storage/update/storageName
// method:post
func (s *SystemStorageApi) UpdateStorageName(c *gin.Context) {
	storageUpdateModel(c, "newName", storageService.UpdateStorageName)
}

// router:/api/storage/update/storageUrlName
// method:post
func (s *SystemStorageApi) UpdateStorageUrl(c *gin.Context) {
	storageUpdateModel(c, "newUrlName", storageService.UpdateStorageUrlName)
}

// router:/api/storage/update/realPath
// method:post
func (s *SystemStorageApi) UpdateStoragePath(c *gin.Context) {
	storageUpdateModel(c, "newPath", storageService.UpdateStoragePath)
}

// router:/api/storage/update/status
// method:post
func (s *SystemStorageApi) UpdateStorageStatus(c *gin.Context) {
	var model request.UpdateGeneric
	if err := c.ShouldBindJSON(&model); err != nil {
		log.Println(err)
		response.Fail(c)
		return 
	}
	if err := storageService.UpdateStorageStatus(model.StorageID, model.NewStatus); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	response.Ok(c)
}
// update generic function model
func storageUpdateModel(c *gin.Context, param string, updatefunc func(id uint, new string) error) {
	var model request.UpdateGeneric
	if err := c.ShouldBindJSON(&model); err != nil {
		log.Println(err)
		response.Fail(c)
		return 
	}
	var new string 
	switch param {
	case "newName":
		new = model.NewName
	case "newUrlName":
		new = model.NewUrlName
	case "newPath":
		new = model.NewPath
	}
	log.Println(new, model)
	if err := updatefunc(model.StorageID, new); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	response.Ok(c)
}