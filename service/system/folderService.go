package system

import (
	"errors"
	"context"
	"time"
	"fmt"
	"log"
	"fileCollect/utils/cache"
	"fileCollect/global"
	model "fileCollect/model/system"

	"gorm.io/gorm"
)

type FolderService struct{}

// create the folder in the storage
func (s *FolderService) CreateFolder(folderName string, storageId uint, parentFolderId *uint) error {
	db := global.MysqlDB
	var pfi uint 
	// root directory translate
	if parentFolderId == nil {
		pfi = 0
	} else {
		pfi = *parentFolderId
	}
	res := db.Create(&model.Folder{
		FolderName: folderName,
		StorageId: storageId,
		ParentFolderId: pfi,
	})
	if res.Error == nil {
		rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
		if err := rcache.Del(fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, pfi)); err != nil {
			log.Println("service/system/folderService.go CreateFolder methods:" + err.Error())
		}
	}
	return res.Error
}

// update the folder Name
func (s *FolderService) UpdateFolderName(folderId, storageId uint, newName string) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := rcache.Del(fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, folderId)); err != nil {
		log.Println("service/system/folderService.go UpdateFolderName methods:" + err.Error())
	}
	res := db.Model(&model.Folder{}).Where("id = ?", folderId).Update("FolderName", newName)
	return res.Error
}


// delete the folder information 
func (s *FolderService) DeleteFolder(folderId, storageId uint) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5 * time.Minute)
	if err := rcache.Del(fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, folderId)); err != nil {
		log.Println("service/system/folderService.go DeleteFolder methods:" + err.Error())
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Folder{Model : gorm.Model{ID: folderId}}).Association("Files").Clear(); err != nil {
			return err 
		}
		if tmp := tx.Where("folder_id is NULL").Delete(&model.File{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		if tmp := tx.Where("id = ?", folderId).Delete(&model.Folder{}); tmp.Error != nil {
			return tmp.Error
		}
		// commit
		return nil 
	})
	return err
}
