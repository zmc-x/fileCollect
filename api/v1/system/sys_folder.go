package system

import (
	"context"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"fileCollect/utils/cache"
	"fileCollect/utils/zaplog"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemFolderApi struct{}

// router:/api/folder/createFolder
// method:post
func (sf *SystemFolderApi) CreateFolder(c *gin.Context) {
	var createFolderInfo request.CreateFolderInfo
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
	err = os.Mkdir(filepath.Join(storagePath, createFolderInfo.Path, createFolderInfo.FolderName), 0644)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	// database
	err = folderService.CreateFolder(createFolderInfo.FolderName, createFolderInfo.StorageKey, createFolderInfo.Path)
	if err != nil {
		defer os.Remove(filepath.Join(storagePath, createFolderInfo.Path, createFolderInfo.FolderName))
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	response.Ok(c)
}


// router:/api/folder/deletefolder/
// method:delete
func (sf *SystemFolderApi) DeleteFolders(c *gin.Context) {
	var info request.DeleteFolderInfo
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
	for _, v := range info.Folders {
		if err := folderService.DeleteFolder(v.FolderName, info.Path, info.StorageKey); err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
			continue
		}
		// delete system folder
		if err := os.RemoveAll(filepath.Join(storagePath, info.Path, v.FolderName)); err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
		}
	}
	response.Ok(c)
}

// router:/api/folder/updatefolder
// method:post
func (sf *SystemFolderApi) UpdateFolder(c *gin.Context) {
	var info request.UpdateFolderInfo
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
	nName, oName := filepath.Join(folderPre, info.FolderNewName), filepath.Join(folderPre, info.FolderName)
	// update system folder
	if err := os.Rename(oName, nName); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := folderService.UpdateFolderName(info.FolderName, info.Path, info.StorageKey, info.FolderNewName); err != nil {
		// restore
		defer os.Rename(nName, oName)
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	response.Ok(c)
}