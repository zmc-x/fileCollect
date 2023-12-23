package system

import (
	"context"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"fileCollect/utils/cache"
	"fileCollect/utils/zaplog"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemFileApi struct{}

// router: /api/file/uploadFiles
// method: post
func (sf *SystemFileApi) UploadFiles(c *gin.Context) {
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	form, err := c.MultipartForm()
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	files := form.File["uploads"]
	// get other information
	var info request.FileInfo
	if err := c.ShouldBind(&info); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := rc.Del("FileList:" + info.StorageKey + ":" + info.Path); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	storageKey, path := info.StorageKey, info.Path
	// query real path according by mysql
	storagePath, err := storageService.QueryStorageRealPath(storageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	fileLen := len(files)
	for _, file := range files {
		fileSrc := filepath.Join(storagePath, path, file.Filename)
		err := fileService.StoreFile(file, fileSrc)
		if err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
			continue
		}
		fileLen--
	}
	if fileLen != 0 {
		zaplog.GetLogLevel(zaplog.WARN, "Some files failed to be uploaded")
		response.FailWithMsg(c, "Some files failed to be uploaded")
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "upload files successfully")
	response.OkWithMsg(c, "Upload successfully")
}

// router:/api/file/deleteFile/
// method: delete
func (sf *SystemFileApi) DeleteFiles(c *gin.Context) {
	var files request.FileArray
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := c.ShouldBindJSON(&files); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := rc.Del("FileList:" + files.StorageKey + ":" + files.Path); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	// get storage real path
	storagePath, err := storageService.QueryStorageRealPath(files.StorageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	fileLen := len(files.Files)
	for _, file := range files.Files {
		if err := fileService.DeleteFile(storagePath, filepath.Join(files.Path, file.FileName)); err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
			continue
		}
		fileLen--
	}
	if fileLen != 0 {
		zaplog.GetLogLevel(zaplog.WARN, "Some files failed to be deleted")
		response.FailWithMsg(c, "Some files failed to be deleted")
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "delete files successfully")
	response.Ok(c)
}

// router:/api/file/updateFile
// method: post
func (sf *SystemFileApi) UpdateFileName(c *gin.Context) {
	var updateNameReq request.UpdateNameReq
	rc := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := c.ShouldBindJSON(&updateNameReq); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := rc.Del("FileList:" + updateNameReq.StorageKey + ":" + updateNameReq.Path); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
	}
	// get storage real path
	storagePath, err := storageService.QueryStorageRealPath(updateNameReq.StorageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	// rename file name
	// file path prefix
	pathPre := filepath.Join(storagePath, updateNameReq.Path)
	if err := fileService.UpdateFileName(pathPre, updateNameReq.FileName, updateNameReq.NewFileName); err != nil {
		zaplog.GetLogLevel(zaplog.WARN, err.Error())
		response.Fail(c)
		return
	}
	zaplog.GetLogLevel(zaplog.INFO, "updateFileName successfully")
	response.Ok(c)
}

// router: /api/file/download
// method: post
func (sf *SystemFileApi) Download(c *gin.Context) {
	var info request.DownMode
	var path string
	if err := c.ShouldBindJSON(&info); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if realpath, err := storageService.QueryStorageRealPath(info.StorageKey); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	} else {
		path = filepath.Join(realpath, info.Path)
	}
	lenFile, lenFolder := len(info.Files), len(info.Folders)
	if lenFile > 1 || lenFolder > 0 {
		// translate the zip file
		zipPath, zipFile, comp, err := fileService.DownloadCompressFile(info.Files, info.Folders, path)
		defer os.Remove(zipPath)
		if err != nil {
			zaplog.GetLogLevel(zaplog.ERROR, err.Error())
			response.Fail(c)
			comp.Close()
			zipFile.Close()
			return
		}
		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(filepath.Base(zipPath))))
		// close file
		comp.Close()
		zipFile.Close()
		c.File(zipPath)
	} else {
		// single file
		src, err := fileService.Download(info.Files, path)
		if err != nil {
			zaplog.GetLogLevel(zaplog.ERROR, err.Error())
			response.Fail(c)
			return
		}
		// set header
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(filepath.Base(src))))
		// Serve the file for download
		c.File(src)
	}
	zaplog.GetLogLevel(zaplog.INFO, "download files successfully")
}


