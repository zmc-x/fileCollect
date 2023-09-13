package system

import (
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"
)

type FileService struct{}

// store files' information into database
func (s *FileService) StoreFile(fileSize uint, fileName, storageKey, path string) error {
	db := global.MysqlDB
	var storageId, folderId uint
	var err error
	if storageId, err = getStorageId(storageKey); err != nil {
		return err
	}
	if folderId, err = getFolderId(path, storageId); err != nil {
		return err
	}
	// Check whether the directory stores the same name
	if temp := db.Where("storage_id = ? and parent_folder_id = ? and folder_name = ?", storageId, folderId, fileName).Find(&model.Folder{}); temp.RowsAffected != 0 {
		return errors.New("this directory contains a file or folder with the same name")
	}
	res := db.Create(&model.File{
		FileSize:  fileSize,
		FolderId:  folderId,
		StorageId: storageId,
		FileName:  fileName,
	})
	return res.Error
}

// update file related information
func (s *FileService) UpdateFileName(storageKey, path, newName, fileName string) error {
	db := global.MysqlDB
	var storageId, folderId uint 
	var err error
	if storageId, err = getStorageId(storageKey); err != nil {
		return err
	}
	if folderId, err = getFolderId(path, storageId); err != nil {
		return err
	}
	res := db.Model(&model.File{}).Where("storage_id = ? and file_name = ? and folder_id = ?", storageId, fileName, folderId).Update("FileName", newName)
	return res.Error
}

// delete file record
func (s *FileService) DeleteFile(storageKey, fileName, path string) error {
	db := global.MysqlDB
	var storageId, folderId uint 
	var err error
	if storageId, err = getStorageId(storageKey); err != nil {
		return err
	}
	if folderId, err = getFolderId(path, storageId); err != nil {
		return err
	}
	res := db.Where("storage_id = ? and file_name = ? and folder_id = ?", storageId, fileName, folderId).Unscoped().Delete(&model.File{})
	return res.Error
}
