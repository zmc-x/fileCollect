package system

import (
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"

	"gorm.io/gorm"
)

type FolderService struct{}

// create the folder in the storage
func (fs *FolderService) CreateFolder(folderName, storageKey, path string) error {
	db := global.MysqlDB
	var storageId, folderId uint 
	var err error
	if storageId, err = getStorageId(storageKey); err != nil {
		return err
	}
	if folderId, err = getFolderId(path, storageId); err != nil {
		return err
	}
	// Check whether there is the same name as the file
	if tmp := db.Where("storage_id = ? and file_name = ? and folder_id = ?", storageId, folderName, folderId).Find(&model.File{}); tmp.RowsAffected != 0 {
		return errors.New("this directory contains a file or folder with the same name")
	}
	res := db.Create(&model.Folder{
		FolderName: folderName,
		StorageId: storageId,
		ParentFolderId: folderId,
	})
	return res.Error
}

// update the folder Name
func (fs *FolderService) UpdateFolderName(folderName, path, storageKey, newName string) error {
	db := global.MysqlDB
	var storageId, folderId uint 
	var err error
	if storageId, err = getStorageId(storageKey); err != nil {
		return err
	}
	if folderId, err = getFolderId(path, storageId); err != nil {
		return err
	}
	res := db.Model(&model.Folder{}).Where("storage_id = ? and folder_name = ? and parent_folder_id = ?", storageId, folderName, folderId).Update("FolderName", newName)
	return res.Error
}


// delete the folder information 
func (fs *FolderService) DeleteFolder(folderName, path, storageKey string) error {
	db := global.MysqlDB
	var storageId, parFolderId uint
	var folder model.Folder 
	var err error
	if storageId, err = getStorageId(storageKey); err != nil {
		return err
	}
	if parFolderId, err = getFolderId(path, storageId); err != nil {
		return err
	}
	// Look for the number of the folder named folderName
	if temp := db.Where("storage_id = ? and folder_name = ? and parent_folder_id = ?", storageId, folderName, parFolderId).Find(&folder); temp.RowsAffected == 0 {
		return errors.New("the corresponding record cannot be queried")
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&folder).Association("Files").Clear(); err != nil {
			return err 
		}
		if tmp := tx.Where("folder_id is NULL").Unscoped().Delete(&model.File{}); tmp.Error != nil {
			return errors.New("these records delete error")
		}
		if tmp := tx.Unscoped().Delete(&model.Folder{}, folder.ID); tmp.Error != nil {
			return tmp.Error
		}
		// commit
		return nil 
	})
	return err
}

// Check whether the directory exists in the system
func (fs *FolderService) FolderExist(storageKey, path string) (err error) {
	var storageId uint 
	storageId, _ = getStorageId(storageKey)
	_, err = getFolderId(path, storageId)
	return err
}
