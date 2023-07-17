package system

import (
	"context"
	"fileCollect/global"
	model "fileCollect/model/system"
	"fileCollect/utils/cache"
	"fmt"
	"log"
	"time"
)

type FileService struct{}

// store files' information into database
func (s *FileService) StoreFile(fileSize, storageId uint, fileName string, folderId *uint) error {
	db := global.MysqlDB
	var tmp uint
	// root directory translate
	if folderId == nil {
		tmp = 0
	} else {
		tmp = *folderId
	}
	res := db.Create(&model.File{
		FileSize:  fileSize,
		FolderId:  tmp,
		StorageId: storageId,
		FileName:  fileName,
	})
	if res.Error == nil {
		// delete cache key-value
		rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
		if err := rcache.Del(fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, tmp)); err != nil {
			log.Println("service/system/fileService.go StorageFile methods:" + err.Error())
		}
	}
	return res.Error
}

// update file related information
func (s *FileService) UpdateFileName(fileId, storageId, parentFolderId uint, newName string) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del(fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, parentFolderId)); err != nil {
		log.Println("service/system/fileService.go UpdateFileName methods:" + err.Error())
	}
	res := db.Model(&model.File{}).Where("id = ?", fileId).Update("FileName", newName)
	return res.Error
}

// delete file record
func (s *FileService) DeleteFile(fileId, storageId, parentFolderId uint) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del(fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, parentFolderId)); err != nil {
		log.Println("service/system/fileService.go DeleteFile methods:" + err.Error())
	}
	res := db.Where("id = ?", fileId).Delete(&model.File{})
	return res.Error
}
