package system

import (
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"
	"fileCollect/model/system/response"
	"time"

	"gorm.io/gorm"
)

type StorageService struct{}

type validPair struct {
	StorageUrlName string
	DeadLine       time.Time
}

var (
	storageSet      map[validPair]bool = make(map[validPair]bool)
	storagelist     []validPair
	storagedeadline chan validPair = make(chan validPair)
	done            chan struct{}  = make(chan struct{})
	// preserve resource
	token chan struct{} = make(chan struct{}, 1)
)

// create the storage
func (s *StorageService) CreateStorage(storageName, storageUrlName, storageRealPath string, deadline time.Time) error {
	db := global.MysqlDB
	res := db.Create(&model.Storage{
		StorageName:     storageName,
		StorageUrlName:  storageUrlName,
		StorageRealPath: storageRealPath,
		Status:          true,
		DeadLine:        deadline,
	})
	// start timer
	if res.Error == nil {
		storageSet[validPair{StorageUrlName: storageUrlName, DeadLine: deadline}] = true
		go validTimer(validPair{StorageUrlName: storageUrlName, DeadLine: deadline})
	}
	return res.Error
}

// update the storage's Name
func (s *StorageService) UpdateStorageName(storageKey, newName string) error {
	db := global.MysqlDB
	res := db.Model(&model.Storage{}).Where("storage_url_name = ?", storageKey).Update("StorageName", newName)
	return res.Error
}

// update the storage's url name
func (s *StorageService) UpdateStorageUrlName(storageKey, newUrlName string) error {
	db := global.MysqlDB
	res := db.Model(&model.Storage{}).Where("storage_url_name = ?", storageKey).Update("StorageUrlName", newUrlName)
	return res.Error
}

// update the storage's path
func (s *StorageService) UpdateStoragePath(storageKey, newPath string) error {
	db := global.MysqlDB
	// delete the global realPath
	delete(global.RealPath, storageKey)
	res := db.Model(&model.Storage{}).Where("storage_url_name = ?", storageKey).Update("StorageRealPath", newPath)
	return res.Error
}

// update the storage's status
func (s *StorageService) UpdateStorageStatus(storageKey string, newStatus bool, deadLine time.Time) error {
	db := global.MysqlDB
	res := db.Model(&model.Storage{}).Where("storage_url_name = ?", storageKey).Updates(map[string]interface{}{
		"status":    newStatus,
		"dead_line": deadLine,
	})
	if res.Error == nil {
		// get token
		token <- struct{}{}
		close(done)
		findEle := func(validTime validPair, status bool) {
			for k, v := range storageSet {
				if !v {continue}
				if k.StorageUrlName == validTime.StorageUrlName {
					storageSet[validTime] = status
					return
				}
			}
		}
		findEle(validPair{StorageUrlName: storageKey, DeadLine: deadLine}, newStatus)
		done = make(chan struct{})
		for k, v := range storageSet {
			if !v {
				continue
			}
			go validTimer(k)
		}
		<- token
	}
	return res.Error
}

// delete the storage
// files and foldes will be delete if they in this storage
func (s *StorageService) DeleteStorage(storageKey string) error {
	db := global.MysqlDB
	storage := model.Storage{}
	// get Storage Id
	if t, err := getStorageId(storageKey); err != nil {
		return err
	} else {
		storage.ID = t
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
		if tmp := tx.Where("storage_id is NULL").Unscoped().Delete(&model.File{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		if tmp := tx.Where("storage_id is NULL").Unscoped().Delete(&model.Folder{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		// delete the storage
		if tmp := tx.Where("storage_url_name = ?", storageKey).Unscoped().Delete(&model.Storage{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		return nil
	})
	return err
}

// query the storage file by storageKey and path
func (s *StorageService) QueryFiles(storageKey, path string) ([]response.StorageFileList, error) {
	// query the file
	var res []response.StorageFileList
	var storageId, folderId uint
	var err error
	files, folders := []model.File{}, []model.Folder{}
	db := global.MysqlDB
	if storageId, err = getStorageId(storageKey); err != nil {
		return res, err
	}
	if folderId, err = getFolderId(path, storageId); err != nil {
		return res, err
	}
	// get all files from storage
	if tmp := db.Where("storage_id = ? and folder_id = ?", storageId, folderId).Find(&files); tmp.Error != nil {
		return res, tmp.Error
	}
	if tmp := db.Where("storage_id = ? and parent_folder_id = ?", storageId, folderId).Find(&folders); tmp.Error != nil {
		return res, tmp.Error
	}
	// return the result
	for _, v := range files {
		res = append(res, response.StorageFileList{
			UpdateAt: v.UpdatedAt,
			FName:    v.FileName,
			FSize:    v.FileSize,
			FType:    response.File,
		})
	}
	for _, v := range folders {
		res = append(res, response.StorageFileList{
			UpdateAt: v.UpdatedAt,
			FName:    v.FolderName,
			FSize:    0,
			FType:    response.Folder,
		})
	}
	return res, nil
}

// query the storage information
func (s *StorageService) QueryStorageInfo() (res []response.StorageInfo, err error) {
	db := global.MysqlDB
	t := []model.Storage{}
	tmp := db.Select("ID", "StorageName", "StorageUrlName", "Status").Find(&t)
	if tmp.RowsAffected == 0 {
		err = errors.New("this system don't have storage")
		return
	}
	for _, v := range t {
		res = append(res, response.StorageInfo{
			StorageName: v.StorageName,
			StorageKey:  v.StorageUrlName,
			Status:      v.Status,
		})
	}
	return
}

// query the storage real path
func (s *StorageService) QueryStorageRealPath(storageKey string) (res string, err error) {
	db := global.MysqlDB
	t := model.Storage{}
	if global.RealPath == nil {
		global.RealPath = make(map[string]string)
	}
	if _, ok := global.RealPath[storageKey]; ok {
		res, err = global.RealPath[storageKey], nil
		return
	}
	tmp := db.Select("StorageRealPath").Where("storage_url_name = ?", storageKey).Find(&t)
	if tmp.RowsAffected == 0 {
		err = errors.New("this system don't have storage")
		return
	}
	res, err = t.StorageRealPath, nil
	global.RealPath[storageKey] = res
	return
}

// excute when server starts
// Close the storage source at a specific time
func InitTimer() {
	// storage change status when current time exceeds the deadline
	db := global.MysqlDB
	db.Model(&model.Storage{}).Where("status = ?", true).Find(&storagelist)
	for _, storage := range storagelist {
		storage := storage
		storageSet[storage] = true
		go validTimer(storage)
	}
	go receiveStorageKey()
}

// storage timer
func validTimer(valid validPair) {
	timer := time.NewTimer(time.Until(valid.DeadLine))
	select {
	case <-done:
		timer.Stop()
		return
	case <-timer.C:
	}
	storagedeadline <- valid
}

func receiveStorageKey() {
	s := &StorageService{}
	for {
		storage := <-storagedeadline
		storageSet[storage] = false
		t, _ := time.Parse(global.Format, "9999-01-01 00:00:00(CST)")
		s.UpdateStorageStatus(storage.StorageUrlName, false, t)
	}
}
