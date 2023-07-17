package system

import (
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SystemStorageApi struct{}

// router:/api/storage/createStorage
// method:post
func (s *SystemStorageApi) CreateStorage(c *gin.Context) {
	var info request.StorageInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		processError(c, "api/v1/system/sys_storage.go CreateStorage method:", err)
		return
	}
	// store
	if err := storageService.CreateStorage(info.StorageName, info.StorageURLName, info.StorageRealPath); err != nil {
		processError(c, "api/v1/system/sys_storage.go CreateStorage method:", err)
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
		processError(c, "api/v1/system/sys_storage.go UpdateStorageStatus method:", err)
		return 
	}
	if err := storageService.UpdateStorageStatus(model.StorageID, model.NewStatus); err != nil {
		processError(c, "api/v1/system/sys_storage.go UpdateStorageStatus method:", err)
		return
	}
	response.Ok(c)
}

// update generic function model
func storageUpdateModel(c *gin.Context, param string, updatefunc func(id uint, new string) error) {
	var model request.UpdateGeneric
	if err := c.ShouldBindJSON(&model); err != nil {
		processError(c, "api/v1/system/sys_storage.go storageUpdateModel function:", err)
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
	if err := updatefunc(model.StorageID, new); err != nil {
		processError(c, "api/v1/system/sys_storage.go storageUpdateModel function:", err)
		return
	}
	response.Ok(c)
}

// router:/api/storage/delete/{storageId}
// method:delete
func (s *SystemStorageApi) DeleteStorage(c *gin.Context) {
	tStorageId := c.Param("storageId")
	storageId, err := strconv.Atoi(tStorageId)
	if err != nil {
		processError(c, "api/v1/system/sys_storage.go DeleteStorage method:", err)
		return
	}
	if err := storageService.DeleteStorage(uint(storageId)); err != nil {
		processError(c, "api/v1/system/sys_storage.go DeleteStorage method:", err)
		return
	}
	response.Ok(c)
}

// router:/api/storage/query/storageInfo
// method:get
func (s *SystemStorageApi) QueryStorageInfo(c *gin.Context) {
	t, err := storageService.QueryStorageInfo()
	if err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryStorageInfo method:", err)
		return
	}
	data := response.StorageInfo{}
	for _, v := range t {
		data.StorageList = append(data.StorageList, response.StorageList{
			StorageID: v.Id,
			StorageName: v.StorageName,
			StorageURLName: v.StorageUrlName,
			Status: v.Status,
			Path: v.Path,
		})
	}
	response.OkWithData(c, data)
}


// router:/api/storage/query/list
// method:get
func (s *SystemStorageApi) QueryFilesList(c *gin.Context) {
	// get query params
	tStorageId, tFolderId := c.Query("storageId"), c.Query("folderId")
	storageId, err := strconv.Atoi(tStorageId)
	if err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryFilesList method:", err)
		return
	}
	folderId, err := strconv.Atoi(tFolderId)
	if err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryFilesList method:", err)
		return
	}
	t, err := storageService.QueryFiles(uint(storageId), uint(folderId))
	if err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryFilesList method:", err)
		return
	}
	data := response.FilesInfo{}
	for _, v := range t {
		data.FileList = append(data.FileList, response.FileList{
			FileID: v.ID,
			FName: v.FName,
			FSize: v.FSize,
			FType: v.FType,
			UpdateAt: v.UpdateAt,
		})
	}
	response.OkWithData(c, data)
}