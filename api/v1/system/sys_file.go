package system

import (
	"context"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"fileCollect/utils/cache"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemFileApi struct{}


// router: /api/file/uploadFiles
// method: post
func (sf *SystemFileApi) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err != nil {
		processError(c, "api/v1/system/sys_file.go UploadFiles method:", err)
		return
	}
	files := form.File["uploads"]
	// get other information
	var info request.FileInfo
	if err := c.ShouldBind(&info); err != nil {
		processError(c, "api/v1/system/sys_file.go UploadFiles method:", err)
		return
	}
	storageKey, path := info.StorageKey, info.Path
	// query real path according by redis or mysql
	storagePath, err := storageService.QueryStorageRealPath(storageKey)
	if err != nil {
		processError(c, "api/v1/system/sys_file.go UploadFiles method:", err)
		return
	}
	// Check whether the directory exists in the system
	if err := folderService.FolderExist(storageKey, path); err != nil {
		processError(c, "api/v1/system/sys_file.go UploadFiles method:", err)
		return
	}
	for _, file := range files {
		fileSrc := filepath.Join(storagePath, path, file.Filename)
		err := c.SaveUploadedFile(file, fileSrc)
		if err != nil {
			log.Println("api/v1/system/sys_file.go UploadFiles method:" + err.Error())
			continue
		}
		err = fileService.StoreFile(uint(file.Size), file.Filename, storageKey, path)
		if err != nil {
			// delete file if store file to database have error
			defer os.Remove(fileSrc)
			log.Println("api/v1/system/sys_file.go UploadFiles method:" + err.Error())
		}
	}
	if err := rc.Del("FileList:" + info.StorageKey + ":" + info.Path); err != nil {
		log.Println("api/v/system/sys_file.go UploadFiles method:" + err.Error())
	}
	response.OkWithMsg(c, "Upload successfully")
}


// router:/api/file/deleteFile/
// method: delete
func (sf *SystemFileApi) DeleteFiles(c *gin.Context) {
	var files request.FileArray
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := c.ShouldBindJSON(&files); err != nil {
		processError(c, "api/v1/system/sys_file.go DeleteFiles method:", err)
		return
	}
	// get storage real path
	storagePath, err := storageService.QueryStorageRealPath(files.StorageKey)
	if err != nil {
		processError(c, "api/v1/system/sys_file.go DeleteFiles method:", err)
		return
	}
	for _, v := range files.Files {
		if err := fileService.DeleteFile(files.StorageKey, v.FileName, files.Path); err != nil {
			log.Println("api/v1/system/sys_file.go DeleteFiles method:" + err.Error())
			continue
		}
		// delete the file
		if err := os.Remove(filepath.Join(storagePath, files.Path, v.FileName)); err != nil {
			log.Println("api/v1/system/sys_file.go DeleteFiles method:" + err.Error())
		}
	}
	if err := rc.Del("FileList:" + files.StorageKey + ":" + files.Path); err != nil {
		log.Println("api/v1/system/sys_file.go DeleteFiles method:" + err.Error())
	}
	response.Ok(c)
}

// router:/api/file/updateFile
// method: post
func (sf *SystemFileApi) UpdateFileName(c *gin.Context) {
	var updateNameReq request.UpdateNameReq
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := c.ShouldBindJSON(&updateNameReq); err != nil {
		processError(c, "api/v1/system/sys_file.go UpdateFileName method:", err)
		return
	}
	// get storage real path
	storagePath, err := storageService.QueryStorageRealPath(updateNameReq.StorageKey)
	if err != nil {
		processError(c, "api/v1/system/sys_file.go UpdateFileName method:", err)
		return
	}
	// rename file name
	// file path prefix
	pathPre := filepath.Join(storagePath, updateNameReq.Path)
	if err := os.Rename(filepath.Join(pathPre, updateNameReq.FileName), filepath.Join(pathPre, updateNameReq.NewFileName)); err != nil {
		processError(c, "api/v1/system/sys_file.go UpdateFileName method:", err)
		return
	}
	if err := fileService.UpdateFileName(updateNameReq.StorageKey, updateNameReq.Path, updateNameReq.NewFileName, updateNameReq.FileName); err != nil {
		// restore 
		defer os.Rename(filepath.Join(pathPre, updateNameReq.NewFileName), filepath.Join(pathPre, updateNameReq.FileName))
		processError(c, "api/v1/system/sys_file.go UpdateFileName method:", err)
		return
	}
	if err := rc.Del("FileList:" + updateNameReq.StorageKey + ":" + updateNameReq.Path); err != nil {
		log.Println("api/v1/system/sys_file.go UpdateFileName method:" + err.Error())
	}
	response.Ok(c)
}