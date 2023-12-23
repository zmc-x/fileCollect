package system

import (
	"context"
	"encoding/json"
	"fileCollect/global"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	sysres "fileCollect/model/system/response"
	"fileCollect/utils/cache"
	"fileCollect/utils/zaplog"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemStorageApi struct{}

// router:/api/storage/createStorage
// method:post
func (s *SystemStorageApi) CreateStorage(c *gin.Context) {
	var info request.StorageInfo
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rc.Del("storageInfo"); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	tempValidTime, err := time.Parse(global.Format, info.DeadLine+"(CST)")
	if err != nil {
		tempValidTime, _ = time.Parse(global.Format, "9999-01-01 00:00:00(CST)")
	}
	// store
	if err := storageService.CreateStorage(info.StorageName, info.StorageURLName, info.StorageRealPath, tempValidTime); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "create storage successfully")
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
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rc.Del("storageInfo"); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	if err := c.ShouldBindJSON(&model); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	t, err := time.Parse(global.Format, model.DeadLine+"(CST)")
	if err != nil {
		t, _ = time.Parse(global.Format, "9999-01-01 00:00:00(CST)")
	}
	if err := storageService.UpdateStorageStatus(model.StorageKey, model.NewStatus, t); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "update storage status successfully")
	response.Ok(c)
}

// update generic function model
func storageUpdateModel(c *gin.Context, param string, updatefunc func(storageKey, newN string) error) {
	var model request.UpdateGeneric
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rc.Del("storageInfo"); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	if err := c.ShouldBindJSON(&model); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
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
	if err := updatefunc(model.StorageKey, new); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "update storage information successfully")
	response.Ok(c)
}

// router:/api/storage/delete/{storageKey}
// method:delete
func (s *SystemStorageApi) DeleteStorage(c *gin.Context) {
	storageKey := c.Param("storageKey")
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rc.Del("storageInfo"); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	if err := storageService.DeleteStorage(storageKey); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "delete storage successfully")
	response.Ok(c)
}

// router:/api/storage/query/storageInfo
// method:get
func (s *SystemStorageApi) QueryStorageInfo(c *gin.Context) {
	data := sysres.Storages{}
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	// storageList
	if str, err := rc.Get("storageInfo"); err == nil {
		json.Unmarshal([]byte(str), &data.StorageList)
		response.OkWithData(c, data)
		return
	}
	data, err := storageService.QueryStorageInfo()
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	// data.StorageList []struct -> json
	if tmp, err := json.Marshal(data.StorageList); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	} else {
		rc.Set("storageInfo", string(tmp))
	}
	zaplog.GetLogLevel(zaplog.INFO, "query storage list successfully")
	response.OkWithData(c, data)
}

// router:/api/storage/query/list
// method:post
func (s *SystemStorageApi) QueryFilesList(c *gin.Context) {
	var req request.ReqStorageList
	var res sysres.Files
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := c.ShouldBindJSON(&req); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	// FileList:{storageKey}:{path}
	if str, err := rc.Get("FileList:" + req.StorageKey + ":" + req.Path); err == nil {
		json.Unmarshal([]byte(str), &res.FileList)
		response.OkWithData(c, res)
		return
	}
	// call the storageService
	res, err := storageService.QueryFiles(req.StorageKey, req.Path)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if temp, err := json.Marshal(res.FileList); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	} else {
		rc.Set("FileList:"+req.StorageKey+":"+req.Path, string(temp))
	}
	zaplog.GetLogLevel(zaplog.INFO, "query storage file successfully")
	response.OkWithData(c, res)
}
