package system

import (
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"
	"fileCollect/model/system/response"

	"gorm.io/gorm"
)

type StorageService struct{}

// create the storage
func (s *StorageService) CreateStorage(storageName, storageUrlName, storageRealPath string) error {
	db := global.MysqlDB
	res := db.Create(&model.Storage{
		StorageName: storageName,
		StorageUrlName: storageUrlName,
		StorageRealPath: storageRealPath,
		Status: true,
	})
	return res.Error
}

// update the storage's Name
func (s *StorageService) UpdateStorageName(id uint, newName string) error {
	db := global.MysqlDB
	if err := checkRecordById(id, model.Storage{}); err != nil {
		return err
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("StorageName", newName)
	return res.Error
}

// update the storage's url name
func (s *StorageService) UpdateStorageUrlName(id uint, newUrlName string) error {
	db := global.MysqlDB
	if err := checkRecordById(id, model.Storage{}); err != nil {
		return err 
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("StorageUrlName", newUrlName)
	return res.Error
}

// update the storage's path
func (s *StorageService) UpdateStoragePath(id uint, newPath string) error {
	db := global.MysqlDB
	if err := checkRecordById(id, model.Storage{}); err != nil {
		return err
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("StorageRealPath", newPath)
	return res.Error
}

// update the storage's status
func (s *StorageService) UpdateStorageStatus(id uint, newStatus bool) error {
	db := global.MysqlDB
	if err := checkRecordById(id, model.Storage{}); err != nil {
		return err
	}
	res := db.Model(&model.Storage{}).Where("id = ?", id).Update("Status", newStatus)
	return res.Error
}

// delete the storage
// files and foldes will be delete if they in this storage
func (s *StorageService) DeleteStorage(id uint) error {
	db := global.MysqlDB
	storage := model.Storage{
		Model: gorm.Model{ID : id},
	}
	// clear the relation
	// start transaction
	err := db.Transaction(func (tx *gorm.DB) error {
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
			Model: v.Model,
			FName: v.FileName,
			FSize: v.FileSize,
			FType: response.File,
		})
	}
	for _, v := range folders {
		res = append(res, response.StorageFileList{
			Model: v.Model,
			FName: v.FolderName,
			FSize: 0,
			FType: response.Folder,
		})
	}
	return res, nil
}

// query the storage information
func (s *StorageService) QueryStorageInfo() (res []response.StorageInfo, err error) {
	db := global.MysqlDB
	t := []model.Storage{}
	tmp := db.Select("ID", "StorageName", "StorageRealPath", "StorageUrlName", "Status").Find(&t)
	if tmp.RowsAffected == 0 {
		err = errors.New("this system don't have storage")
		return
	}
	for _, v := range t {
		res = append(res, response.StorageInfo{
			Id: v.ID,
			StorageName: v.StorageName,
			StorageUrlName: v.StorageUrlName,
			StorageRealPath: v.StorageRealPath,
			// the feild express the storate root catalogue
			Path: "/",
			Status: v.Status,
		})
	}
	return
}