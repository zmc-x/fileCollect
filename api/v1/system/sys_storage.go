package system

import (
	"context"
	"encoding/json"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"fileCollect/utils/cache"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemStorageApi struct{}

// router:/api/storage/createStorage
// method:post
func (s *SystemStorageApi) CreateStorage(c *gin.Context) {
	var info request.StorageInfo
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := c.ShouldBindJSON(&info); err != nil {
		processError(c, "api/v1/system/sys_storage.go CreateStorage method:", err)
		return
	}
	// store
	if err := storageService.CreateStorage(info.StorageName, info.StorageURLName, info.StorageRealPath); err != nil {
		processError(c, "api/v1/system/sys_storage.go CreateStorage method:", err)
		return
	}
	if err := rc.Del("storageInfo"); err != nil {
		log.Println("api/v1/system/sys_storage.go CreateStorage method:" + err.Error())
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
	if err := storageService.UpdateStorageStatus(model.StorageKey, model.NewStatus); err != nil {
		processError(c, "api/v1/system/sys_storage.go UpdateStorageStatus method:", err)
		return
	}
	response.Ok(c)
}

// update generic function model
func storageUpdateModel(c *gin.Context, param string, updatefunc func(storageKey, newN string) error) {
	var model request.UpdateGeneric
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
	if err := updatefunc(model.StorageKey, new); err != nil {
		processError(c, "api/v1/system/sys_storage.go storageUpdateModel function:", err)
		return
	}
	if err := rc.Del("storageInfo"); err != nil {
		log.Println("api/v1/system/sys_storage.go storageUpdateModel function:" + err.Error())
	}
	response.Ok(c)
}

// router:/api/storage/delete/{storageKey}
// method:delete
func (s *SystemStorageApi) DeleteStorage(c *gin.Context) {
	storageKey := c.Param("storageKey")
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := storageService.DeleteStorage(storageKey); err != nil {
		processError(c, "api/v1/system/sys_storage.go DeleteStorage method:", err)
		return
	}
	if err := rc.Del("storageInfo"); err != nil {
		log.Println("api/v1/system/sys_storage.go DeleteStorage method:" + err.Error())
	}
	response.Ok(c)
}

// router:/api/storage/query/storageInfo
// method:get
func (s *SystemStorageApi) QueryStorageInfo(c *gin.Context) {
	data := response.StorageInfo{}
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	// storageList
	if str, err := rc.Get("storageInfo"); err == nil {
		json.Unmarshal([]byte(str), &data.StorageList)
		response.OkWithData(c, data)
		return
	}
	t, err := storageService.QueryStorageInfo()
	if err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryStorageInfo method:", err)
		return
	}
	for _, v := range t {
		data.StorageList = append(data.StorageList, response.StorageList{
			StorageName: v.StorageName,
			Status: v.Status,
			StorageKey: v.StorageKey,
			Path: "/",
		})
	}
	// data.StorageList []struct -> json
	if tmp, err := json.Marshal(data.StorageList); err != nil {
		log.Println("api/v1/system/sys_storage.go QueryStorageInfo method:" + err.Error())
	} else {
		rc.Set("storageInfo", string(tmp))
	}
	response.OkWithData(c, data)
}


// router:/api/storage/query/list
// method:post
func (s *SystemStorageApi) QueryFilesList(c *gin.Context) {
	var req request.ReqStorageList
	var res response.FilesInfo
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := c.ShouldBindJSON(&req); err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryFilesList method:", err)
		return 
	}
	// FileList:{storageKey}:{path}
	if str, err := rc.Get("FileList:" + req.StorageKey + ":" + req.Path); err == nil {
		json.Unmarshal([]byte(str), &res.FileList)
		response.OkWithData(c, res)
		return
	}
	// call the storageService
	if t, err := storageService.QueryFiles(req.StorageKey, req.Path); err != nil {
		processError(c, "api/v1/system/sys_storage.go QueryFilesList method:", err)
		return 
	} else {
		for _, v := range t {
			res.FileList = append(res.FileList, response.FileList{
				FName: v.FName,
				FSize: v.FSize,
				FType: v.FType,
				UpdateAt: v.UpdateAt.Format("2006-01-02 15:04:05"),
			})
		}
	}
	if temp, err := json.Marshal(res.FileList); err != nil {
		log.Println("api/v1/system/sys_storage.go QueryFilesList method:" + err.Error())
	} else {
		rc.Set("FileList:" + req.StorageKey + ":" + req.Path, string(temp))
	}
	response.OkWithData(c, res)
}