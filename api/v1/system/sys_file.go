package system

import (
	"archive/zip"
	"context"
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"fileCollect/utils/cache"
	"fileCollect/utils/zaplog"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemFileApi struct{}


// router: /api/file/uploadFiles
// method: post
func (sf *SystemFileApi) UploadFiles(c *gin.Context) {
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
		log.Println("api/v/system/sys_file.go UploadFiles method:" + err.Error())
	}
	storageKey, path := info.StorageKey, info.Path
	// query real path according by mysql
	storagePath, err := storageService.QueryStorageRealPath(storageKey)
	if err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	// Check whether the directory exists in the system
	if err := folderService.FolderExist(storageKey, path); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	for _, file := range files {
		fileSrc := filepath.Join(storagePath, path, file.Filename)
		err := c.SaveUploadedFile(file, fileSrc)
		if err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
			continue
		}
		err = fileService.StoreFile(uint(file.Size), file.Filename, storageKey, path)
		if err != nil {
			// delete file if store file to database have error
			defer os.Remove(fileSrc)
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
		}
	}
	response.OkWithMsg(c, "Upload successfully")
}


// router:/api/file/deleteFile/
// method: delete
func (sf *SystemFileApi) DeleteFiles(c *gin.Context) {
	var files request.FileArray
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
	for _, v := range files.Files {
		if err := fileService.DeleteFile(files.StorageKey, v.FileName, files.Path); err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
			continue
		}
		// delete the file
		if err := os.Remove(filepath.Join(storagePath, files.Path, v.FileName)); err != nil {
			zaplog.GetLogLevel(zaplog.WARN, err.Error())
		}
	}
	response.Ok(c)
}

// router:/api/file/updateFile
// method: post
func (sf *SystemFileApi) UpdateFileName(c *gin.Context) {
	var updateNameReq request.UpdateNameReq
	rc := cache.SetRedisStore(context.Background(), 5 * time.Minute)
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
	if err := os.Rename(filepath.Join(pathPre, updateNameReq.FileName), filepath.Join(pathPre, updateNameReq.NewFileName)); err != nil {
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	if err := fileService.UpdateFileName(updateNameReq.StorageKey, updateNameReq.Path, updateNameReq.NewFileName, updateNameReq.FileName); err != nil {
		// restore 
		defer os.Rename(filepath.Join(pathPre, updateNameReq.NewFileName), filepath.Join(pathPre, updateNameReq.FileName))
		zaplog.GetLogLevel(zaplog.ERROR, err.Error())
		response.Fail(c)
		return
	}
	response.Ok(c)
}


// router: /api/file/download
// method: post
func (sf *SystemFileApi) Download(c *gin.Context) {
	var info request.DownMode
	var mark string
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
		mark = filepath.Join(realpath, info.Path)
	}
	lenFile, lenFolder := len(info.Files), len(info.Folders)
	if lenFile > 1 || lenFolder > 0 {
		
		// translate the zip file
		path := []string{}
		var findFile func(src string) 
		findFile = func(src string) {
			dir, _ := os.Stat(src)
			if !dir.IsDir() {
				path = append(path, src)
				return
			}
			files, _ := os.ReadDir(src)
			for _, v := range files {
				findFile(filepath.Join(src, v.Name()))
			}
		}
		zipname := strconv.Itoa(int(time.Now().Unix())) + ".zip"
		zipPath := filepath.Join(mark, zipname)
		zipFile, err := os.Create(zipPath)
		if err != nil {
			zaplog.GetLogLevel(zaplog.ERROR, err.Error())
			response.Fail(c)
			return
		}
		defer os.Remove(zipPath)
		zipWrite := zip.NewWriter(zipFile)
		for _, folder := range info.Folders {
			findFile(filepath.Join(mark, folder))
			err = createZip(zipWrite, path, mark)
			if err != nil {
				zaplog.GetLogLevel(zaplog.ERROR, err.Error())
				response.Fail(c)
				return
			}
			path = nil
		}
		for _, file := range info.Files {
			path = append(path, filepath.Join(mark, file))
		}
		if err = createZip(zipWrite, path, mark); err != nil {
			zaplog.GetLogLevel(zaplog.ERROR, err.Error())
			response.Fail(c)
			return
		}
		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(zipname)))
		// close file
		zipWrite.Close()
		zipFile.Close()
		c.File(zipPath)
	} else {
		// single file
		filename := info.Files[0]
		src := filepath.Join(mark, filename)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			zaplog.GetLogLevel(zaplog.ERROR, err.Error())
			response.Fail(c)
			return
		}
		// set header
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(filename)))
		// Serve the file for download
		c.File(src)
	}
}

// create zip archive
func createZip(zipWrite *zip.Writer, path []string, prefix string) (err error){
	for _, v := range path {
		zipName, err := filepath.Rel(prefix, v)
		if err != nil {
			return err
		}
		dstf, err := zipWrite.Create(zipName)
		if err != nil {
			return err
		}
		srcf, err := os.Open(v)
		if err != nil {
			return err
		}
		defer srcf.Close()
		_, err = io.Copy(dstf, srcf)
		if err != nil {
			return err
		}
	}
	return nil
}