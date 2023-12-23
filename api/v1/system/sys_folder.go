package system

import (
	"context"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"fileCollect/utils/cache"
	"fileCollect/utils/zaplog"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemFolderApi struct{}

// router:/api/folder/createFolder
// method:post
func (sf *SystemFolderApi) CreateFolder(c *gin.Context) {
	var createFolderInfo request.CreateFolderInfo
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := c.ShouldBindJSON(&createFolderInfo); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := rc.Del("FileList:" + createFolderInfo.StorageKey + ":" + createFolderInfo.Path); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	// get the storageRealPath
	storagePath, err := storageService.QueryStorageRealPath(createFolderInfo.StorageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	// create folder in system
	err = folderService.CreateFolder(storagePath, filepath.Join(createFolderInfo.Path, createFolderInfo.FolderName))
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "create folder successfully")
	response.Ok(c)
}

// router:/api/folder/deletefolder/
// method:delete
func (sf *SystemFolderApi) DeleteFolders(c *gin.Context) {
	var info request.DeleteFolderInfo
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := c.ShouldBindJSON(&info); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := rc.Del("FileList:" + info.StorageKey + ":" + info.Path); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	// get the storageRealPath
	storagePath, err := storageService.QueryStorageRealPath(info.StorageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	folderLen := len(info.Folders)
	for _, folder := range info.Folders {
		if err := folderService.DeleteFolder(storagePath, filepath.Join(info.Path, folder.FolderName)); err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
			continue
		}
		folderLen--
	}
	if folderLen != 0 {
		zaplog.GetLogLevel(zaplog.WARN, "Some folders failed to be deleted")
		response.FailWithMsg(c, "Some folders failed to be deleted")
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "delete folder successfully")
	response.Ok(c)
}

// router:/api/folder/updatefolder
// method:post
func (sf *SystemFolderApi) UpdateFolder(c *gin.Context) {
	var info request.UpdateFolderInfo
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := c.ShouldBindJSON(&info); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := rc.Del("FileList:" + info.StorageKey + ":" + info.Path); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	// get the storageRealPath
	storagePath, err := storageService.QueryStorageRealPath(info.StorageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	folderPre := filepath.Join(storagePath, info.Path)
	// update system folder
	err = folderService.UpdateFolderName(folderPre, info.FolderName, info.FolderNewName)
	if err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "update folder information successfully")
	response.Ok(c)
}
