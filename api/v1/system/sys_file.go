package system

import (
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type SystemFileApi struct{}


// router: /api/file/uploadFiles
// method: post
func (sf *SystemFileApi) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
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
	folderId, storageId, path := info.FolderID, info.StorageID, info.Path
	// query real path according by redis or mysql
	storagePath, err := storageService.QueryStorageRealPath(storageId)
	if err != nil {
		processError(c, "api/v1/system/sys_file.go UploadFiles method:", err)
		return
	}
	for _, file := range files {
		err := c.SaveUploadedFile(file, filepath.Join(storagePath, path, file.Filename))
		if err != nil {
			log.Println("api/v1/system/sys_file.go UploadFiles method:" + err.Error())
			continue
		} 
		switch path {
		case "/": 
			err = fileService.StoreFile(uint(file.Size), storageId, file.Filename, nil)
		default: 
			err = fileService.StoreFile(uint(file.Size), storageId, file.Filename, &folderId)
		}
		// don't return
		if err != nil {
			// delete file if store file to database have error
			defer os.Remove(filepath.Join(storagePath, path, file.Filename))
			log.Println("api/v1/system/sys_file.go UploadFiles method:" + err.Error())
		}
	}
	response.OkWithMsg(c, "Upload successfully")
}


// router:/api/file/deleteFile/
// method: delete
func (sf *SystemFileApi) DeleteFiles(c *gin.Context) {
	var files request.FileArray
	if err := c.ShouldBindJSON(&files); err != nil {
		processError(c, "api/v1/system/sys_file.go DeleteFiles method:", err)
		return
	}
	// get storage real path
	storagePath, err := storageService.QueryStorageRealPath(files.StorageId)
	if err != nil {
		processError(c, "api/v1/system/sys_file.go DeleteFiles method:", err)
		return
	}
	for _, v := range files.Files {
		if err := fileService.DeleteFile(v.FileID, files.StorageId, files.ParentFolderId); err != nil {
			log.Println("api/v1/system/sys_file.go DeleteFiles method:" + err.Error())
			continue
		}
		// delete the file
		if err := os.Remove(filepath.Join(storagePath, files.Path, v.FileName)); err != nil {
			log.Println("api/v1/system/sys_file.go DeleteFiles method:" + err.Error())
		}
	}
	response.Ok(c)
}

// router:/api/file/updateFile
// method: post
func (sf *SystemFileApi) UpdateFileName(c *gin.Context) {
	var updateNameReq request.UpdateNameReq
	if err := c.ShouldBindJSON(&updateNameReq); err != nil {
		processError(c, "api/v1/system/sys_file.go UpdateFileName method:", err)
		return
	}
	// get storage real path
	storagePath, err := storageService.QueryStorageRealPath(updateNameReq.StorageId)
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
	if err := fileService.UpdateFileName(updateNameReq.FileID, updateNameReq.StorageId, updateNameReq.ParentFolderId, updateNameReq.NewFileName); err != nil {
		// restore 
		defer os.Rename(filepath.Join(pathPre, updateNameReq.NewFileName), filepath.Join(pathPre, updateNameReq.FileName))
		processError(c, "api/v1/system/sys_file.go UpdateFileName method:", err)
		return
	}
	response.Ok(c)
}