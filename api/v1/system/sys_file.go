package system

import (
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type SystemFileApi struct{}


// router: /api/file/uploadFiles
// method: post
func (sf *SystemFileApi) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("parsed multipartForm error, the error is " + err.Error())
		response.Fail(c)
		return
	}
	files := form.File["uploads"]
	// get other information
	var info request.FileInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	storagePath, folderId, storageId, path := info.StorageRealPath, info.FolderID, info.StorageID, info.Path
	for _, file := range files {
		err := c.SaveUploadedFile(file, filepath.Join(storagePath, path, file.Filename))
		if err != nil {
			log.Println(err)
			response.Fail(c)
			return 
		} 
		switch path {
		case "/": 
			err = fileService.StoreFile(uint(file.Size), storageId, file.Filename, nil)
		default: 
			err = fileService.StoreFile(uint(file.Size), storageId, file.Filename, &folderId)
		}
		if err != nil {
			log.Println(err)
		}
	}
	response.OkWithMsg(c, "Upload successfully")
}


// router:/api/file/deleteFile/
// method: delete
func (sf *SystemFileApi) DeleteFiles(c *gin.Context) {
	var files request.FileId
	if err := c.ShouldBindJSON(&files); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	for _, v := range files.Files {
		if err := fileService.DeleteFile(v); err != nil {
			log.Println(err)
		}
	}
	response.Ok(c)
}

// router:/api/file/updateFile
// method: post
func (sf *SystemFileApi) UpdateFileName(c *gin.Context) {
	var updateNameReq request.UpdateNameReq
	if err := c.ShouldBindJSON(&updateNameReq); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	if err := fileService.UpdateFileName(updateNameReq.FileID, updateNameReq.NewFileName); err != nil {
		log.Println(err)
		response.Fail(c)
		return
	}
	response.Ok(c)
}