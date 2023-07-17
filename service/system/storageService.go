package system

import (
	"context"
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"
	"fileCollect/model/system/response"
	"fileCollect/utils/cache"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StorageService struct{}

// create the storage
func (s *StorageService) CreateStorage(storageName, storageUrlName, storageRealPath string) error {
	db := global.MysqlDB
	res := db.Create(&model.Storage{
		StorageName:     storageName,
		StorageUrlName:  storageUrlName,
		StorageRealPath: storageRealPath,
		Status:          true,
	})
	// delete cache key-value
	if res.Error == nil {
		rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
		if err := rcache.Del("storageList"); err != nil {
			log.Println("service/system storageService.go CreateStorage method:" + err.Error())
		}
	}
	return res.Error
}

// update the storage's Name
func (s *StorageService) UpdateStorageName(id uint, newName string) error {
	db := global.MysqlDB
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del("storageList"); err != nil {
		log.Println("service/system storageService.go UpdateStorageName method:" + err.Error())
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("StorageName", newName)
	return res.Error
}

// update the storage's url name
func (s *StorageService) UpdateStorageUrlName(id uint, newUrlName string) error {
	db := global.MysqlDB
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del("storageList"); err != nil {
		log.Println("service/system storageService.go UpdateStoragePath method:" + err.Error())
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("StorageUrlName", newUrlName)
	return res.Error
}

// update the storage's path
func (s *StorageService) UpdateStoragePath(id uint, newPath string) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del(fmt.Sprintf("storageId:%d", id)); err != nil {
		log.Println("service/system storageService.go UpdateStoragePath method:" + err.Error())
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("StorageRealPath", newPath)
	return res.Error
}

// update the storage's status
func (s *StorageService) UpdateStorageStatus(id uint, newStatus bool) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del("storageList"); err != nil {
		log.Println("service/system storageService.go UpdateStorageStatus method:" + err.Error())
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("Status", newStatus)
	return res.Error
}

// delete the storage
// files and foldes will be delete if they in this storage
func (s *StorageService) DeleteStorage(id uint) error {
	db := global.MysqlDB
	// delete cache key-value
	rcache := cache.SetRedisStore(context.Background(), 5*time.Minute)
	if err := rcache.Del(fmt.Sprintf("storageId:%d", id)); err != nil {
		log.Println("service/system storageService.go DeleteStorage method:" + err.Error())
	}
	if err := rcache.Del("storageList"); err != nil {
		log.Println("service/system storageService.go DeleteStorage method:" + err.Error())
	}
	storage := model.Storage{
		Model: gorm.Model{ID: id},
	}
	// clear the relation
	// start transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&storage).Association("Files").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&storage).Association("Folders").Clear(); err != nil {
			return err
		}
		// delete all files and all folders in the storage
		if tmp := tx.Where("storage_id is NULL").Delete(&model.File{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		if tmp := tx.Where("storage_id is NULL").Delete(&model.Folder{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		// delete the storage
		if tmp := tx.Where("id = ?", id).Delete(&model.Storage{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		return nil
	})
	return err
}

// query the storage file by storageId and folderId
func (s *StorageService) QueryFiles(storageId, folderId uint) ([]response.StorageFileList, error) {
	// query the file
	var res []response.StorageFileList
	files, folders := []model.File{}, []model.Folder{}
	db := global.MysqlDB
	// query the redis cache
	rcache := cache.SetRedisStore(context.Background(), time.Minute*5)
	key := fmt.Sprintf("files:storageId:%d:folderId:%d", storageId, folderId)
	tRes, err := rcache.Get(key)
	lenTRes := len(tRes)
	// the tRes contain at least 2 element "[]"
	if err == nil && lenTRes > 2 {
		// processing string
		// [{} {}] -> } }
		tmp := strings.ReplaceAll(tRes[1:lenTRes-1], "{", "")
		// start split
		resTmp := strings.Split(tmp[:len(tmp)-1], "}")
		for _, v := range resTmp {
			spliRes := strings.Split(v, " ")
			// traslate
			id, _ := strconv.Atoi(spliRes[0])
			// parse layout: 2006-01-02 15:04:05 -0700 MST
			updateAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", spliRes[1]+" "+spliRes[2]+" "+spliRes[3]+" "+spliRes[4])
			fSize, _ := strconv.Atoi(spliRes[6])
			fType, _ := strconv.ParseBool(spliRes[7])
			res = append(res, response.StorageFileList{
				ID:       uint(id),
				UpdateAt: updateAt,
				FName:    spliRes[5],
				FSize:    uint(fSize),
				FType:    fType,
			})
		}
		return res, err
	}
	// get all files from storage
	if tmp := db.Where("storage_id = ? and folder_id = ?", storageId, folderId).Find(&files); tmp.Error != nil {
		return res, tmp.Error
	}
	if tmp := db.Where("storage_id = ? and parent_folder_id = ?", storageId, folderId).Find(&folders); tmp.Error != nil {
		return res, tmp.Error
	}
	var str string
	// return the result
	for _, v := range files {
		res = append(res, response.StorageFileList{
			ID:       v.ID,
			UpdateAt: v.UpdatedAt,
			FName:    v.FileName,
			FSize:    v.FileSize,
			FType:    response.File,
		})
		str += fmt.Sprintf("{%d %s %s %d %v}", v.ID, v.UpdatedAt.String(), v.FileName, v.FileSize, response.File)
	}
	for _, v := range folders {
		res = append(res, response.StorageFileList{
			ID:       v.ID,
			UpdateAt: v.UpdatedAt,
			FName:    v.FolderName,
			FSize:    0,
			FType:    response.Folder,
		})
		str += fmt.Sprintf("{%d %s %s 0 %v}", v.ID, v.UpdatedAt.String(), v.FolderName, response.Folder)
	}
	// set the str to the cache
	if err := rcache.Set(key, "["+str+"]"); err != nil {
		log.Println("service/system/storageService.go QueryFiles method:" + err.Error())
	}
	return res, nil
}

// query the storage information
func (s *StorageService) QueryStorageInfo() (res []response.StorageInfo, err error) {
	db := global.MysqlDB
	t := []model.Storage{}
	// query the redis cache
	rcache := cache.SetRedisStore(context.Background(), time.Minute*5)
	key := "storageList"
	tRes, tErr := rcache.Get(key)
	lenTRes := len(tRes)
	// contain at least 2 element
	if tErr == nil && lenTRes > 2 {
		err = tErr
		// processing string
		// [{} {}] -> } }
		tmp := strings.ReplaceAll(tRes[1:lenTRes-1], "{", "")
		// start split
		resTmp := strings.Split(tmp[:len(tmp)-1], "}")
		for _, v := range resTmp {
			spliRes := strings.Split(v, " ")
			// traslate
			id, _ := strconv.Atoi(spliRes[0])
			status, _ := strconv.ParseBool(spliRes[4])
			res = append(res, response.StorageInfo{
				Id:             uint(id),
				StorageName:    spliRes[1],
				StorageUrlName: spliRes[2],
				Path:           spliRes[3],
				Status:         status,
			})
		}
		return
	}
	tmp := db.Select("ID", "StorageName", "StorageUrlName", "Status").Find(&t)
	if tmp.RowsAffected == 0 {
		err = errors.New("this system don't have storage")
		return
	}
	var str string
	for _, v := range t {
		res = append(res, response.StorageInfo{
			Id:             v.ID,
			StorageName:    v.StorageName,
			StorageUrlName: v.StorageUrlName,
			// the feild express the storate root catalogue
			Path:   "/",
			Status: v.Status,
		})
		str += fmt.Sprintf("{%d %s %s / %v}", v.ID, v.StorageName, v.StorageUrlName, v.Status)
	}
	// write data to the cache
	if err := rcache.Set(key, "["+str+"]"); err != nil {
		log.Println("service/system/storageService.go QueryStorageInfo:" + err.Error())
	}
	return
}

// query the storage real path
func (s *StorageService) QueryStorageRealPath(id uint) (res string, err error) {
	db := global.MysqlDB
	// query the redis cache
	rcache := cache.SetRedisStore(context.Background(), time.Minute*5)
	// access cache
	key := fmt.Sprintf("storageId:%d", id)
	res, err = rcache.Get(key)
	// Existence of data
	if err == nil && len(res) != 0 {
		return
	}
	t := model.Storage{}
	tmp := db.Select("StorageRealPath").Find(&t)
	if tmp.RowsAffected == 0 {
		err = errors.New("this system don't have storage")
		return
	}
	// write the data to the cache
	if err := rcache.Set(key, t.StorageRealPath); err != nil {
		log.Println("service/system/storageService.go QueryStorageRealPath:" + err.Error())
	}
	// reset the error
	res, err = t.StorageRealPath, nil
	return
}